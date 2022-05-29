use std::rc::Rc;

use crate::object::{EmeraldObject, UnderlyingValueType};

pub fn new_class(name: &str) -> Rc<EmeraldObject> {
    let underlying_value = UnderlyingValueType::Class(name.to_string());

    Rc::from(EmeraldObject {
        class: None,
        q_super: None,
        built_in_method_set: Default::default(),
        underlying_value,
    })
}
