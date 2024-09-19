use std::str::CharIndices;

pub(crate) struct CharPositions<'a> {
    inner: CharIndices<'a>,
    line: usize,
    column: usize,
    prev: Option<char>,
}

#[derive(Clone, Debug, PartialEq)]
pub(crate) struct Position {
    line: usize,
    column: usize,
    offset: usize,
}

impl Position {
    pub(crate) fn new(line: usize, column: usize, offset: usize) -> Self {
        Self {
            line,
            column,
            offset,
        }
    }
}

impl<'a> CharPositions<'a> {
    pub(crate) fn new(inner: CharIndices<'a>) -> Self {
        Self {
            inner,
            line: 0,
            column: 0,
            prev: None,
        }
    }
}

impl<'a> Iterator for CharPositions<'a> {
    type Item = (Position, char);

    fn next(&mut self) -> Option<Self::Item> {
        let (offset, next) = self.inner.next()?;

        match self.prev {
            Some('\n') => {
                self.line += 1;
                self.column = 0;
            }
            Some(_) => {
                self.column += 1;
            }
            None => {}
        }

        self.prev = Some(next);

        let position = Position {
            line: self.line,
            column: self.column,
            offset,
        };

        return Some((position, next));
    }
}

#[cfg(test)]
mod tests {
    use crate::lexer::char_positions::Position;

    use super::CharPositions;

    #[test]
    fn char_positions_works() {
        let mut char_pos = CharPositions::new("a\nbc".char_indices());

        assert_eq!(
            char_pos.next(),
            Some((
                Position {
                    line: 0,
                    column: 0,
                    offset: 0
                },
                'a'
            ))
        );
        assert_eq!(
            char_pos.next(),
            Some((
                Position {
                    line: 0,
                    column: 1,
                    offset: 1
                },
                '\n'
            ))
        );
        assert_eq!(
            char_pos.next(),
            Some((
                Position {
                    line: 1,
                    column: 0,
                    offset: 2
                },
                'b'
            ))
        );
        assert_eq!(
            char_pos.next(),
            Some((
                Position {
                    line: 1,
                    column: 1,
                    offset: 3
                },
                'c'
            ))
        );
    }
}
