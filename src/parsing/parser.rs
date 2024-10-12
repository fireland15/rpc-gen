use std::{clone, iter::Peekable};

use crate::parsing::{
    ast::{Identifier, ModelFieldDefinition},
    char_positions::CharPositionIterators,
};

use super::{
    ast::{ModelDefinition, ServiceDefinition, Type},
    token::{KeywordKind, Span, Token, TokenKind},
    tokens::Tokens,
};

pub fn parse(source: &str) -> Result<ServiceDefinition, ParseError> {
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
    },
}

struct Parser<'a> {
    source: &'a str,
    tokens: Peekable<Tokens<'a>>,
}

impl<'a> Parser<'a> {
    fn parse(&mut self) -> Result<ServiceDefinition, Vec<ParseError>> {
        let mut service = ServiceDefinition {
            models: Vec::new(),
            methods: Vec::new(),
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
                    kind: TokenKind::Keyword(KeywordKind::Rpc),
                    ..
                } => {
                    todo!("parse method definition");
                }
                _ => {
                    let token = self.tokens.next().expect("already checked that is Some");
                    self.unexpected(TokenKind::Keyword(KeywordKind::Model), &token);
                    self.unexpected(TokenKind::Keyword(KeywordKind::Rpc), &token);
                }
            }
        }
        todo!();
    }

    fn unexpected(&self, expected: TokenKind, token: &Token) -> ParseError {
        ParseError::UnexpectedSymbol {
            expected: expected,
            actual: token.kind.clone(),
            actual_text: self.get_text(&token.span).into(),
        }
    }

    fn get_text(&self, span: &Span) -> &'a str {
        &self.source[span.start.index..span.end.index + 1]
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
            let identifier = Identifier { token };

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
            name: Identifier {
                token: identifier_token,
            },
            fields,
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
            identifier: Identifier { token },
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

#[cfg(test)]
mod tests {
    use crate::parsing::{
        ast::{Identifier, ModelDefinition, ModelFieldDefinition, Type},
        char_positions::{CharPositionIterators, Position},
        token::{Span, Token, TokenKind},
        tokens::Tokens,
    };

    use super::Parser;

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
                        token: Token {
                            kind: TokenKind::Identifier,
                            span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 5, 5))
                        }
                    }
                }),
                parameters: vec![
                    Type::Type {
                        identifier: Identifier {
                            token: Token {
                                kind: TokenKind::Identifier,
                                span: Span::starts(Position::new(0, 7, 7))
                                    .ends(Position::new(0, 12, 12))
                            }
                        }
                    },
                    Type::Type {
                        identifier: Identifier {
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
                    Span::starts(Position::new(0, 6, 6)).ends(Position::new(0, 10, 10))
                ),
                fields: vec![
                    ModelFieldDefinition::new(
                        Identifier::new(
                            Span::starts(Position::new(1, 4, 18)).ends(Position::new(1, 7, 21))
                        ),
                        Type::new(
                            Span::starts(Position::new(1, 9, 23)).ends(Position::new(1, 14, 28))
                        )
                    ),
                    ModelFieldDefinition::new(
                        Identifier::new(
                            Span::starts(Position::new(2, 4, 34)).ends(Position::new(2, 14, 44))
                        ),
                        Type::Optional {
                            inner: Box::new(Type::new(
                                Span::starts(Position::new(2, 16, 46))
                                    .ends(Position::new(2, 21, 51))
                            )),
                        }
                    ),
                    ModelFieldDefinition::new(
                        Identifier::new(
                            Span::starts(Position::new(3, 4, 58)).ends(Position::new(3, 12, 66))
                        ),
                        Type::Array {
                            inner: Box::new(Type::new(
                                Span::starts(Position::new(3, 14, 68))
                                    .ends(Position::new(3, 19, 73))
                            )),
                        }
                    )
                ]
            }
        )
    }
}
