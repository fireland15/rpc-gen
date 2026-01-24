use crate::generation::protocol::{ModelDefinition, TypeDefinition};
use std::fs::File;

use super::protocol::{Protocol, TypeRef};
use crate::generation::protocol;
use crate::generation::tera_helpers::register_tera_testers;
use convert_case::{Case, Casing};
use serde::Serialize;
use tera::{Context, Tera, Value};

pub fn generate_typescript(protocol: &Protocol) {
    let mut tera = Tera::new("templates/**/*").unwrap();
    register_tera_testers(&mut tera);
    pub fn none(value: Option<&Value>, _args: &[Value]) -> tera::Result<bool> {
        Ok(value.unwrap().is_null())
    }
    tera.register_tester("null", none);

    let f = File::create("api.ts").unwrap();

    let mut ts_client = TsClientDefinition {
        scalars: Vec::new(),
        enums: Vec::new(),
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
        TypeDefinition::Model(ModelDefinition { fields, .. }) => {
            let type_def = TsTypeDefinition {
                name: name.clone(),
                fields: fields
                    .iter()
                    .map(|f| TsFieldDefinition {
                        name: f.name.to_case(Case::Camel),
                        ty: type_ref_string(&f.ty),
                    })
                    .collect(),
            };
            ts_client.types.push(type_def);
        }
        TypeDefinition::Enum(protocol::EnumDefinition { variants, .. }) => {
            let enum_def = TsEnumDefinition {
                name: name.clone(),
                variants: variants
                    .iter()
                    .map(|v| TsEnumVariant {
                        name: v.name.clone(),
                        value: v.serialized_value.clone(),
                    })
                    .collect(),
            };
            ts_client.enums.push(enum_def);
        }
    });

    protocol.methods.iter().for_each(|(_, method)| {
        let m = TsMethodDefinition::new(
            &method.name,
            &method.path,
            method
                .parameters
                .iter()
                .map(|(x, y)| TsMethodParameter {
                    name: x.to_case(Case::Camel),
                    ty: type_ref_string(&y),
                })
                .collect(),
            &method.returns,
        );
        ts_client.methods.push(m);
    });

    tera.render_to(
        "typescript/ts_client.template",
        &Context::from_serialize(ts_client).unwrap(),
        f,
    )
    .unwrap();
}

#[derive(Serialize)]
struct TsClientDefinition {
    scalars: Vec<TsScalarDefinition>,
    types: Vec<TsTypeDefinition>,
    methods: Vec<TsMethodDefinition>,
    enums: Vec<TsEnumDefinition>,
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

#[derive(Serialize)]
struct TsEnumDefinition {
    name: String,
    variants: Vec<TsEnumVariant>,
}

#[derive(Serialize)]
struct TsEnumVariant {
    name: String,
    value: String,
}

#[derive(Debug, Serialize)]
pub struct TsMethodDefinition {
    name: String,
    path: String,
    parameters: Vec<TsMethodParameter>,
    returns: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct TsMethodParameter {
    name: String,
    ty: String,
}

impl TsMethodDefinition {
    fn new(
        name: &str,
        path: &str,
        parameters: Vec<TsMethodParameter>,
        type_ref: &Option<TypeRef>,
    ) -> Self {
        Self {
            name: name.to_case(Case::Camel),
            path: path.into(),
            parameters,
            returns: type_ref.as_ref().map(type_ref_string),
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
            if name == "string" {
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
