use std::{iter::Peekable, str::FromStr};

use super::{
    char_positions::CharPositions,
    token::{KeywordKind, Span, Token, TokenKind},
};

pub struct Tokens<'a> {
    inner: Peekable<CharPositions<'a>>,
}

impl<'a> Tokens<'a> {
    pub fn new(inner: CharPositions<'a>) -> Self {
        Self {
            inner: inner.peekable(),
        }
    }

    fn consume_whitespace(&mut self) {
        while self.inner.next_if(|(ch, _)| ch.is_whitespace()).is_some() {
            // no-op
        }
    }
}

impl<'a> Iterator for Tokens<'a> {
    type Item = Token;

    fn next(&mut self) -> Option<Self::Item> {
        self.consume_whitespace();

        let Some((ch, start)) = self.inner.next() else {
            return None;
        };

        match ch {
            ch if unicode_xid::UnicodeXID::is_xid_start(ch) => {
                let mut s = String::from(ch);
                let mut end = start.clone();
                while let Some((ch, pos)) = self
                    .inner
                    .next_if(|(ch, _)| unicode_xid::UnicodeXID::is_xid_continue(*ch))
                {
                    s.push(ch);
                    end = pos
                }

                if let Ok(kind) = KeywordKind::from_str(&s) {
                    Some(Token {
                        kind: TokenKind::Keyword(kind),
                        span: Span::starts(start).ends(end),
                    })
                } else {
                    Some(Token {
                        kind: TokenKind::Identifier,
                        span: Span::starts(start).ends(end),
                    })
                }
            }
            '{' => Some(Token {
                kind: TokenKind::LeftSquiggle,
                span: Span::starts(start),
            }),
            '}' => Some(Token {
                kind: TokenKind::RightSquiggle,
                span: Span::starts(start),
            }),
            '(' => Some(Token {
                kind: TokenKind::LeftParenthesis,
                span: Span::starts(start),
            }),
            ')' => Some(Token {
                kind: TokenKind::RightParenthesis,
                span: Span::starts(start),
            }),
            '<' => Some(Token {
                kind: TokenKind::LeftAngle,
                span: Span::starts(start),
            }),
            '>' => Some(Token {
                kind: TokenKind::RightAngle,
                span: Span::starts(start),
            }),
            '[' => Some(Token {
                kind: TokenKind::LeftSquare,
                span: Span::starts(start),
            }),
            ']' => Some(Token {
                kind: TokenKind::RightSquare,
                span: Span::starts(start),
            }),
            '?' => Some(Token {
                kind: TokenKind::Question,
                span: Span::starts(start),
            }),
            ',' => Some(Token {
                kind: TokenKind::Comma,
                span: Span::starts(start),
            }),
            _ => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::parsing::{
        char_positions::{CharPositionIterators, Position},
        token::{KeywordKind, Span, Token, TokenKind},
    };

    use super::Tokens;

    #[test]
    fn consume_whitespace() {
        let mut tokens = Tokens::new("    a".char_indices().char_positions());
        tokens.consume_whitespace();

        assert_eq!(tokens.inner.next(), Some(('a', Position::new(0, 4, 4))));
    }

    #[test]
    fn tokenizes_identifier() {
        let mut tokens = Tokens::new("    a".char_indices().char_positions());
        let next = tokens.next().expect("a token");
        assert_eq!(
            next,
            Token {
                kind: TokenKind::Identifier,
                span: Span::starts(Position::new(0, 4, 4)).ends(Position::new(0, 4, 4))
            }
        )
    }

    #[test]
    fn tokenizes_keywords() {
        let mut tokens = Tokens::new("query model".char_indices().char_positions());
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::Keyword(KeywordKind::Query),
                span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 4, 4))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::Keyword(KeywordKind::Model),
                span: Span::starts(Position::new(0, 6, 6)).ends(Position::new(0, 10, 10))
            }
        );
        assert_eq!(tokens.next(), None);
    }

    #[test]
    fn tokenizes_symbols() {
        let mut tokens = Tokens::new("{ } < > ( ) [ ] ? ,".char_indices().char_positions());
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::LeftSquiggle,
                span: Span::starts(Position::new(0, 0, 0)).ends(Position::new(0, 0, 0))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::RightSquiggle,
                span: Span::starts(Position::new(0, 2, 2)).ends(Position::new(0, 2, 2))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::LeftAngle,
                span: Span::starts(Position::new(0, 4, 4)).ends(Position::new(0, 4, 4))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::RightAngle,
                span: Span::starts(Position::new(0, 6, 6)).ends(Position::new(0, 6, 6))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::LeftParenthesis,
                span: Span::starts(Position::new(0, 8, 8)).ends(Position::new(0, 8, 8))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::RightParenthesis,
                span: Span::starts(Position::new(0, 10, 10)).ends(Position::new(0, 10, 10))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::LeftSquare,
                span: Span::starts(Position::new(0, 12, 12)).ends(Position::new(0, 12, 12))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::RightSquare,
                span: Span::starts(Position::new(0, 14, 14)).ends(Position::new(0, 14, 14))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::Question,
                span: Span::starts(Position::new(0, 16, 16)).ends(Position::new(0, 16, 16))
            }
        );
        assert_eq!(
            tokens.next().expect("a token"),
            Token {
                kind: TokenKind::Comma,
                span: Span::starts(Position::new(0, 18, 18)).ends(Position::new(0, 18, 18))
            }
        );
        assert_eq!(tokens.next(), None);
    }

    #[test]
    fn returns_none_at_end() {
        let mut tokens = Tokens::new("    ".char_indices().char_positions());
        assert_eq!(None, tokens.next());
    }
}
