use tera::Tera;

use super::protocol::Protocol;

pub fn generate_go_echo_server(protocol: &Protocol) {
    let mut tera = Tera::default();
    tera.add_template_file("../../templates/typescript/ts_sdk.tera", Some("ts-client"))
        .unwrap();
}
