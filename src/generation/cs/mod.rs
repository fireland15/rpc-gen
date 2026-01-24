mod model;

pub(crate) use crate::generation::cs::model::Config;
use crate::generation::cs::model::{type_ref_string, CsFile, CsMethod, CsModel, CsProperty};
use crate::generation::protocol;
use crate::generation::tera_helpers::register_tera_testers;
use convert_case::{Case, Casing};
use std::fs;
use std::fs::File;
use std::path::Path;
use tera::{Context, Tera};

pub fn generate_csharp(protocol: &protocol::Protocol, cfg: &Config) -> Result<(), ()> {
    let mut tera = Tera::new("templates/**/*").unwrap();

    register_tera_testers(&mut tera);

    let file_path = Path::new(cfg.model_dir.as_str());
    fs::create_dir_all(file_path).unwrap();

    for (_, ty) in protocol.types.iter() {
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

    render_interfaces(&mut tera, protocol, cfg)?;
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

fn render_cs_model(tera: &mut Tera, model: &CsModel, cfg: &Config) -> Result<(), ()> {
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

fn render_interfaces(
    tera: &mut Tera,
    protocol: &protocol::Protocol,
    cfg: &Config,
) -> Result<(), ()> {
    let file_path = Path::new(cfg.interfaces_dir.as_str());
    fs::create_dir_all(file_path).unwrap();

    let mut it = protocol.methods.iter().map(|(n, m)| CsMethod::new(&m, cfg));

    for method in it {
        let file_name = format!("I{}Handler.cs", &method.name);
        let file_path = file_path.join(&file_name);
        let f = File::create(file_path).unwrap();
        tera.render_to(
            "cs/common/interface.tera",
            &Context::from_serialize(CsFile {
                namespace: cfg.interfaces_namespace.clone(),
                model: method,
            })
            .unwrap(),
            f,
        )
        .unwrap();
    }

    Ok(())
}
