//! Structures for managing function registrations.

use std::{collections::HashMap, sync::Arc};

/// Prefix used for environment variables that hold exported function definitions.
pub(crate) const EXPORTED_FUNCTION_ENV_VAR_PREFIX: &str = "BASH_FUNC_";
/// Suffix used for environment variables that hold exported function definitions.
pub(crate) const EXPORTED_FUNCTION_ENV_VAR_SUFFIX: &str = "%%";

/// An environment for defined, named functions.
#[derive(Clone, Default)]
pub struct FunctionEnv {
    functions: HashMap<String, FunctionRegistration>,
}

impl FunctionEnv {
    /// Tries to retrieve the registration for a function by name.
    ///
    /// # Arguments
    ///
    /// * `name` - The name of the function to retrieve.
    pub fn get(&self, name: &str) -> Option<&FunctionRegistration> {
        self.functions.get(name)
    }

    /// Tries to retrieve a mutable reference to the registration for a function by name.
    ///
    /// # Arguments
    ///
    /// * `name` - The name of the function to retrieve.
    pub fn get_mut(&mut self, name: &str) -> Option<&mut FunctionRegistration> {
        self.functions.get_mut(name)
    }

    /// Unregisters a function from the environment.
    ///
    /// # Arguments
    ///
    /// * `name` - The name of the function to remove.
    pub fn remove(&mut self, name: &str) -> Option<FunctionRegistration> {
        self.functions.remove(name)
    }

    /// Updates a function registration in this environment.
    ///
    /// # Arguments
    ///
    /// * `name` - The name of the function to update.
    /// * `definition` - The new definition for the function.
    pub fn update(&mut self, name: String, definition: Arc<brush_parser::ast::FunctionDefinition>) {
        // N.B. If the function was already registered, preserve its exported attribute.
        let exported = self
            .functions
            .get(&name)
            .is_some_and(FunctionRegistration::is_exported);
        self.functions.insert(
            name,
            FunctionRegistration {
                definition,
                exported,
            },
        );
    }

    /// Returns an iterator over the functions registered in this environment.
    pub fn iter(&self) -> impl Iterator<Item = (&String, &FunctionRegistration)> {
        self.functions.iter()
    }

    /// Returns an iterator over the functions registered in this environment that
    /// are marked for export to child processes.
    pub fn iter_exported(&self) -> impl Iterator<Item = (&String, &FunctionRegistration)> {
        self.functions.iter().filter(|(_, v)| v.is_exported())
    }
}

/// Encapsulates a registration for a defined function.
#[derive(Clone)]
pub struct FunctionRegistration {
    /// The definition of the function.
    pub(crate) definition: Arc<brush_parser::ast::FunctionDefinition>,
    /// Whether or not the function is marked for export to child processes.
    exported: bool,
}

impl FunctionRegistration {
    /// Returns whether or not the function is exported to child processes.
    pub fn is_exported(&self) -> bool {
        self.exported
    }

    /// Marks the function as exported to child processes.
    pub fn export(&mut self) {
        self.exported = true;
    }

    /// Marks the function as not exported to child processes.
    pub fn unexport(&mut self) {
        self.exported = false;
    }
}
