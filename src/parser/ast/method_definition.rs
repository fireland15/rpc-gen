use crate::parser::{Parse, ParseError, TokenStream};

use super::identifier::Identifier;

pub(crate) struct MethodDefinition {
    name: Identifier,
}

impl<'a> Parse<'a> for MethodDefinition {
    fn parse(tokens: &mut TokenStream<'a>) -> Result<Self, ParseError> {
        todo!()
    }
}
