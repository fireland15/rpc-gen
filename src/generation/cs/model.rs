use crate::generation::protocol;
use convert_case::{Case, Casing};
use serde::Serialize;
use std::collections::BTreeMap;

pub struct Config {
    pub out_dir: String,
    pub namespace: String,
    pub model_dir: String,
    pub model_namespace: String,
    pub interfaces_dir: String,
    pub interfaces_namespace: String,
    pub scalar_map: BTreeMap<String, String>,
}

#[derive(Debug, Serialize)]
pub struct CsFile<T> {
    pub namespace: String,
    pub model: T,
}

#[derive(Debug, Serialize)]
pub struct CsModel {
    pub name: String,
    pub properties: Vec<CsProperty>,
}

#[derive(Debug, Serialize)]
pub struct CsProperty {
    pub name: String,
    pub ty: String,
    pub required: bool,
    pub json_property_name: String,
}

impl CsModel {
    pub fn new(name: &str) -> Self {
        Self {
            name: name.into(),
            properties: vec![],
        }
    }

    pub fn from_protocol(model_def: &protocol::ModelDefinition, cfg: &Config) -> Self {
        Self {
            name: model_def.name.to_case(Case::UpperCamel),
            properties: model_def
                .fields
                .iter()
                .map(|f| CsProperty {
                    name: f.name.to_case(Case::UpperCamel),
                    ty: type_ref_string(&f.ty, &cfg.scalar_map),
                    required: !matches!(f.ty, protocol::TypeRef::Optional { .. }),
                    json_property_name: f.name.to_case(Case::Camel),
                })
                .collect(),
        }
    }
}

#[derive(Debug, Serialize)]
pub struct CsMethod {
    pub name: String,
    pub path: String,
    pub parameters: Vec<MethodParameter>,
    pub return_ty: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct MethodParameter {
    // The name of the method parameter used in the request object.
    name: String,
    // the name of the method parameter used as an argument.
    argument_name: String,
    ty: String,
}

impl CsMethod {
    pub fn new(method: &protocol::Method, cfg: &Config) -> Self {
        Self {
            name: method.name.to_case(Case::UpperCamel),
            path: method.path.clone(),
            parameters: method
                .parameters
                .iter()
                .map(|(n, ty)| MethodParameter {
                    name: n.to_case(Case::UpperCamel),
                    argument_name: n.to_case(Case::Camel),
                    ty: type_ref_string(ty, &cfg.scalar_map),
                })
                .collect(),
            return_ty: method
                .returns
                .as_ref()
                .map(|ty| type_ref_string(&ty, &cfg.scalar_map)),
        }
    }
}

pub fn type_ref_string(
    type_ref: &protocol::TypeRef,
    scalar_map: &BTreeMap<String, String>,
) -> String {
    match &type_ref {
        protocol::TypeRef::Named { name } => {
            if scalar_map.contains_key(name) {
                return scalar_map[name].clone();
            } else if name == "void" {
                name.into()
            } else {
                name.to_case(Case::UpperCamel)
            }
        }
        protocol::TypeRef::Optional { inner } => format!("{}?", type_ref_string(inner, scalar_map)),
        protocol::TypeRef::Array { inner } => format!("{}[]", type_ref_string(inner, scalar_map)),
        protocol::TypeRef::Generic { inner, parameters } => {
            let generic_params = parameters
                .iter()
                .map(|x| type_ref_string(x, scalar_map))
                .collect::<Vec<_>>()
                .join(", ");
            format!("{}<{}>", type_ref_string(inner, scalar_map), generic_params)
        }
    }
}
