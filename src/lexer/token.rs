use super::char_positions::Position;

#[derive(Debug, PartialEq)]
pub(crate) enum TokenType {
    Identifier(String),
    Keyword(KeywordKind),
    LeftParenthesis,
    RightParenthesis,
    LeftSquiggle,
    RightSquiggle,
    LeftSquare,
    RightSquare,
    Question,
    EndOfInput,
    UnrecognizedCharacter(char),
}

impl TokenType {
    pub fn to_string(&self) -> String {
        match self {
            TokenType::Identifier(identifier) => identifier.clone(),
            TokenType::Keyword(keyword_kind) => match keyword_kind {
                KeywordKind::Model => "model".into(),
                KeywordKind::Rpc => "rpc".into(),
            },
            TokenType::LeftParenthesis => "(".into(),
            TokenType::RightParenthesis => ")".into(),
            TokenType::LeftSquiggle => "{".into(),
            TokenType::RightSquiggle => "}".into(),
            TokenType::LeftSquare => "[".into(),
            TokenType::RightSquare => "]".into(),
            TokenType::Question => "?".into(),
            TokenType::EndOfInput => "EOF".into(),
            TokenType::UnrecognizedCharacter(ch) => ch.to_string(),
        }
    }
}

#[derive(Debug, PartialEq)]
pub(crate) enum KeywordKind {
    Model,
    Rpc,
}

pub fn to_keyword(str: &str) -> Option<KeywordKind> {
    match str {
        "model" => Some(KeywordKind::Model),
        "rpc" => Some(KeywordKind::Rpc),
        _ => None,
    }
}

#[derive(Debug, PartialEq)]
pub(crate) struct Token {
    pub token_type: TokenType,
    pub span: Span,
}

#[derive(Clone, Debug, PartialEq)]
pub struct Span {
    pub start: Position,
    pub end: Position,
}

impl Span {
    pub(crate) fn range(start: Position, end: Position) -> Self {
        Self { start, end }
    }

    pub(crate) fn at(position: Position) -> Self {
        Self {
            start: position.clone(),
            end: position,
        }
    }
}

impl Token {
    pub(crate) fn new(token_type: TokenType, span: Span) -> Self {
        Self { token_type, span }
    }
}
