use std::rc::Rc;

use crate::object;
use crate::object::{EmeraldObject, UnderlyingValueType};

pub const NAME: &str = "Integer";

pub fn em_class() -> Rc<EmeraldObject> {
    object::class::new_class(NAME)
}

pub fn em_instance(class: Rc<EmeraldObject>, val: i64) -> Rc<EmeraldObject> {
    Rc::from(EmeraldObject {
        class: Some(class),
        q_super: None,
        built_in_method_set: Default::default(),
        underlying_value: UnderlyingValueType::Integer(val),
    })
}
