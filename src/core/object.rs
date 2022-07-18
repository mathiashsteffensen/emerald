use lazy_static::lazy_static;
use std::collections::HashMap;
use std::default::Default;
use std::sync::Arc;

use crate::core;
use crate::object::{BuiltInMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};

pub const NAME: &str = "Object";

lazy_static! {
    pub static ref EM_MAIN_OBJ: Arc<EmeraldObject> = em_instance();
}

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::basic_object::NAME),
        HashMap::from([
            ("puts".to_string(), em_object_puts as BuiltInMethod),
            ("to_s".to_string(), em_object_to_s as BuiltInMethod),
        ]),
    ))
    .unwrap()
}

pub fn em_instance() -> Arc<EmeraldObject> {
    EmeraldObject::new_instance(NAME)
}

fn em_object_puts(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    for arg in args {
        let stringified = arg
            .send(arg.clone(), "to_s", Arc::clone(&ctx), Default::default())
            .unwrap();

        match &stringified.underlying_value {
            UnderlyingValueType::String(s) => {
                println!("{}", s)
            }
            _ => panic!(
                "to_s did not return a string, got {}",
                stringified.class_name()
            ),
        }
    }

    Arc::clone(&core::nil_class::EM_NIL)
}

fn em_object_to_s(
    ctx: Arc<ExecutionContext>,
    _args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    let q_self = ctx.borrow_self();
    let class_name = q_self.class_name();

    core::string::em_instance(format!("#<{class_name}:{q_self:p}>"))
}
