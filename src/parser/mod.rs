use std::iter::Peekable;

use ast::ServiceDefinition;

use crate::lexer::{
    token::{Span, TokenType},
    Lexer,
};

mod ast;

type TokenStream<'a> = Peekable<Lexer<'a>>;

#[derive(Debug, PartialEq)]
enum ParseError {
    UnexpectedToken {
        expected_token_type: TokenType,
        received_token_type: TokenType,
        span: Span,
    },
    UnexpectedEndOfInput,
}

impl ParseError {
    fn unexpected(
        expected_token_type: TokenType,
        received_token_type: TokenType,
        span: Span,
    ) -> ParseError {
        ParseError::UnexpectedToken {
            expected_token_type,
            received_token_type,
            span,
        }
    }
}

pub(crate) struct Parser<'a> {
    tokens: TokenStream<'a>,
}

impl<'a> Parser<'a> {
    pub(crate) fn parse() -> Result<ServiceDefinition, ()> {
        todo!();
    }
}

trait Parse<'a>
where
    Self: Sized,
{
    fn parse(tokens: &mut TokenStream<'a>) -> Result<Self, ParseError>;
}
