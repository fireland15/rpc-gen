use super::*;
use std::collections::BTreeMap;

pub fn build_protocol(ast: &ast::ProtocolDefinition) -> Result<Protocol, Vec<ProtocolError>> {
    let mut builder = ProtocolBuilder::new();

    for model in &ast.models {
        builder.add_model(
            &model.name.text,
            model
                .fields
                .iter()
                .map(|f| (f.name.text.clone(), make_type_ref(&f.ty)))
                .collect(),
        );
    }

    for scalar in &ast.scalars {
        builder.add_scalar(&scalar.name.text, &scalar.serialized.text);
    }

    for method in &ast.methods {
        builder.add_method(
            &method.name.text,
            method
                .parameters
                .iter()
                .map(|p| (p.name.text.clone(), make_type_ref(&p.ty)))
                .collect(),
            method.return_ty.as_ref().map(make_type_ref),
        );
    }

    for e in &ast.enums {
        builder.add_enum(
            &e.name.text,
            e.variants
                .iter()
                .map(|v| EnumVariant {
                    name: v.name.text.clone(),
                    serialized_value: v.name.text.to_case(Case::UpperSnake),
                })
                .collect(),
        );
    }

    builder.build()
}

#[derive(Debug, Clone)]
pub enum ProtocolError {
    DuplicateType { name: String },
    DuplicateMethod { name: String },
}

#[derive(Default)]
pub struct ProtocolBuilder {
    types: BTreeMap<String, TypeDefinition>,
    methods: BTreeMap<String, Method>,
    errors: Vec<ProtocolError>,
}

impl ProtocolBuilder {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn add_model(&mut self, name: &str, fields: BTreeMap<String, TypeRef>) {
        if self.types.contains_key(name) {
            self.errors
                .push(ProtocolError::DuplicateType { name: name.into() });
            return;
        }

        let fields = fields
            .into_iter()
            .map(|(name, ty)| (name.clone(), FieldDefinition::new(name.clone(), ty)))
            .collect();

        self.types.insert(
            name.into(),
            TypeDefinition::Model(ModelDefinition {
                name: name.into(),
                fields,
            }),
        );
    }

    pub fn add_enum(&mut self, name: &str, variants: Vec<EnumVariant>) {
        if self.types.contains_key(name) {
            self.errors
                .push(ProtocolError::DuplicateType { name: name.into() });
            return;
        }

        self.types.insert(
            name.into(),
            TypeDefinition::Enum(EnumDefinition {
                name: name.into(),
                variants,
            }),
        );
    }

    pub fn add_scalar(&mut self, name: &str, serialized_type: &str) {
        if self.types.contains_key(name) {
            self.errors
                .push(ProtocolError::DuplicateType { name: name.into() });
            return;
        }

        self.types.insert(
            name.into(),
            TypeDefinition::Scalar(ScalarDefinition {
                name: name.into(),
                serialized_type: serialized_type.into(),
            }),
        );
    }

    pub fn add_method(
        &mut self,
        name: &str,
        parameters: Vec<(String, TypeRef)>,
        returns: Option<TypeRef>,
    ) {
        if self.methods.contains_key(name) {
            self.errors
                .push(ProtocolError::DuplicateMethod { name: name.into() });
            return;
        }

        let parameters = parameters
            .into_iter()
            .map(|(name, ty)| Parameter::new(name, ty))
            .collect();

        self.methods.insert(
            name.into(),
            Method {
                name: name.into(),
                path: name.to_case(Case::Snake),
                parameters,
                returns,
            },
        );
    }

    pub fn build(self) -> Result<Protocol, Vec<ProtocolError>> {
        if self.errors.is_empty() {
            Ok(Protocol {
                schema: Schema {
                    types: self.types,
                    methods: self.methods,
                },
            })
        } else {
            Err(self.errors)
        }
    }
}
