use tera::Error;

#[derive(Debug)]
pub enum GeneratorError {
    TemplateError(String),
    FileError(String),
}

impl GeneratorError {
    pub fn from_tera(err: Error) -> Self {
        Self::TemplateError(err.to_string())
    }

    pub fn file_error(err: std::io::Error) -> Self {
        Self::FileError(err.to_string())
    }
}
