use crate::generation::cs;
use crate::generation::cs::generate_csharp;
use generation::{protocol::build_protocol, typescript::generate_typescript};
use parsing::parser;
use std::collections::BTreeMap;
use std::{fs::File, io::Read};

mod generation;
mod parsing;

fn main() {
    let mut file = File::open("journal.rpc").expect("Failed to open file");
    let mut contents = String::new();
    file.read_to_string(&mut contents)
        .expect("problem reading file");

    let ast = parser::parse(&contents).expect("problem parsing");
    let p = build_protocol(&ast).unwrap();
    generate_typescript(&p);

    let mut cs_scalar_map = BTreeMap::new();
    cs_scalar_map.insert("string".into(), "string".into());
    cs_scalar_map.insert("date".into(), "DateTime".into());
    cs_scalar_map.insert("int".into(), "int".into());
    cs_scalar_map.insert("uuid".into(), "Guid".into());

    generate_csharp(
        &p,
        &cs::Config {
            out_dir: "out".into(),
            namespace: "Test.Namespace".into(),
            model_dir: "out/models".into(),
            model_namespace: "Test.Namespace.Models".to_string(),
            interfaces_dir: "out/Abstractions".into(),
            interfaces_namespace: "Test.Namespace.Abstractions".to_string(),
            scalar_map: cs_scalar_map,
        },
    )
    .unwrap();
}
