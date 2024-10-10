use std::str::FromStr;

use super::char_positions::Position;

#[derive(Clone, Debug, PartialEq)]
pub enum TokenKind {
    Identifier,
    Keyword(KeywordKind),
    LeftParenthesis,
    RightParenthesis,
    LeftSquiggle,
    RightSquiggle,
    LeftAngle,
    RightAngle,
}

#[derive(Clone, Debug, PartialEq)]
pub enum KeywordKind {
    Rpc,
    Model,
}

impl FromStr for KeywordKind {
    type Err = ();

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "rpc" => Ok(KeywordKind::Rpc),
            "model" => Ok(KeywordKind::Model),
            _ => Err(()),
        }
    }
}

#[derive(Clone, Debug, PartialEq)]
pub struct Span {
    pub start: Position,
    pub end: Position,
}

impl Span {
    pub fn starts(start: Position) -> Self {
        Self {
            end: start.clone(),
            start,
        }
    }

    pub fn ends(mut self, end: Position) -> Self {
        self.end = end;
        self
    }
}

#[derive(Clone, Debug, PartialEq)]
pub struct Token {
    pub kind: TokenKind,
    pub span: Span,
}
