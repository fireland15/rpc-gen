mod builder;

use crate::parsing::ast;
use convert_case::{Case, Casing};
use serde::{Deserialize, Serialize};
use std::collections::BTreeMap;

pub use builder::build_protocol;

fn make_type_ref(ty: &ast::Type) -> TypeRef {
    match ty {
        ast::Type::Type { identifier } => TypeRef::Named {
            name: identifier.text.clone(),
        },
        ast::Type::Optional { inner } => TypeRef::Optional {
            inner: Box::new(make_type_ref(inner)),
        },
        ast::Type::Array { inner } => TypeRef::Array {
            inner: Box::new(make_type_ref(inner)),
        },
        ast::Type::Generic { inner, parameters } => TypeRef::Generic {
            inner: Box::new(make_type_ref(inner)),
            parameters: parameters.iter().map(make_type_ref).collect(),
        },
    }
}

/// Represents a complete RPC protocol definition.
///
/// A protocol consists of:
/// - A set of named type definitions (scalars, models, enums)
/// - A set of RPC methods callable over HTTP
///
/// This is the central, fully-resolved representation used by
/// downstream tooling (e.g. code generation or schema export).
#[derive(Clone, Debug, Serialize)]
pub struct Protocol {
    schema: Schema,
}

impl Protocol {
    pub fn schema(&self) -> &Schema {
        &self.schema
    }

    pub fn types(&self) -> &BTreeMap<String, TypeDefinition> {
        &self.schema.types
    }

    pub fn methods(&self) -> &BTreeMap<String, Method> {
        &self.schema.methods
    }
}

#[derive(Clone, Debug, Serialize)]
pub struct Schema {
    /// All type definitions used by the protocol, indexed by name.
    ///
    /// This includes scalars, models, and enums.
    pub types: BTreeMap<String, TypeDefinition>,

    /// RPC methods exposed by the protocol, indexed by method name.
    pub methods: BTreeMap<String, Method>,
}

impl Schema {
    /// Returns an iterator over scalar type definitions in the protocol
    pub fn scalars(&self) -> impl Iterator<Item = &ScalarDefinition> {
        self.types.iter().filter_map(|(_, x)| {
            if let TypeDefinition::Scalar(s) = x {
                return Some(s);
            }
            return None;
        })
    }

    /// Returns an iterator over enum definitions in the protocol
    pub fn enums(&self) -> impl Iterator<Item = &EnumDefinition> {
        self.types.iter().filter_map(|(_, x)| {
            if let TypeDefinition::Enum(e) = x {
                Some(e)
            } else {
                None
            }
        })
    }

    /// Returns an iterator over the model definitions in the protocol.
    pub fn models(&self) -> impl Iterator<Item = &ModelDefinition> {
        self.types.iter().filter_map(|(_, x)| {
            if let TypeDefinition::Model(m) = x {
                Some(m)
            } else {
                None
            }
        })
    }
}

/// Represents a single RPC method.
///
/// Each method maps to an HTTP request of the form:
/// `POST {BASE_URL}/{path}`
///
/// Parameters are serialized into the JSON request body,
/// and return values are deserialized from a successful (2xx)
/// JSON response.
#[derive(Clone, Debug, Serialize)]
pub struct Method {
    /// The logical name of the method as defined in the protocol.
    pub name: String,

    /// The URL path segment derived from the method name.
    ///
    /// This is typically a snake_case transformation of `name`.
    pub path: String,

    /// Parameters accepted by the method.
    ///
    /// These are serialized into the request body as JSON.
    pub parameters: Vec<Parameter>,

    /// The return type of the method, if any.
    ///
    /// If present, the response body of a successful request
    /// is deserialized as this type.
    pub returns: Option<TypeRef>,
}

/// Represents a single parameter to an RPC method.
#[derive(Clone, Debug, Serialize)]
pub struct Parameter {
    /// The logical name of the parameter as defined in the protocol.
    pub name: String,

    /// The name used when serializing this parameter to JSON.
    ///
    /// By convention this is a snake_case version of `name`.
    pub serialized_name: String,

    /// The parameter's type.
    pub ty: TypeRef,
}

impl Parameter {
    fn new(name: String, ty: TypeRef) -> Self {
        Self {
            serialized_name: name.to_case(Case::Snake),
            name,
            ty,
        }
    }
}

/// A named type definition in the RPC protocol.
///
/// Every type referenced by a [`TypeRef`] must correspond to one
/// of these definitions.
#[derive(Clone, Debug, Serialize)]
pub enum TypeDefinition {
    /// A scalar type that maps directly to a JSON primitive.
    Scalar(ScalarDefinition),

    /// A structured object type with named fields.
    Model(ModelDefinition),

    /// A closed set of named values.
    Enum(EnumDefinition),
}

/// Definition of a scalar type in the RPC protocol.
///
/// Scalars represent atomic values that map directly to
/// JSON primitives (e.g. string, number, boolean).
#[derive(Clone, Debug, Serialize)]
pub struct ScalarDefinition {
    /// The protocol-level name of the scalar.
    pub name: String,

    /// The JSON type this scalar serializes to.
    ///
    /// Examples: `"string"`, `"number"`, `"boolean"`.
    pub serialized_type: String,
}

/// Definition of a structured model (object) type.
///
/// Models consist of a fixed set of named fields, each with
/// an associated type.
#[derive(Clone, Debug, Serialize)]
pub struct ModelDefinition {
    /// The name of the model type.
    pub name: String,

    /// Fields belonging to this model.
    pub fields: BTreeMap<String, FieldDefinition>,
}

/// Definition of a single field within a model.
#[derive(Clone, Debug, Serialize)]
pub struct FieldDefinition {
    /// The logical name of the field in the protocol.
    pub name: String,

    /// The name used when serializing this field to JSON.
    ///
    /// Typically, a snake_case version of `name`.
    pub serialized_name: String,

    /// The field's type.
    pub ty: TypeRef,
}

impl FieldDefinition {
    fn new(name: String, ty: TypeRef) -> Self {
        Self {
            serialized_name: name.to_case(Case::Snake),
            name,
            ty,
        }
    }
}

/// Definition of an enum type in the RPC protocol.
///
/// Enum values are serialized as strings.
#[derive(Clone, Debug, Serialize)]
pub struct EnumDefinition {
    /// The name of the enum type.
    pub name: String,

    /// All possible variants of the enum.
    pub variants: Vec<EnumVariant>,
}

/// Definition of a single enum variant.
#[derive(Clone, Debug, Serialize)]
pub struct EnumVariant {
    /// The logical name of the variant.
    pub name: String,

    /// The value used when serializing this variant to JSON.
    pub serialized_value: String,
}

/// A reference to a type used in parameters, fields, and return values.
///
/// `TypeRef` is recursive and can represent:
/// - Named types (scalars, models, enums)
/// - Optional values
/// - Arrays
/// - Generic types with parameters
///
/// This enum describes *usage*, not definition. All named references
/// must correspond to a [`TypeDefinition`] in the protocol.
#[derive(Debug, PartialEq, Clone, Serialize, Deserialize)]
pub enum TypeRef {
    /// A reference to a named type.
    ///
    /// The name must correspond to a registered type definition.
    Named { name: String },

    /// An optional (nullable) value.
    Optional { inner: Box<TypeRef> },

    /// A list of values of the same type.
    Array { inner: Box<TypeRef> },

    /// A generic type application.
    ///
    /// For example: `Result<T, E>` or `Map<K, V>`.
    Generic {
        /// The generic type being applied.
        inner: Box<TypeRef>,

        /// Type parameters supplied to the generic.
        parameters: Vec<TypeRef>,
    },
}
