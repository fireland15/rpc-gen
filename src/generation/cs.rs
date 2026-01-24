use crate::generation::protocol;
use convert_case::{Case, Casing};
use serde::Serialize;
use std::collections::BTreeMap;
use std::fs;
use std::fs::File;
use std::path::Path;
use serde_json::Value;
use tera::{Context, Tera};

pub struct Config {
    pub out_dir: String,
    pub namespace: String,
    pub model_dir: String,
    pub model_namespace: String,
    pub scalar_map: BTreeMap<String, String>,
}

pub fn generate_csharp(protocol: &protocol::Protocol, cfg: &Config) -> Result<(), ()> {
    let mut tera = Tera::new("templates/**/*").unwrap();

    tera.register_tester("none", |value: Option<&Value>, _args: &[Value]| -> tera::Result<bool> {
        match value {
            None => Ok(true),
            Some(x) => Ok(x.is_null())
        }
    });

    let file_path = Path::new(cfg.model_dir.as_str());
    fs::create_dir_all(file_path).unwrap();

    for (name, ty) in protocol.types.iter() {
        match ty {
            protocol::TypeDefinition::Enum(enum_def) => {
                render_enum(&mut tera, enum_def, cfg)?;
            }
            protocol::TypeDefinition::Model(model_def) => {
                render_model(&mut tera, model_def, cfg)?;
            }
            _ => {}
        }
    }

    for (_, method) in protocol.methods.iter() {
        let name = format!("{}Parameters", &method.name);
        if method.parameters.len() == 0 {
            continue;
        }
        let mut params_class = CsModel::new(&name);
        for (name, ty) in method.parameters.iter() {
            params_class.properties.push(CsProperty {
                name: name.to_case(Case::UpperCamel),
                ty: type_ref_string(&ty, &cfg.scalar_map),
                required: !matches!(ty, protocol::TypeRef::Optional { .. }),
                json_property_name: name.to_case(Case::Camel),
            })
        }
        render_cs_model(&mut tera, &params_class, cfg)?;
    }

    render_methods(&mut tera, protocol, cfg)?;

    Ok(())
}

fn render_enum(
    tera: &mut Tera,
    enum_def: &protocol::EnumDefinition,
    cfg: &Config,
) -> Result<(), ()> {
    let file_path = Path::new(cfg.model_dir.as_str());
    let file_name = format!("{}.cs", &enum_def.name);
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/common/enum.tera",
        &Context::from_serialize(CsFile {
            namespace: cfg.model_namespace.clone(),
            model: enum_def,
        })
        .unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

fn render_model(
    tera: &mut Tera,
    model_def: &protocol::ModelDefinition,
    cfg: &Config,
) -> Result<(), ()> {
    let file_path = Path::new(cfg.model_dir.as_str());
    let file_name = format!("{}.cs", &model_def.name);
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/common/model.tera",
        &Context::from_serialize(CsFile {
            namespace: cfg.model_namespace.clone(),
            model: CsModel::from_protocol(model_def, cfg),
        })
        .unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

fn render_cs_model(
    tera: &mut Tera,
    model: &CsModel,
    cfg: &Config,
) -> Result<(), ()> {
    let file_path = Path::new(cfg.model_dir.as_str());
    let file_name = format!("{}.cs", &model.name);
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/common/model.tera",
        &Context::from_serialize(CsFile {
            namespace: cfg.model_namespace.clone(),
            model,
        })
        .unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

fn render_methods(tera: &mut Tera, protocol: &protocol::Protocol, cfg: &Config) -> Result<(), ()> {
    let file_path = Path::new(cfg.out_dir.as_str());
    let file_name = "Protocol.cs";
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/minimal_api/protocol.tera",
        &Context::from_serialize(CsFile {
            namespace: cfg.namespace.clone(),
            model: protocol
                .methods
                .iter()
                .map(|(n, m)| CsMethod::new(&m, cfg))
                .collect::<Vec<CsMethod>>(),
        })
        .unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

#[derive(Debug, Serialize)]
pub struct CsFile<T> {
    pub namespace: String,
    pub model: T,
}

#[derive(Debug, Serialize)]
pub struct CsModel {
    name: String,
    properties: Vec<CsProperty>,
}

#[derive(Debug, Serialize)]
pub struct CsProperty {
    name: String,
    ty: String,
    required: bool,
    json_property_name: String,
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
    name: String,
    path: String,
    parameters: Vec<MethodParameter>,
    return_ty: Option<String>,
    serialization_strategy: protocol::SerializationStrategy,
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
            serialization_strategy: method.serialization_strategy.clone(),
        }
    }
}

fn type_ref_string(type_ref: &protocol::TypeRef, scalar_map: &BTreeMap<String, String>) -> String {
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
