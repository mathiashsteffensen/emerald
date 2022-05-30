use std::collections::HashMap;
use std::rc::Rc;

pub mod class;

pub type BuiltInMethod = fn(Rc<ExecutionContext>, Vec<Rc<EmeraldObject>>) -> Rc<EmeraldObject>;

#[derive(Debug)]
pub enum UnderlyingValueType {
    Integer(i64),
    String(String),
    Class(String),
    None,
}

#[derive(Debug)]
pub struct EmeraldObject {
    pub class: Option<Rc<EmeraldObject>>,
    pub q_super: Option<Rc<EmeraldObject>>,
    pub built_in_method_set: HashMap<String, Rc<BuiltInMethod>>,
    pub underlying_value: UnderlyingValueType,
}

impl EmeraldObject {
    pub fn send(
        &self,
        name: &str,
        ctx: Rc<ExecutionContext>,
        args: Vec<Rc<EmeraldObject>>,
    ) -> Result<Rc<EmeraldObject>, String> {
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

    pub fn method(&self, name: &str) -> Option<Rc<BuiltInMethod>> {
        match self.class.clone() {
            Some(class) => match class.built_in_method_set.get(name) {
                Some(method) => Some(Rc::clone(method)),
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

    pub fn borrow_class(&self) -> &Rc<EmeraldObject> {
        self.class.as_ref().unwrap()
    }
}

pub struct ExecutionContext {
    pub outer: Option<Rc<ExecutionContext>>,
    pub const_map: HashMap<String, Rc<EmeraldObject>>,
    pub q_self: Rc<EmeraldObject>,
}

impl ExecutionContext {
    pub fn new(
        q_self: Rc<EmeraldObject>,
        const_map: HashMap<String, Rc<EmeraldObject>>,
    ) -> ExecutionContext {
        ExecutionContext {
            outer: None,
            const_map,
            q_self,
        }
    }

    pub fn with_outer(q_self: Rc<EmeraldObject>, outer: Rc<ExecutionContext>) -> ExecutionContext {
        ExecutionContext {
            outer: Some(Rc::clone(&outer)),
            const_map: outer.const_map.clone(),
            q_self,
        }
    }

    pub fn get_const(&self, name: &str) -> Rc<EmeraldObject> {
        Rc::clone(self.const_map.get(name).unwrap())
    }

    pub fn borrow_self(&self) -> Rc<EmeraldObject> {
        Rc::clone(&self.q_self)
    }
}
