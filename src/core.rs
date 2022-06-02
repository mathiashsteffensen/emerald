use lazy_static::lazy_static;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

use crate::object::{EmeraldObject, UnderlyingValueType};

pub mod all;
pub mod basic_object;
pub mod integer;
pub mod object;
pub mod string;
pub mod symbol;
pub mod true_class;

lazy_static! {
    static ref EM_CLASS_MAP: Mutex<HashMap<String, Arc<EmeraldObject>>> =
        Mutex::new(HashMap::new());
}

pub fn em_define_class(class: Arc<EmeraldObject>) -> Result<(), String> {
    match &class.underlying_value {
        UnderlyingValueType::Class(name) => {
            EM_CLASS_MAP.lock().unwrap().insert(name.clone(), class);
            Ok(())
        }
        _ => Err(format!(
            "em_define_class expected class but got {:?}",
            class
        )),
    }
}

pub fn em_get_class(name: &str) -> Option<Arc<EmeraldObject>> {
    let map = EM_CLASS_MAP.lock().unwrap();
    let class = map.get(name);

    match class {
        Some(class) => Some(Arc::clone(class)),
        None => None,
    }
}
