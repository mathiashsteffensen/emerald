use crate::core;
use crate::object::{BuiltInMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};
use std::collections::HashMap;
use std::ops::Add;
use std::sync::Arc;

pub const NAME: &str = "Symbol";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::object::NAME),
        HashMap::from([("inspect".to_string(), em_symbol_inspect as BuiltInMethod)]),
    ))
    .unwrap()
}

pub fn em_instance(val: String) -> Arc<EmeraldObject> {
    Arc::from(EmeraldObject::new_instance_with_underlying_value(
        NAME,
        UnderlyingValueType::Symbol(val),
    ))
}

fn em_symbol_inspect(
    ctx: Arc<ExecutionContext>,
    _args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    extract_underlying_value(ctx, |val| {
        core::string::em_instance(":".to_string().add(&*val.clone()))
    })
}

fn extract_underlying_value<F, TReturns>(ctx: Arc<ExecutionContext>, cb: F) -> TReturns
where
    F: Fn(&String) -> TReturns,
{
    match &ctx.q_self.underlying_value {
        UnderlyingValueType::Symbol(val) => cb(val),
        _ => unreachable!(),
    }
}
