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
            _ => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::parsing::{
        char_positions::{CharPositionIterators, Position},
        token::{Span, Token, TokenKind},
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
    fn returns_none_at_end() {
        let mut tokens = Tokens::new("    ".char_indices().char_positions());
        assert_eq!(None, tokens.next());
    }
}
