use std::{fs::File, io::Read};

use generation::{protocol::build_protocol, typescript::generate_typescript};
use parsing::parser;

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
}
