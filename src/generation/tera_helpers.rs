use serde_json::Value;
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
