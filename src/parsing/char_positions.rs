use std::str::CharIndices;

pub struct CharPositions<'a> {
    inner: CharIndices<'a>,
    line: usize,
    column: usize,
    prev: Option<char>,
}

impl<'a> CharPositions<'a> {
    fn new(char_indices: CharIndices<'a>) -> Self {
        Self {
            inner: char_indices,
            line: 0,
            column: 0,
            prev: None,
        }
    }
}

#[derive(Clone, Debug, PartialEq)]
pub struct Position {
    pub line: usize,
    pub column: usize,
    pub index: usize,
}

impl Position {
    pub fn new(line: usize, column: usize, index: usize) -> Self {
        Self {
            line,
            column,
            index,
        }
    }
}

impl<'a> Iterator for CharPositions<'a> {
    type Item = (char, Position);

    fn next(&mut self) -> Option<Self::Item> {
        let (index, ch) = self.inner.next()?;

        let column = match self.prev {
            Some('\n') => {
                self.line += 1;
                self.column = 0;
                0
            }
            Some(_) => {
                self.column += 1;
                self.column
            }
            None => self.column,
        };

        self.prev = Some(ch);

        return Some((
            ch,
            Position {
                line: self.line,
                column,
                index,
            },
        ));
    }
}

pub trait CharPositionIterators<'a> {
    fn char_positions(self) -> CharPositions<'a>;
}

impl<'a> CharPositionIterators<'a> for CharIndices<'a> {
    fn char_positions(self) -> CharPositions<'a> {
        CharPositions::new(self)
    }
}

#[cfg(test)]
mod tests {
    use crate::parsing::char_positions::Position;

    use super::CharPositions;

    #[test]
    fn iterates() {
        let mut char_positions = CharPositions::new("abcd\nef".char_indices());

        assert_eq!(char_positions.next(), Some(('a', Position::new(0, 0, 0))));
        assert_eq!(char_positions.next(), Some(('b', Position::new(0, 1, 1))));
        assert_eq!(char_positions.next(), Some(('c', Position::new(0, 2, 2))));
        assert_eq!(char_positions.next(), Some(('d', Position::new(0, 3, 3))));
        assert_eq!(char_positions.next(), Some(('\n', Position::new(0, 4, 4))));
        assert_eq!(char_positions.next(), Some(('e', Position::new(1, 0, 5))));
        assert_eq!(char_positions.next(), Some(('f', Position::new(1, 1, 6))));
    }
}
