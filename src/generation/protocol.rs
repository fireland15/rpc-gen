use std::collections::BTreeMap;
use convert_case::{Case, Casing};
use serde::Serialize;

use crate::parsing::ast::{self};

pub fn build_protocol(x: &ast::ProtocolDefinition) -> Result<Protocol, ()> {
    let mut p = Protocol::new();
    x.models.iter().for_each(|m| {
        p.add_model(
            &m.name.text,
            m.fields
                .iter()
                .map(|x| (x.name.text.clone(), make_type_ref(&x.ty)))
                .collect::<BTreeMap<String, TypeRef>>(),
        )
        .unwrap();
    });
    x.scalars.iter().for_each(|s| {
        p.add_scalar(&s.name.text, &s.serialized.text).unwrap();
    });
    x.methods.iter().for_each(|method| {
        p.add_method(
            &method.name.text,
            method
                .parameters
                .iter()
                .map(|x| (x.name.text.clone(), make_type_ref(&x.ty)))
                .collect(),
            method.return_ty.as_ref().map(make_type_ref),
        )
        .unwrap();
    });
    x.enums.iter().for_each(|e| {
        p.add_enum(
            &e.name.text,
            e.variants
                .iter()
                .map(|x| EnumVariant {
                    name: x.name.text.clone(),
                    serialized_value: x.name.text.to_case(Case::UpperSnake),
                })
                .collect(),
        )
        .unwrap();
    });
    Ok(p)
}

fn make_type_ref(ty: &ast::Type) -> TypeRef {
    match ty {
        ast::Type::Type { identifier } => TypeRef::Named {
            name: identifier.text.clone(),
        },
        ast::Type::Optional { inner } => TypeRef::Optional {
            inner: Box::new(make_type_ref(inner)),
        },
        ast::Type::Array { inner } => TypeRef::Array {
            inner: Box::new(make_type_ref(inner)),
        },
        ast::Type::Generic { inner, parameters } => TypeRef::Generic {
            inner: Box::new(make_type_ref(inner)),
            parameters: parameters.iter().map(make_type_ref).collect(),
        },
    }
}

#[derive(Clone, Debug, Serialize)]
pub struct Protocol {
    pub types: BTreeMap<String, TypeDefinition>,
    pub methods: BTreeMap<String, Method>,
}

#[derive(Clone, Debug, Serialize)]
pub struct Method {
    pub name: String,
    pub parameters: Vec<(String, TypeRef)>,
    pub returns: Option<TypeRef>,
    pub serialization_strategy: SerializationStrategy,
}

#[derive(Debug, Clone, Serialize)]
pub enum SerializationStrategy {
    Multipart,
    Json,
}

impl Protocol {
    pub fn new() -> Self {
        Self {
            types: BTreeMap::new(),
            methods: BTreeMap::new(),
        }
    }

    pub fn add_model(&mut self, name: &str, fields: BTreeMap<String, TypeRef>) -> Result<(), ()> {
        if self
            .types
            .insert(name.into(), TypeDefinition::Object { fields })
            .is_some()
        {
            Err(())
        } else {
            Ok(())
        }
    }

    pub fn add_enum(&mut self, name: &str, variants: Vec<EnumVariant>) -> Result<(), ()> {
        if self
            .types
            .insert(name.into(), TypeDefinition::Enum { variants })
            .is_some()
        {
            Err(())
        } else {
            Ok(())
        }
    }

    pub fn add_scalar(&mut self, name: &str, serialized_type: &str) -> Result<(), ()> {
        if self
            .types
            .insert(
                name.into(),
                TypeDefinition::Scalar {
                    serialized_type: serialized_type.into(),
                },
            )
            .is_some()
        {
            Err(())
        } else {
            Ok(())
        }
    }

    pub fn add_method(
        &mut self,
        name: &str,
        parameters: Vec<(String, TypeRef)>,
        returns: Option<TypeRef>,
    ) -> Result<(), ()> {
        let serialization_strategy = if parameters.iter().any(|(_, ty)| self.has_upload_field(ty)) {
            SerializationStrategy::Multipart
        } else {
            SerializationStrategy::Json
        };

        if self
            .methods
            .insert(
                name.into(),
                Method {
                    name: name.into(),
                    parameters,
                    returns,
                    serialization_strategy,
                },
            )
            .is_some()
        {
            Err(())
        } else {
            Ok(())
        }
    }

    pub fn has_upload_field(&self, type_ref: &TypeRef) -> bool {
        match type_ref {
            TypeRef::Named { name } => match self.types.get(name) {
                Some(TypeDefinition::Object { fields }) => {
                    fields.iter().any(|(_, f)| self.has_upload_field(f))
                }
                _ => name == "Upload",
            },
            TypeRef::Optional { inner } => self.has_upload_field(&inner),
            TypeRef::Array { inner } => self.has_upload_field(&inner),
            TypeRef::Generic { inner, .. } => self.has_upload_field(&inner),
        }
    }
}

#[derive(Clone, Debug, Serialize)]
pub enum TypeDefinition {
    Scalar { serialized_type: String },
    Object { fields: BTreeMap<String, TypeRef> },
    Enum { variants: Vec<EnumVariant> },
}

#[derive(Clone, Debug, Serialize)]
pub struct EnumVariant {
    pub name: String,
    pub serialized_value: String,
}

#[derive(Debug, PartialEq, Clone, Serialize)]
pub enum TypeRef {
    Named {
        name: String,
    },
    Optional {
        inner: Box<TypeRef>,
    },
    Array {
        inner: Box<TypeRef>,
    },
    Generic {
        inner: Box<TypeRef>,
        parameters: Vec<TypeRef>,
    },
}
