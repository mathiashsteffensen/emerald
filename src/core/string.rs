use crate::core;
use crate::object::{BuiltInMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};
use std::collections::HashMap;
use std::ops::Add;
use std::sync::Arc;

pub const NAME: &str = "String";

pub fn em_init_class() {
    let method_set = HashMap::from([
        ("to_s".to_string(), em_string_to_s as BuiltInMethod),
        ("inspect".to_string(), em_string_inspect as BuiltInMethod),
    ]);

    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::object::NAME),
        method_set,
    ))
    .unwrap()
}

pub fn em_instance(val: String) -> Arc<EmeraldObject> {
    EmeraldObject::new_instance_with_underlying_value(NAME, UnderlyingValueType::String(val))
}

pub fn em_string_to_s(
    ctx: Arc<ExecutionContext>,
    _args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    ctx.borrow_self()
}

pub fn em_string_inspect(
    ctx: Arc<ExecutionContext>,
    _args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    if let UnderlyingValueType::String(str) = &ctx.borrow_self().underlying_value {
        let mut out = "\"".to_string();

        out.push_str(str.as_str());

        em_instance(out.add("\""))
    } else {
        Arc::clone(&core::true_class::EM_TRUE)
    }
}
