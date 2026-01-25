use crate::generation::tera_helpers::{
    register_tera_filters, register_tera_testers, RenderContext,
};

use crate::protocol;
use crate::protocol::TypeDefinition;
use convert_case::{Case, Casing};
use serde::Serialize;
use std::fs;
use std::fs::File;
use std::io::Write;
use std::path::Path;
use tera::{Context, Error, Tera};

pub struct Config {
    pub out_dir: String,
    pub sdk_file: Option<String>,
}

#[derive(Debug)]
pub enum GeneratorError {
    TemplateError(String),
    FileError(String),
}

impl GeneratorError {
    fn from_tera(err: Error) -> Self {
        Self::TemplateError(err.to_string())
    }

    fn file_error(err: std::io::Error) -> Self {
        Self::FileError(err.to_string())
    }
}
pub fn generate_typescript(
    protocol: &protocol::Protocol,
    cfg: &Config,
) -> Result<(), GeneratorError> {
    let tera = setup_tera()?;

    if let Some(sdk_file) = &cfg.sdk_file {
        let mut f = create_output_file(&cfg.out_dir, sdk_file)?;
        render_sdk(&tera, protocol, &mut f)?;
        f.flush().map_err(GeneratorError::file_error)?;
    }

    Ok(())
}

fn setup_tera() -> Result<Tera, GeneratorError> {
    let mut tera = Tera::new("templates/**/*").map_err(GeneratorError::from_tera)?;

    register_tera_testers(&mut tera);
    register_tera_filters(&mut tera, type_ref_string);

    Ok(tera)
}

fn create_output_file(out_dir: &str, file_name: &str) -> Result<File, GeneratorError> {
    let dir = Path::new(out_dir);
    fs::create_dir_all(dir).map_err(GeneratorError::file_error)?;
    let path = dir.join(file_name);
    File::create(path).map_err(GeneratorError::file_error)
}

fn render_sdk(
    tera: &Tera,
    protocol: &protocol::Protocol,
    mut out: impl Write,
) -> Result<(), GeneratorError> {
    tera.render_to("typescript/preamble.tera", &Context::new(), &mut out)
        .map_err(GeneratorError::from_tera)?;

    fn new_ctx<T: Serialize>(rc: RenderContext<T>) -> Result<Context, GeneratorError> {
        Context::from_serialize(rc).map_err(GeneratorError::from_tera)
    }

    for type_def in protocol.types().values() {
        let (template, ctx) = match type_def {
            TypeDefinition::Scalar(s) => {
                ("typescript/scalar.tera", new_ctx(RenderContext::new(s))?)
            }
            TypeDefinition::Model(m) => ("typescript/model.tera", new_ctx(RenderContext::new(m))?),
            TypeDefinition::Enum(e) => ("typescript/enum.tera", new_ctx(RenderContext::new(e))?),
        };

        tera.render_to(template, &ctx, &mut out)
            .map_err(GeneratorError::from_tera)?;
    }

    for (_, method) in protocol.methods() {
        tera.render_to(
            "typescript/function.tera",
            &Context::from_serialize(RenderContext::new(method))
                .map_err(GeneratorError::from_tera)?,
            &mut out,
        )
        .map_err(GeneratorError::from_tera)?;
    }

    Ok(())
}

fn type_ref_string(type_ref: &protocol::TypeRef) -> String {
    match &type_ref {
        protocol::TypeRef::Named { name } => {
            if name == "string" {
                name.into()
            } else if name == "void" {
                name.into()
            } else {
                name.to_case(Case::UpperCamel)
            }
        }
        protocol::TypeRef::Optional { inner } => format!("{} | null", type_ref_string(inner)),
        protocol::TypeRef::Array { inner } => format!("{}[]", type_ref_string(inner)),
        protocol::TypeRef::Generic { inner, parameters } => {
            let generic_params = parameters
                .iter()
                .map(type_ref_string)
                .collect::<Vec<_>>()
                .join(", ");
            format!("{}<{}>", type_ref_string(inner), generic_params)
        }
    }
}
