use crate::generation::protocol::TypeDefinition;
use std::fs::File;

use super::protocol::{Protocol, SerializationStrategy, TypeRef};
use convert_case::{Case, Casing};
use serde::Serialize;
use tera::{Context, Tera, Value};

pub fn generate_typescript(protocol: &Protocol) {
    let mut tera = Tera::new("templates/**/*").unwrap();
    
    pub fn none(value: Option<&Value>, _args: &[Value]) -> tera::Result<bool> {
        Ok(value.unwrap().is_null())
    }
    tera.register_tester("null", none);

    let f = File::create("api.ts").unwrap();

    let mut ts_client = TsClientDefinition {
        scalars: Vec::new(),
        types: Vec::new(),
        methods: Vec::new(),
    };

    protocol.types.iter().for_each(|(name, ty)| match ty {
        TypeDefinition::Scalar { serialized_type } => {
            let scalar_def = TsScalarDefinition {
                name: name.to_case(Case::UpperCamel),
                ty: serialized_type.clone(),
                brand: name.to_case(Case::Camel),
            };
            ts_client.scalars.push(scalar_def);
        }
        TypeDefinition::Object { fields } => {
            let type_def = TsTypeDefinition {
                name: name.clone(),
                fields: fields
                    .iter()
                    .map(|(name, ty)| TsFieldDefinition {
                        name: name.clone(),
                        ty: type_ref_string(ty),
                    })
                    .collect(),
            };
            ts_client.types.push(type_def);
        }
    });

    protocol.methods.iter().for_each(|(_, method)| {
        let m = TsMethodDefinition::new(
            &method.name,
            method
                .parameters
                .iter()
                .map(|(x, y)| TsMethodParameter {
                    name: x.clone(),
                    ty: type_ref_string(&y),
                    client_side: false,
                })
                .collect(),
            &method.returns,
            method.serialization_strategy.clone(),
        );
        ts_client.methods.push(m);
    });

    tera.render_to("typescript/ts_client.template", &Context::from_serialize(ts_client).unwrap(), f)
        .unwrap();
}

#[derive(Serialize)]
struct TsClientDefinition {
    scalars: Vec<TsScalarDefinition>,
    types: Vec<TsTypeDefinition>,
    methods: Vec<TsMethodDefinition>,
}

#[derive(Serialize)]
struct TsScalarDefinition {
    name: String,
    ty: String,
    brand: String,
}

#[derive(Serialize)]
struct TsTypeDefinition {
    name: String,
    fields: Vec<TsFieldDefinition>,
}

#[derive(Debug, Serialize)]
pub struct TsMethodDefinition {
    name: String,
    parameters: Vec<TsMethodParameter>,
    returns: Option<String>,
    serialization_strategy: String,
}

#[derive(Debug, Serialize)]
pub struct TsMethodParameter {
    name: String,
    ty: String,

    /// Controls whether the method parameter is sent in the request body
    client_side: bool,
}

impl TsMethodDefinition {
    fn new(
        name: &str,
        parameters: Vec<TsMethodParameter>,
        type_ref: &Option<TypeRef>,
        serialization_strategy: SerializationStrategy,
    ) -> Self {
        Self {
            name: name.to_case(convert_case::Case::Camel),
            parameters,
            returns: type_ref.as_ref().map(type_ref_string),
            serialization_strategy: match serialization_strategy {
                SerializationStrategy::Multipart => "multipart".into(),
                SerializationStrategy::Json => "json".into(),
            },
        }
    }
}

#[derive(Serialize)]
struct TsFieldDefinition {
    name: String,
    ty: String,
}

fn type_ref_string(type_ref: &TypeRef) -> String {
    match &type_ref {
        TypeRef::Named { name } => {
            if name == "Upload" {
                "File".into()
            } else if name == "string" {
                name.into()
            } else if name == "void" {
                name.into()
            } else {
                name.to_case(Case::UpperCamel)
            }
        }
        TypeRef::Optional { inner } => format!("{} | null", type_ref_string(inner)),
        TypeRef::Array { inner } => format!("{}[]", type_ref_string(inner)),
        TypeRef::Generic { inner, parameters } => {
            let generic_params = parameters
                .iter()
                .map(type_ref_string)
                .collect::<Vec<_>>()
                .join(", ");
            format!("{}<{}>", type_ref_string(inner), generic_params)
        }
    }
}
