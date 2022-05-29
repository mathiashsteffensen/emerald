use std::collections::HashMap;
use std::rc::Rc;

use crate::core;
use crate::object::EmeraldObject;

pub fn map() -> HashMap<String, Rc<EmeraldObject>> {
    HashMap::from([
        (core::integer::NAME.to_string(), core::integer::em_class()),
        (core::string::NAME.to_string(), core::string::em_class()),
    ])
}
