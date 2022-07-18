use std::collections::HashMap;
use std::fmt;
use std::sync::{Arc, Mutex};

use crate::compiler::bytecode::Bytecode;
use crate::{core, kernel};

pub type BuiltInMethod = fn(Arc<ExecutionContext>, Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject>;

pub enum EmeraldMethod {
    BuiltIn(BuiltInMethod),
    Compiled(Block),
}

#[derive(Debug, Clone)]
pub struct Block {
    pub bytecode: Bytecode,
    pub arity: u8,
    pub num_locals: u16,
}

impl Block {
    pub fn new(bytecode: Bytecode, arity: u8, num_locals: u16) -> Block {
        Block {
            bytecode,
            arity,
            num_locals,
        }
    }
}

#[derive(Debug)]
pub enum UnderlyingValueType {
    None,
    Class(String),
    Symbol(String),
    String(String),
    Proc(Block),
    Integer(i64),
    True,
    False,
    Nil,
}

impl fmt::Display for UnderlyingValueType {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        let format = match self {
            UnderlyingValueType::String(string) => string.clone(),
            UnderlyingValueType::Integer(int) => int.to_string(),
            _ => self.to_string(),
        };

        write!(f, "{}", format)
    }
}

#[derive(Debug)]
pub struct EmeraldObject {
    _class: Option<Arc<EmeraldObject>>,
    pub q_super: Option<Arc<EmeraldObject>>,
    pub built_in_method_set: HashMap<String, BuiltInMethod>,
    pub defined_method_set: Mutex<HashMap<String, Block>>,
    pub underlying_value: UnderlyingValueType,
}

impl EmeraldObject {
    pub fn new_instance(class_name: &str) -> Arc<EmeraldObject> {
        Arc::from(EmeraldObject {
            _class: Some(core::em_get_class(class_name).unwrap()),
            q_super: None,
            built_in_method_set: Default::default(),
            defined_method_set: Default::default(),
            underlying_value: UnderlyingValueType::None,
        })
    }

    pub fn new_instance_with_underlying_value(
        class_name: &str,
        val: UnderlyingValueType,
    ) -> Arc<EmeraldObject> {
        Arc::from(EmeraldObject {
            _class: Some(core::em_get_class(class_name).unwrap()),
            q_super: None,
            built_in_method_set: Default::default(),
            defined_method_set: Default::default(),
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
            _class: None,
            q_super,
            built_in_method_set,
            defined_method_set: Default::default(),
            underlying_value,
        })
    }

    pub fn send(
        &self,
        receiver: Arc<EmeraldObject>,
        name: &str,
        ctx: Arc<ExecutionContext>,
        args: Vec<Arc<EmeraldObject>>,
    ) -> Result<Arc<EmeraldObject>, String> {
        if let Some(method) = self.method(name) {
            match method {
                EmeraldMethod::BuiltIn(method) => Ok(method(
                    Arc::from(ExecutionContext::with_outer(receiver, ctx)),
                    args,
                )),
                EmeraldMethod::Compiled(block) => {
                    Ok(kernel::execute_bytecode(block.bytecode, args))
                }
            }
        } else {
            Err(format!(
                "Undefined method {} for {:?}",
                name,
                self.class_name(),
            ))
        }
    }

    pub fn define_method(&self, name: String, block: Block) {
        self.class()
            .defined_method_set
            .lock()
            .unwrap()
            .insert(name, block);
    }

    pub fn method(&self, name: &str) -> Option<EmeraldMethod> {
        let class = self.class();

        if let Some(method) = class.built_in_method_set.get(name) {
            return Some(EmeraldMethod::BuiltIn(*method));
        }

        if let Some(method) = class.defined_method_set.lock().unwrap().get(name) {
            return Some(EmeraldMethod::Compiled(method.clone()));
        }

        return None;
    }

    pub fn responds_to(&self, name: &str) -> bool {
        matches!(self.method(name), Some(_))
    }

    pub fn class(&self) -> Arc<EmeraldObject> {
        self._class.as_ref().unwrap().clone()
    }

    // Returns the class name as a string
    pub fn class_name(&self) -> String {
        if let UnderlyingValueType::Class(name) = &self.class().underlying_value {
            name.clone()
        } else {
            unreachable!()
        }
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
