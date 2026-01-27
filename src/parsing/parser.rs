use std::iter::Peekable;

use super::{
    ast::{
        MethodDefinition, MethodParameter, ModelDefinition, ProtocolDefinition, ScalarDefinition,
        Type,
    },
    token::{KeywordKind, Span, Token, TokenKind},
    tokens::Tokens,
};
use crate::parsing::ast::{EnumDefinition, EnumVariant};
use crate::parsing::{
    ast::{Identifier, ModelFieldDefinition},
    char_positions::CharPositionIterators,
};
use crate::protocol;
use crate::protocol::MethodKind;

pub fn parse(source: &str) -> Result<ProtocolDefinition, Vec<ParseError>> {
    let mut parser = Parser {
        source,
        tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
    };
    parser.parse()
}

#[derive(Debug, Clone, PartialEq)]
pub enum ParseError {
    UnexpectedEndOfInput,
    UnexpectedSymbol {
        expected: TokenKind,
        actual: TokenKind,
        actual_text: String,
        span: Span,
    },
}

struct Parser<'a> {
    source: &'a str,
    tokens: Peekable<Tokens<'a>>,
}

impl<'a> Parser<'a> {
    fn parse(&mut self) -> Result<ProtocolDefinition, Vec<ParseError>> {
        let mut service = ProtocolDefinition {
            models: Vec::new(),
            methods: Vec::new(),
            scalars: Vec::new(),
            enums: Vec::new(),
        };
        let mut errors = Vec::new();
        while let Some(token) = self.tokens.peek() {
            match token {
                Token {
                    kind: TokenKind::Keyword(KeywordKind::Model),
                    ..
                } => {
                    match self.parse_model_definition() {
                        Ok(def) => service.models.push(def),
                        Err(err) => errors.push(err),
                    };
                }
                Token {
                    kind: TokenKind::Keyword(KeywordKind::Scalar),
                    ..
                } => match self.parse_scalar_definition() {
                    Ok(def) => service.scalars.push(def),
                    Err(err) => errors.push(err),
                },
                Token {
                    kind: TokenKind::Keyword(KeywordKind::Enum),
                    ..
                } => {
                    match self.parse_enum_definition() {
                        Ok(def) => service.enums.push(def),
                        Err(err) => errors.push(err),
                    };
                }
                Token {
                    kind: TokenKind::Keyword(KeywordKind::Query),
                    ..
                } => {
                    match self.parse_method_definition(KeywordKind::Query) {
                        Ok(def) => service.methods.push(def),
                        Err(err) => errors.push(err),
                    };
                }
                Token {
                    kind: TokenKind::Keyword(KeywordKind::Mutation),
                    ..
                } => {
                    match self.parse_method_definition(KeywordKind::Mutation) {
                        Ok(def) => service.methods.push(def),
                        Err(err) => errors.push(err),
                    };
                }
                Token {
                    kind: TokenKind::Keyword(KeywordKind::Stream),
                    ..
                } => {
                    match self.parse_method_definition(KeywordKind::Stream) {
                        Ok(def) => service.methods.push(def),
                        Err(err) => errors.push(err),
                    };
                }
                _ => {
                    let token = self.tokens.next().expect("already checked that is Some");
                    errors.push(self.unexpected(TokenKind::Keyword(KeywordKind::Model), &token));
                    errors.push(self.unexpected(TokenKind::Keyword(KeywordKind::Query), &token));
                    errors.push(self.unexpected(TokenKind::Keyword(KeywordKind::Mutation), &token));
                    errors.push(self.unexpected(TokenKind::Keyword(KeywordKind::Stream), &token));
                    errors.push(self.unexpected(TokenKind::Keyword(KeywordKind::Scalar), &token));
                }
            }
        }
        if errors.len() > 0 {
            return Err(errors);
        }
        Ok(service)
    }

    fn unexpected(&self, expected: TokenKind, token: &Token) -> ParseError {
        ParseError::UnexpectedSymbol {
            expected,
            actual: token.kind.clone(),
            actual_text: self.get_text(&token.span).into(),
            span: token.span.clone(),
        }
    }

    fn get_text(&self, span: &Span) -> &'a str {
        &self.source[span.start.index..span.end.index + 1]
    }

    fn parse_scalar_definition(&mut self) -> Result<ScalarDefinition, ParseError> {
        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Keyword(KeywordKind::Scalar),
            ..
        } = token
        else {
            return Err(self.unexpected(TokenKind::Keyword(KeywordKind::Scalar), &token));
        };

        let identifier_token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Identifier,
            ..
        } = identifier_token
        else {
            return Err(self.unexpected(TokenKind::Identifier, &identifier_token));
        };

        let serialized_token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Identifier,
            ..
        } = serialized_token
        else {
            return Err(self.unexpected(TokenKind::Identifier, &serialized_token));
        };

        Ok(ScalarDefinition {
            name: Identifier::new(identifier_token.span, &self.source),
            serialized: Identifier::new(serialized_token.span, &self.source),
        })
    }

    fn parse_method_definition(
        &mut self,
        method_kind: KeywordKind,
    ) -> Result<MethodDefinition, ParseError> {
        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Keyword(kk),
            ..
        } = &token
        else {
            return Err(self.unexpected(TokenKind::Keyword(method_kind), &token));
        };

        if *kk != method_kind {
            return Err(self.unexpected(TokenKind::Keyword(method_kind), &token));
        }

        let identifier_token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Identifier,
            ..
        } = identifier_token
        else {
            return Err(self.unexpected(TokenKind::Identifier, &identifier_token));
        };

        let Some(Token {
            kind: TokenKind::LeftParenthesis,
            ..
        }) = self.tokens.peek()
        else {
            return Err(self.unexpected(TokenKind::LeftParenthesis, &token));
        };
        self.tokens.next();

        // parse parameters

        let mut parameters = Vec::new();

        loop {
            // Identifier
            let Some(Token {
                kind: TokenKind::Identifier,
                ..
            }) = self.tokens.peek()
            else {
                break;
            };
            let token = self.tokens.next().expect("already checked this was Some");
            let identifier = Identifier::new(token.span, &self.source);

            // Type
            let ty = self.parse_type()?;

            parameters.push(MethodParameter {
                name: identifier,
                ty,
            });

            // Optional comma
            if self
                .tokens
                .next_if(|t| t.kind == TokenKind::Comma)
                .is_none()
            {
                break;
            }
        }

        let Some(Token {
            kind: TokenKind::RightParenthesis,
            ..
        }) = self.tokens.peek()
        else {
            return Err(self.unexpected(TokenKind::RightParenthesis, &token));
        };
        self.tokens.next();

        if self
            .tokens
            .peek()
            .is_some_and(|t| t.kind == TokenKind::Identifier)
        {
            let return_ty = self.parse_type()?;
            Ok(MethodDefinition {
                name: Identifier::new(identifier_token.span, &self.source),
                kind: to_method_kind(&method_kind),
                parameters,
                return_ty: Some(return_ty),
            })
        } else {
            Ok(MethodDefinition {
                name: Identifier::new(identifier_token.span, &self.source),
                kind: to_method_kind(&method_kind),
                parameters,
                return_ty: None,
            })
        }
    }

    fn parse_model_definition(&mut self) -> Result<ModelDefinition, ParseError> {
        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Keyword(KeywordKind::Model),
            ..
        } = token
        else {
            return Err(self.unexpected(TokenKind::Keyword(KeywordKind::Model), &token));
        };

        let identifier_token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Identifier,
            ..
        } = identifier_token
        else {
            return Err(self.unexpected(TokenKind::Identifier, &identifier_token));
        };

        let Some(Token {
            kind: TokenKind::LeftSquiggle,
            ..
        }) = self.tokens.peek()
        else {
            return Err(self.unexpected(TokenKind::LeftSquiggle, &token));
        };
        self.tokens.next();

        // parse field definitions

        let mut fields = Vec::new();

        loop {
            let Some(Token {
                kind: TokenKind::Identifier,
                ..
            }) = self.tokens.peek()
            else {
                break;
            };
            let token = self.tokens.next().expect("already checked this was Some");
            let identifier = Identifier::new(token.span, &self.source);

            let ty = self.parse_type()?;

            fields.push(ModelFieldDefinition {
                name: identifier,
                ty,
            });
        }

        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::RightSquiggle,
            ..
        } = token
        else {
            return Err(self.unexpected(TokenKind::RightSquiggle, &token));
        };

        return Ok(ModelDefinition {
            name: Identifier::new(identifier_token.span, &self.source),
            fields,
        });
    }

    fn parse_enum_definition(&mut self) -> Result<EnumDefinition, ParseError> {
        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Keyword(KeywordKind::Enum),
            ..
        } = token
        else {
            return Err(self.unexpected(TokenKind::Keyword(KeywordKind::Enum), &token));
        };

        let identifier_token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Identifier,
            ..
        } = identifier_token
        else {
            return Err(self.unexpected(TokenKind::Identifier, &identifier_token));
        };

        let Some(Token {
            kind: TokenKind::LeftSquiggle,
            ..
        }) = self.tokens.peek()
        else {
            return Err(self.unexpected(TokenKind::LeftSquiggle, &token));
        };
        self.tokens.next();

        // parse enum variants

        let mut variants = Vec::new();

        loop {
            let Some(Token {
                kind: TokenKind::Identifier,
                ..
            }) = self.tokens.peek()
            else {
                break;
            };
            let token = self.tokens.next().expect("already checked this was Some");
            let identifier = Identifier::new(token.span, &self.source);

            variants.push(EnumVariant { name: identifier });
        }

        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::RightSquiggle,
            ..
        } = token
        else {
            return Err(self.unexpected(TokenKind::RightSquiggle, &token));
        };

        return Ok(EnumDefinition {
            name: Identifier::new(identifier_token.span, &self.source),
            variants,
        });
    }

    fn parse_type(&mut self) -> Result<Type, ParseError> {
        let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            kind: TokenKind::Identifier,
            ..
        } = token
        else {
            return Err(self.unexpected(TokenKind::Identifier, &token));
        };

        let mut ty = Type::Type {
            identifier: Identifier::new(token.span, &self.source),
        };

        loop {
            match self.tokens.peek() {
                Some(Token {
                    kind: TokenKind::LeftSquare,
                    ..
                }) => {
                    self.tokens.next();
                    let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
                    let Token {
                        kind: TokenKind::RightSquare,
                        ..
                    } = token
                    else {
                        return Err(self.unexpected(TokenKind::Identifier, &token));
                    };

                    ty = Type::Array {
                        inner: Box::new(ty),
                    };
                }
                Some(Token {
                    kind: TokenKind::LeftAngle,
                    ..
                }) => {
                    self.tokens.next();

                    let mut generic_args = Vec::new();
                    while let Some(Token {
                        kind: TokenKind::Identifier,
                        ..
                    }) = self.tokens.peek()
                    {
                        generic_args.push(self.parse_type()?);
                        if self
                            .tokens
                            .next_if(|t| t.kind == TokenKind::Comma)
                            .is_none()
                        {
                            break;
                        }
                    }

                    let token = self.tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
                    let Token {
                        kind: TokenKind::RightAngle,
                        ..
                    } = token
                    else {
                        return Err(self.unexpected(TokenKind::RightAngle, &token));
                    };

                    ty = Type::Generic {
                        inner: Box::new(ty),
                        parameters: generic_args,
                    };
                }
                Some(Token {
                    kind: TokenKind::Question,
                    ..
                }) => {
                    self.tokens.next();
                    ty = Type::Optional {
                        inner: Box::new(ty),
                    };
                }
                _ => break,
            }
        }

        return Ok(ty);
    }
}

fn to_method_kind(kw: &KeywordKind) -> MethodKind {
    match kw {
        KeywordKind::Query => MethodKind::Query,
        KeywordKind::Mutation => MethodKind::Mutation,
        KeywordKind::Stream => MethodKind::Stream,
        _ => panic!("you shouldn't be trying to map a keyword for methods to method kinds."),
    }
}

#[cfg(test)]
mod tests {
    use super::Parser;
    use crate::parsing::ast::{EnumDefinition, EnumVariant};
    use crate::parsing::{
        ast::{Identifier, ModelDefinition, ModelFieldDefinition, Type},
        char_positions::{CharPositionIterators, Position},
        token::{Span, Token, TokenKind},
        tokens::Tokens,
    };

    #[test]
    fn parses_plain_type() {
        let source = "Apples";
        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let ty = p.parse_type().expect("failed to parse type");
        assert_eq!(
            ty,
            Type::Type {
                identifier: Identifier {
                    text: "Apples".into(),
                    token: Token {
                        kind: TokenKind::Identifier,
                        span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 5, 5))
                    }
                }
            }
        );
    }

    #[test]
    fn parses_array_type() {
        let source = "Apples[]";
        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let ty = p.parse_type().expect("failed to parse type");
        assert_eq!(
            ty,
            Type::Array {
                inner: Box::new(Type::Type {
                    identifier: Identifier {
                        text: "Apples".into(),
                        token: Token {
                            kind: TokenKind::Identifier,
                            span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 5, 5))
                        }
                    }
                })
            }
        );
    }

    #[test]
    fn parses_optional_type() {
        let source = "Apples?";
        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let ty = p.parse_type().expect("failed to parse type");
        assert_eq!(
            ty,
            Type::Optional {
                inner: Box::new(Type::Type {
                    identifier: Identifier {
                        text: "Apples".into(),
                        token: Token {
                            kind: TokenKind::Identifier,
                            span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 5, 5))
                        }
                    }
                })
            }
        );
    }

    #[test]
    fn parses_generic_type() {
        let source = "Apples<Banana, Carrot>";
        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let ty = p.parse_type().expect("failed to parse type");
        assert_eq!(
            ty,
            Type::Generic {
                inner: Box::new(Type::Type {
                    identifier: Identifier {
                        text: "Apples".into(),
                        token: Token {
                            kind: TokenKind::Identifier,
                            span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 5, 5))
                        }
                    }
                }),
                parameters: vec![
                    Type::Type {
                        identifier: Identifier {
                            text: "Banana".into(),
                            token: Token {
                                kind: TokenKind::Identifier,
                                span: Span::starts(Position::new(0, 7, 7))
                                    .ends(Position::new(0, 12, 12))
                            }
                        }
                    },
                    Type::Type {
                        identifier: Identifier {
                            text: "Carrot".into(),
                            token: Token {
                                kind: TokenKind::Identifier,
                                span: Span::starts(Position::new(0, 15, 15))
                                    .ends(Position::new(0, 20, 20))
                            }
                        }
                    }
                ]
            }
        );
    }

    #[test]
    fn parses_complex_type() {
        let source = "Apples<Banana[]?>";
        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let ty = p.parse_type().expect("failed to parse type");
        assert_eq!(
            ty,
            Type::Generic {
                inner: Box::new(Type::Type {
                    identifier: Identifier {
                        text: "Apples".into(),
                        token: Token {
                            kind: TokenKind::Identifier,
                            span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 5, 5))
                        }
                    }
                }),
                parameters: vec![Type::Optional {
                    inner: Box::new(Type::Array {
                        inner: Box::new(Type::Type {
                            identifier: Identifier {
                                text: "Banana".into(),
                                token: Token {
                                    kind: TokenKind::Identifier,
                                    span: Span::starts(Position::new(0, 7, 7))
                                        .ends(Position::new(0, 12, 12))
                                }
                            }
                        })
                    })
                }]
            }
        );
    }

    #[test]
    fn parses_model_definition() {
        let source = r#"model Fruit {
    name string
    costPerUnit number?
    locations string[]
}"#;
        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };
        let model_def = p
            .parse_model_definition()
            .expect("failed to parse model definition");
        assert_eq!(
            model_def,
            ModelDefinition {
                name: Identifier::new(
                    Span::starts(Position::new(0, 6, 6)).ends(Position::new(0, 10, 10)),
                    &source,
                ),
                fields: vec![
                    ModelFieldDefinition::new(
                        Identifier::new(
                            Span::starts(Position::new(1, 4, 18)).ends(Position::new(1, 7, 21)),
                            &source,
                        ),
                        Type::new(
                            Span::starts(Position::new(1, 9, 23)).ends(Position::new(1, 14, 28)),
                            &source,
                        )
                    ),
                    ModelFieldDefinition::new(
                        Identifier::new(
                            Span::starts(Position::new(2, 4, 34)).ends(Position::new(2, 14, 44)),
                            &source,
                        ),
                        Type::Optional {
                            inner: Box::new(Type::new(
                                Span::starts(Position::new(2, 16, 46))
                                    .ends(Position::new(2, 21, 51)),
                                &source,
                            )),
                        }
                    ),
                    ModelFieldDefinition::new(
                        Identifier::new(
                            Span::starts(Position::new(3, 4, 58)).ends(Position::new(3, 12, 66)),
                            &source,
                        ),
                        Type::Array {
                            inner: Box::new(Type::new(
                                Span::starts(Position::new(3, 14, 68))
                                    .ends(Position::new(3, 19, 73)),
                                &source,
                            )),
                        }
                    )
                ]
            }
        )
    }

    #[test]
    fn parses_enum_definition() {
        let source = r#"enum Color {
    Red
    Green
    Blue
}"#;

        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let enum_def = p
            .parse_enum_definition()
            .expect("failed to parse enum definition");

        assert_eq!(
            enum_def,
            EnumDefinition {
                name: Identifier::new(
                    Span::starts(Position::new(0, 5, 5)).ends(Position::new(0, 9, 9)),
                    &source,
                ),
                variants: vec![
                    EnumVariant {
                        name: Identifier::new(
                            Span::starts(Position::new(1, 4, 17)).ends(Position::new(1, 6, 19)),
                            &source,
                        )
                    },
                    EnumVariant {
                        name: Identifier::new(
                            Span::starts(Position::new(2, 4, 25)).ends(Position::new(2, 8, 29)),
                            &source,
                        )
                    },
                    EnumVariant {
                        name: Identifier::new(
                            Span::starts(Position::new(3, 4, 35)).ends(Position::new(3, 7, 38)),
                            &source,
                        )
                    },
                ]
            }
        );
    }

    #[test]
    fn parses_single_variant_enum() {
        let source = r#"enum Truth {
    True
}"#;

        let mut p = Parser {
            source,
            tokens: Tokens::new(source.char_indices().char_positions()).peekable(),
        };

        let enum_def = p
            .parse_enum_definition()
            .expect("failed to parse enum definition");

        assert_eq!(enum_def.variants.len(), 1);
        assert_eq!(enum_def.variants[0].name.text, "True");
    }
}
