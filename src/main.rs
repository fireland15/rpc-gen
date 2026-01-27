use crate::generation::cs::generate_csharp;
use crate::generation::{cs, typescript};
use crate::protocol::build_protocol;
use generation::typescript::generate_typescript;
use parsing::parser;
use std::collections::BTreeMap;
use std::sync::Arc;
use std::{fs::File, io::Read};

mod generation;
mod parsing;
mod protocol;

fn main() {
    let mut file = File::open("journal.rpc").expect("Failed to open file");
    let mut contents = String::new();
    file.read_to_string(&mut contents)
        .expect("problem reading file");

    let ast = parser::parse(&contents).expect("problem parsing");
    let p = build_protocol(&ast).unwrap();
    generate_typescript(
        &p,
        &typescript::Config {
            out_dir: "out/ts_sdk".into(),
            sdk_file: Some("sdk.ts".into()),
        },
    )
    .unwrap();

    let mut cs_scalar_map = BTreeMap::new();
    cs_scalar_map.insert("string".into(), "string".into());
    cs_scalar_map.insert("date".into(), "DateTime".into());
    cs_scalar_map.insert("int".into(), "int".into());
    cs_scalar_map.insert("uuid".into(), "Guid".into());

    let cs_cfg = cs::Config {
        out_dir: "out/cs".into(),
        namespace: "Test.Namespace".into(),
        model_dir: "out/cs/models".into(),
        model_namespace: "Test.Namespace.Models".to_string(),
        interfaces_dir: "out/cs/Abstractions".into(),
        interfaces_namespace: "Test.Namespace.Abstractions".to_string(),
        scalar_map: Arc::new(cs_scalar_map),
    };

    generate_csharp(&p, &cs_cfg).unwrap();
}
