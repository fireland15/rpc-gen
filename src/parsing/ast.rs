use super::token::{Span, Token};
use crate::protocol::MethodKind;

#[derive(Debug, PartialEq, Clone)]
pub struct ProtocolDefinition {
    pub models: Vec<ModelDefinition>,
    pub scalars: Vec<ScalarDefinition>,
    pub enums: Vec<EnumDefinition>,
    pub methods: Vec<MethodDefinition>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct MethodDefinition {
    pub name: Identifier,
    pub kind: MethodKind,
    pub parameters: Vec<MethodParameter>,
    pub return_ty: Option<Type>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct MethodParameter {
    pub name: Identifier,
    pub ty: Type,
}

#[derive(Debug, PartialEq, Clone)]
pub struct ScalarDefinition {
    pub name: Identifier,
    pub serialized: Identifier,
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
pub struct EnumDefinition {
    pub name: Identifier,
    pub variants: Vec<EnumVariant>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct EnumVariant {
    pub name: Identifier,
}

impl EnumVariant {
    pub fn new(name: Identifier) -> Self {
        Self { name }
    }
}

#[derive(Debug, PartialEq, Clone)]
pub struct Identifier {
    pub token: Token,
    pub text: String,
}

impl Identifier {
    pub fn new(span: Span, src: &str) -> Self {
        let text = src[span.start.index..span.end.index + 1].into();
        Identifier {
            token: Token::identifier(span),
            text,
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
    pub fn new(span: Span, src: &str) -> Self {
        Self::Type {
            identifier: Identifier::new(span, src),
        }
    }
}
