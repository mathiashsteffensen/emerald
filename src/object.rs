use std::collections::HashMap;
use std::sync::Arc;

use crate::core;

pub type BuiltInMethod = fn(Arc<ExecutionContext>, Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject>;

#[derive(Debug)]
pub enum UnderlyingValueType {
    None,
    Class(String),
    Symbol(String),
    String(String),
    Integer(i64),
    True,
    False,
    Nil,
}

#[derive(Debug)]
pub struct EmeraldObject {
    pub class: Option<Arc<EmeraldObject>>,
    pub q_super: Option<Arc<EmeraldObject>>,
    pub built_in_method_set: HashMap<String, BuiltInMethod>,
    pub underlying_value: UnderlyingValueType,
}

impl EmeraldObject {
    pub fn new_instance(class_name: &str) -> Arc<EmeraldObject> {
        Arc::from(EmeraldObject {
            class: Some(core::em_get_class(class_name).unwrap()),
            q_super: None,
            built_in_method_set: Default::default(),
            underlying_value: UnderlyingValueType::None,
        })
    }

    pub fn new_instance_with_underlying_value(
        class_name: &str,
        val: UnderlyingValueType,
    ) -> Arc<EmeraldObject> {
        Arc::from(EmeraldObject {
            class: Some(core::em_get_class(class_name).unwrap()),
            q_super: None,
            built_in_method_set: Default::default(),
            underlying_value: val,
        })
    }

    pub fn new_class(
        name: &str,
        q_super: Option<Arc<EmeraldObject>>,
        built_in_method_set: HashMap<String, BuiltInMethod>,
    ) -> Arc<EmeraldObject> {
        let underlying_value = UnderlyingValueType::Class(name.to_string());

        Arc::from(EmeraldObject {
            class: None,
            q_super,
            built_in_method_set,
            underlying_value,
        })
    }

    pub fn send(
        &self,
        name: &str,
        ctx: Arc<ExecutionContext>,
        args: Vec<Arc<EmeraldObject>>,
    ) -> Result<Arc<EmeraldObject>, String> {
        if let Some(method) = self.method(name) {
            Ok(method(ctx, args))
        } else {
            Err(format!(
                "Undefined method {} for {:?}",
                name,
                self.borrow_class().underlying_value
            ))
        }
    }

    pub fn method(&self, name: &str) -> Option<BuiltInMethod> {
        match self.class.clone() {
            Some(class) => match class.built_in_method_set.get(name) {
                Some(method) => Some(*method),
                None => None,
            },
            None => None,
        }
    }

    pub fn responds_to(&self, name: &str) -> bool {
        if let Some(_) = self.method(name) {
            true
        } else {
            false
        }
    }

    pub fn borrow_class(&self) -> &Arc<EmeraldObject> {
        self.class.as_ref().unwrap()
    }
}

pub struct ExecutionContext {
    pub outer: Option<Arc<ExecutionContext>>,
    pub q_self: Arc<EmeraldObject>,
}

impl ExecutionContext {
    pub fn new(q_self: Arc<EmeraldObject>) -> ExecutionContext {
        ExecutionContext {
            outer: None,
            q_self,
        }
    }

    pub fn with_outer(
        q_self: Arc<EmeraldObject>,
        outer: Arc<ExecutionContext>,
    ) -> ExecutionContext {
        ExecutionContext {
            outer: Some(Arc::clone(&outer)),
            q_self,
        }
    }

    pub fn borrow_self(&self) -> Arc<EmeraldObject> {
        Arc::clone(&self.q_self)
    }
}
