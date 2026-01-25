use crate::protocol;
use convert_case::{Case, Casing};
use serde::Serialize;
use serde_json::Value;
use std::collections::HashMap;
use tera::Tera;

pub fn register_tera_testers(tera: &mut Tera) {
    tera.register_tester(
        "none",
        |value: Option<&Value>, _args: &[Value]| -> tera::Result<bool> {
            match value {
                None => Ok(true),
                Some(x) => Ok(x.is_null()),
            }
        },
    );
}

pub fn register_tera_filters<F>(tera: &mut Tera, resolve_type_ref: F)
where
    F: Fn(&protocol::TypeRef) -> String + Send + Sync + 'static,
{
    tera.register_filter(
        "to_case",
        |v: &Value, args: &HashMap<String, Value>| -> tera::Result<Value> {
            match v {
                Value::String(s) => match args.get("case") {
                    None => Ok(v.clone()),
                    Some(Value::String(case)) => match case.as_str() {
                        "snake" => Ok(s.to_case(Case::Snake).into()),
                        "upper_snake" => Ok(s.to_case(Case::UpperSnake).into()),
                        "camel" => Ok(s.to_case(Case::Camel).into()),
                        "upper_camel" => Ok(s.to_case(Case::UpperCamel).into()),
                        _ => Err(tera::Error::msg("unsupported case")),
                    },
                    _ => Err(tera::Error::msg("case arg must be a string")),
                },
                _ => Err(tera::Error::msg("value is not a string")),
            }
        },
    );

    tera.register_filter(
        "resolve_type_ref",
        move |v: &Value, _: &HashMap<String, Value>| -> tera::Result<Value> {
            let c = serde_json::from_value::<protocol::TypeRef>(v.clone())
                .map_err(|e| tera::Error::chain("TypeRef", e))?;
            Ok(resolve_type_ref(&c).into())
        },
    );
}

#[derive(Debug, Serialize)]
pub struct RenderContext<'a, T> {
    item: &'a T,
}

impl<'a, T> RenderContext<'a, T> {
    pub fn new(item: &'a T) -> Self {
        Self { item }
    }
}
