use identifier::Identifier;

use crate::lexer::token::{Span, Token, TokenType};

use super::{Parse, ParseError};

mod identifier;
mod method_definition;
mod r#type;

pub(crate) struct ServiceDefinition {}

impl<'a> Parse<'a> for ServiceDefinition {
    fn parse(tokens: &mut super::TokenStream<'a>) -> Result<Self, ParseError> {
        todo!()
    }
}
