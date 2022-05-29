use std::collections::HashMap;
use std::rc::Rc;

pub mod class;

pub type BuiltInMethod = fn(ExecutionContext, Vec<EmeraldObject>) -> EmeraldObject;

pub enum UnderlyingValueType {
    Integer(i64),
    String(String),
    Class(String),
    None,
}

pub struct EmeraldObject {
    pub class: Option<Rc<EmeraldObject>>,
    pub q_super: Option<Rc<EmeraldObject>>,
    pub built_in_method_set: HashMap<String, BuiltInMethod>,
    pub underlying_value: UnderlyingValueType,
}

pub struct ExecutionContext {
    pub outer: Box<Option<ExecutionContext>>,
    pub q_self: Rc<EmeraldObject>,
}
