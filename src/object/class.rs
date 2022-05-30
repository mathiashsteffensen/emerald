use std::collections::HashMap;
use std::rc::Rc;

use crate::object::{BuiltInMethod, EmeraldObject, UnderlyingValueType};

pub fn new(
    name: &str,
    built_in_method_set: HashMap<String, Rc<BuiltInMethod>>,
) -> Rc<EmeraldObject> {
    let underlying_value = UnderlyingValueType::Class(name.to_string());

    Rc::from(EmeraldObject {
        class: None,
        q_super: None,
        built_in_method_set,
        underlying_value,
    })
}
