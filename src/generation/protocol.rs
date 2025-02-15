pub struct Protocol {}

impl Protocol {
    pub fn new() -> Self {
        Self {}
    }
}

#[derive(Debug, PartialEq, Clone)]
pub enum Type {
    Named {
        name: string,
    },
    Optional {
        inner: Box<Type>,
    },
    Array {
        inner: Box<Type>,
    },
    Generic {
        inner: Box<Type>,
        parameters: Vec<Type>,
    },
}

pub enum NamedType {
    Scalar,
    Object {
        fields: Vec<
    }
}