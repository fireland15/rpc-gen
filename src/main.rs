use std::{fs::File, io::Read};

use parsing::parser;

mod parsing;

fn main() {
    let mut file = File::open("journal.rpc").expect("Failed to open file");
    let mut contents = String::new();
    file.read_to_string(&mut contents)
        .expect("problem reading file");

    let ast = parser::parse(&contents).expect("problem parsing");

    dbg!(ast);
}
