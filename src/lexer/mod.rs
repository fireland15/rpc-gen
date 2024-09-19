use std::iter::Peekable;

use char_positions::CharPositions;
use token::{to_keyword, Span, Token, TokenType};
use unicode_xid::UnicodeXID;

pub mod char_positions;
pub mod token;

pub(crate) struct Lexer<'a> {
    char_positions: Peekable<CharPositions<'a>>,
}

impl<'a> Lexer<'a> {
    pub(crate) fn new(char_positions: CharPositions<'a>) -> Self {
        Self {
            char_positions: char_positions.peekable(),
        }
    }
}

impl<'a> Iterator for Lexer<'a> {
    type Item = Token;

    fn next(&mut self) -> Option<Self::Item> {
        while let Some(_) = self.char_positions.next_if(|(_, ch)| ch.is_whitespace()) {
            // consuming whitespace
        }

        let res = match self.char_positions.next()? {
            (pos, ch) if UnicodeXID::is_xid_start(ch) => {
                let start = pos.clone();
                let mut end = start.clone();
                let mut text = String::new();
                text.push(ch);

                while let Some((pos, ch)) = self
                    .char_positions
                    .next_if(|(_, ch)| UnicodeXID::is_xid_continue(*ch))
                {
                    end = pos.clone();
                    text.push(ch);
                }

                if let Some(kw) = to_keyword(&text) {
                    Token::new(TokenType::Keyword(kw), Span::range(start, end))
                } else {
                    Token::new(TokenType::Identifier(text), Span::range(start, end))
                }
            }
            (pos, '{') => Token::new(TokenType::LeftSquiggle, Span::at(pos)),
            (pos, '}') => Token::new(TokenType::RightSquiggle, Span::at(pos)),
            (pos, '[') => Token::new(TokenType::LeftSquare, Span::at(pos)),
            (pos, ']') => Token::new(TokenType::RightSquare, Span::at(pos)),
            (pos, '(') => Token::new(TokenType::LeftParenthesis, Span::at(pos)),
            (pos, ')') => Token::new(TokenType::RightParenthesis, Span::at(pos)),
            (pos, '?') => Token::new(TokenType::Question, Span::at(pos)),
            (pos, ch) => Token::new(TokenType::UnrecognizedCharacter(ch), Span::at(pos)),
        };
        Some(res)
    }
}

#[cfg(test)]
mod tests {
    use crate::lexer::{
        char_positions::Position,
        token::{Span, Token, TokenType},
    };

    use super::{char_positions::CharPositions, Lexer};

    #[test]
    fn tokenizes() {
        let mut lexer = Lexer::new(CharPositions::new("apples ( ) {} [] ?".char_indices()));
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(
                TokenType::Identifier("apples".into()),
                Span::range(Position::new(0, 0, 0), Position::new(0, 5, 5))
            )
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(TokenType::LeftParenthesis, Span::at(Position::new(0, 7, 7)))
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(
                TokenType::RightParenthesis,
                Span::at(Position::new(0, 9, 9))
            )
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(TokenType::LeftSquiggle, Span::at(Position::new(0, 11, 11)))
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(TokenType::RightSquiggle, Span::at(Position::new(0, 12, 12)))
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(TokenType::LeftSquare, Span::at(Position::new(0, 14, 14)))
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(TokenType::RightSquare, Span::at(Position::new(0, 15, 15)))
        );
        assert_eq!(
            lexer.next().unwrap(),
            Token::new(TokenType::Question, Span::at(Position::new(0, 17, 17)))
        );
        assert_eq!(lexer.next(), None);
    }
}
