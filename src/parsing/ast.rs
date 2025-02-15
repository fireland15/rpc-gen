use super::token::{Span, Token};

#[derive(Debug, PartialEq, Clone)]
pub struct ProtocolDefinition {
    pub models: Vec<ModelDefinition>,
    pub methods: Vec<MethodDefinition>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct MethodDefinition {
    pub name: Identifier,
    pub parameters: Vec<MethodParameter>,
    pub return_ty: Option<Type>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct MethodParameter {
    pub name: Identifier,
    pub ty: Type,
}

#[derive(Debug, PartialEq, Clone)]
pub struct ModelDefinition {
    pub name: Identifier,
    pub fields: Vec<ModelFieldDefinition>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct ModelFieldDefinition {
    pub name: Identifier,
    pub ty: Type,
}

impl ModelFieldDefinition {
    pub fn new(name: Identifier, ty: Type) -> Self {
        Self { name, ty }
    }
}

#[derive(Debug, PartialEq, Clone)]
pub struct Identifier {
    pub token: Token,
}

impl Identifier {
    pub fn new(span: Span) -> Self {
        Identifier {
            token: Token::identifier(span),
        }
    }
}

#[derive(Debug, PartialEq, Clone)]
pub enum Type {
    Type {
        identifier: Identifier,
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

impl Type {
    pub fn new(span: Span) -> Self {
        Self::Type {
            identifier: Identifier::new(span),
        }
    }
}
