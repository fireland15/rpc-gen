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
    LeftSquare,
    RightSquare,
    Question,
    Comma,
}

#[derive(Clone, Debug, PartialEq)]
pub enum KeywordKind {
    Rpc,
    Model,
    Scalar,
    Enum,
}

impl FromStr for KeywordKind {
    type Err = ();

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "rpc" => Ok(KeywordKind::Rpc),
            "model" => Ok(KeywordKind::Model),
            "scalar" => Ok(KeywordKind::Scalar),
            "enum" => Ok(KeywordKind::Enum),
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

impl Token {
    pub fn identifier(span: Span) -> Self {
        Token {
            kind: TokenKind::Identifier,
            span,
        }
    }

    pub fn keyword(kind: KeywordKind, span: Span) -> Self {
        Token {
            kind: TokenKind::Keyword(kind),
            span,
        }
    }
}
