use crate::{
    lexer::token::{Token, TokenType},
    parser::{Parse, ParseError, TokenStream},
};

use super::identifier::Identifier;

#[derive(Debug, PartialEq)]
pub enum TypeVariant {
    Named { name: Identifier },
    Array { inner: Box<Type> },
    Optional { inner: Box<Type> },
}

#[derive(Debug, PartialEq)]
pub struct Type {
    variant: TypeVariant,
}

impl<'a> Parse<'a> for Type {
    fn parse(tokens: &mut TokenStream<'a>) -> Result<Self, ParseError> {
        let name = Identifier::parse(tokens)?;
        let ty = Type {
            variant: TypeVariant::Named { name },
        };
        return parse_inner(ty, tokens);
    }
}

fn parse_inner(ty: Type, tokens: &mut TokenStream) -> Result<Type, ParseError> {
    match tokens.peek() {
        Some(Token {
            token_type: TokenType::LeftSquare,
            ..
        }) => {
            tokens.next();
            let token = tokens.peek().ok_or(ParseError::UnexpectedEndOfInput)?;
            let Token {
                token_type: TokenType::RightSquare,
                ..
            } = token
            else {
                return Err(ParseError::unexpected(
                    TokenType::RightSquare,
                    TokenType::EndOfInput,
                    token.span.clone(),
                ));
            };
            tokens.next();
            let new_type = Type {
                variant: TypeVariant::Array { inner: ty.into() },
            };
            return parse_inner(new_type, tokens);
        }
        Some(Token {
            token_type: TokenType::Question,
            ..
        }) => {
            tokens.next();
            let new_type = Type {
                variant: TypeVariant::Optional { inner: ty.into() },
            };
            return parse_inner(new_type, tokens);
        }
        _ => return Ok(ty),
    };
}

#[cfg(test)]
mod tests {
    use crate::{
        lexer::{
            char_positions::{CharPositions, Position},
            token::Span,
            Lexer,
        },
        parser::{
            ast::{identifier::Identifier, r#type::TypeVariant},
            Parse,
        },
    };

    use super::Type;

    #[test]
    fn parses_type() {
        let mut tokens = Lexer::new(CharPositions::new("apples[]?".char_indices())).peekable();

        let ty = Type::parse(&mut tokens);
        assert_eq!(
            ty,
            Ok(Type {
                variant: TypeVariant::Optional {
                    inner: Box::new(Type {
                        variant: TypeVariant::Array {
                            inner: Box::new(Type {
                                variant: TypeVariant::Named {
                                    name: Identifier {
                                        name: "apples".into(),
                                        span: Span::range(
                                            Position::new(0, 0, 0),
                                            Position::new(0, 5, 5)
                                        )
                                    }
                                }
                            })
                        }
                    })
                }
            })
        )
    }
}
