use crate::{
    lexer::token::{Span, Token, TokenType},
    parser::{Parse, ParseError, TokenStream},
};

#[derive(Debug, PartialEq)]
pub(crate) struct Identifier {
    pub name: String,
    pub span: Span,
}

impl<'a> Parse<'a> for Identifier {
    fn parse(tokens: &mut TokenStream<'a>) -> Result<Self, ParseError> {
        let token = tokens.next().ok_or(ParseError::UnexpectedEndOfInput)?;
        let Token {
            token_type: TokenType::Identifier(name),
            span,
        } = token
        else {
            return Err(ParseError::unexpected(
                token.token_type,
                TokenType::Identifier(String::new()),
                token.span,
            ));
        };

        Ok(Identifier { name, span })
    }
}
