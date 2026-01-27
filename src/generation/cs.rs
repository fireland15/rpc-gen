use std::collections::BTreeMap;

use crate::generation::error::GeneratorError;
use crate::generation::tera_helpers::{register_tera_filters, register_tera_testers};
use crate::protocol;
use convert_case::{Case, Casing};
use serde::Serialize;
use std::fs;
use std::fs::File;
use std::io::Write;
use std::path::Path;
use std::sync::Arc;
use tera::{Context, Tera};

pub struct Config {
    pub out_dir: String,
    pub namespace: String,
    pub model_dir: String,
    pub model_namespace: String,
    pub interfaces_dir: String,
    pub interfaces_namespace: String,
    pub scalar_map: Arc<BTreeMap<String, String>>,
}

pub fn generate_csharp(protocol: &protocol::Protocol, cfg: &Config) -> Result<(), GeneratorError> {
    let tera = setup_tera(cfg)?;

    let file_path = Path::new(cfg.model_dir.as_str());
    fs::create_dir_all(file_path).unwrap();

    for (_, ty) in protocol.types().iter() {
        match ty {
            protocol::TypeDefinition::Enum(enum_def) => {
                render_enum(&tera, enum_def, cfg)?;
            }
            protocol::TypeDefinition::Model(model_def) => {
                render_model(&tera, model_def, cfg)?;
            }
            _ => {}
        }
    }

    // for (_, method) in protocol.schema().methods.iter() {
    //     let name = format!("{}Parameters", &method.name);
    //     if method.parameters.len() == 0 {
    //         continue;
    //     }
    //     let mut params_class = CsModel::new(&name);
    //     for (_, param) in method.parameters.iter() {
    //         params_class.properties.push(CsProperty {
    //             name: name.to_case(Case::UpperCamel),
    //             ty: type_ref_string(&param.ty, &cfg.scalar_map),
    //             required: !matches!(param.ty, protocol::TypeRef::Optional { .. }),
    //             json_property_name: name.to_case(Case::Camel),
    //         })
    //     }
    //     render_cs_model(&mut tera, &params_class, cfg)?;
    // }

    render_interfaces(&tera, protocol, cfg)?;
    render_methods(&tera, protocol, cfg)?;

    Ok(())
}

fn setup_tera(cfg: &Config) -> Result<Tera, GeneratorError> {
    let mut tera = Tera::new("templates/**/*").map_err(GeneratorError::from_tera)?;
    register_tera_testers(&mut tera);

    let scalar_map = Arc::clone(&cfg.scalar_map);
    register_tera_filters(&mut tera, move |type_ref| {
        type_ref_string(type_ref, &scalar_map)
    });
    Ok(tera)
}

fn render_enum(
    tera: &Tera,
    enum_def: &protocol::EnumDefinition,
    cfg: &Config,
) -> Result<(), GeneratorError> {
    let file_path = Path::new(cfg.model_dir.as_str());
    let file_name = format!("{}.cs", &enum_def.name);
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/common/enum.tera",
        &Context::from_serialize(CsFile::new(cfg.model_namespace.as_str(), &enum_def)).unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

fn render_model(
    tera: &Tera,
    model_def: &protocol::ModelDefinition,
    cfg: &Config,
) -> Result<(), GeneratorError> {
    let file_path = Path::new(cfg.model_dir.as_str());
    let file_name = format!("{}.cs", &model_def.name);
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/common/model.tera",
        &Context::from_serialize(CsFile::new(cfg.model_namespace.as_str(), &model_def)).unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

fn render_methods(
    tera: &Tera,
    protocol: &protocol::Protocol,
    cfg: &Config,
) -> Result<(), GeneratorError> {
    let file_path = Path::new(cfg.out_dir.as_str());
    let file_name = "Protocol.cs";
    let file_path = file_path.join(&file_name);
    let f = File::create(file_path).unwrap();

    tera.render_to(
        "cs/minimal_api/protocol.tera",
        &Context::from_serialize(CsFile::new(cfg.model_namespace.as_str(), protocol.schema())).unwrap(),
        f,
    )
    .unwrap();

    Ok(())
}

fn render_interfaces(
    tera: &Tera,
    protocol: &protocol::Protocol,
    cfg: &Config,
) -> Result<(), GeneratorError> {
    let file_path = Path::new(cfg.interfaces_dir.as_str());
    fs::create_dir_all(file_path).unwrap();

    for (_, method) in protocol.methods().iter() {
        let file_name = format!("I{}Handler.cs", &method.name);
        let file_path = file_path.join(&file_name);
        let mut f = File::create(file_path).unwrap();
        tera.render_to(
            "cs/common/interface.tera",
            &Context::from_serialize(CsFile::new(cfg.model_namespace.as_str(), &method)).unwrap(),
            &f,
        )
        .unwrap();
        f.flush().unwrap();
    }

    Ok(())
}

fn type_ref_string(type_ref: &protocol::TypeRef, scalar_map: &BTreeMap<String, String>) -> String {
    match &type_ref {
        protocol::TypeRef::Named { name } => {
            if scalar_map.contains_key(name) {
                scalar_map[name].clone()
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

#[derive(Debug, Serialize)]
struct CsFile<'a, T> {
    namespace: String,
    item: &'a T,
}

impl<'a, T> CsFile<'a, T> {
    fn new(namespace: &str, item: &'a T) -> Self {
        Self {
            namespace: namespace.to_string(),
            item,
        }
    }
}
