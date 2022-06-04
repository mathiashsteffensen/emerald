use std::collections::HashMap;
use std::default::Default;
use std::sync::Arc;

use crate::core;
use crate::object::{BuiltInMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};

pub const NAME: &str = "Object";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::basic_object::NAME),
        HashMap::from([("puts".to_string(), em_object_puts as BuiltInMethod)]),
    ))
    .unwrap()
}

pub fn em_instance() -> Arc<EmeraldObject> {
    EmeraldObject::new_instance(NAME)
}

fn em_object_puts(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    for arg in args {
        let stringified = arg
            .send("to_s", Arc::clone(&ctx), Default::default())
            .unwrap();

        match &stringified.underlying_value {
            UnderlyingValueType::String(s) => {
                println!("{}", s)
            }
            _ => panic!("to_s did not return a string"),
        }
    }

    Arc::clone(&core::nil_class::EM_NIL)
}
