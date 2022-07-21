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
        ("+".to_string(), em_string_add as BuiltInMethod),
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

pub fn em_string_add(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    string_infix_operator(ctx, args, "+", |left, right| left + &*right)
}

fn string_infix_operator<F>(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
    op: &str,
    cb: F,
) -> Arc<EmeraldObject>
where
    F: Fn(String, String) -> String,
{
    infix_operator(ctx, args, op, |left, right| em_instance(cb(left, right)))
}

// fn boolean_infix_operator<F>(
//     ctx: Arc<ExecutionContext>,
//     args: Vec<Arc<EmeraldObject>>,
//     op: &str,
//     cb: F,
// ) -> Arc<EmeraldObject>
//     where
//         F: Fn(String, String) -> bool,
// {
//     infix_operator(ctx, args, op, |left, right| {
//         if cb(left, right) {
//             Arc::clone(&core::true_class::EM_TRUE)
//         } else {
//             Arc::clone(&core::false_class::EM_FALSE)
//         }
//     })
// }

fn infix_operator<F>(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
    op: &str,
    cb: F,
) -> Arc<EmeraldObject>
where
    F: Fn(String, String) -> Arc<EmeraldObject>,
{
    extract_underlying_value(
        ctx.borrow_self(),
        {
            |left| {
                extract_underlying_value(
                    Arc::clone(args.get(0).unwrap()),
                    |right| cb(left.clone(), right),
                    format!("String#{} was not passed a string", op),
                )
            }
        },
        format!("Calling String#{} on not a string?!?!?!", op),
    )
}

fn extract_underlying_value<F, TReturns>(obj: Arc<EmeraldObject>, cb: F, reject: String) -> TReturns
where
    F: Fn(String) -> TReturns,
{
    match &obj.underlying_value {
        UnderlyingValueType::String(val) => cb(val.clone()),
        _ => panic!("{}, got {}", reject, obj.class_name()),
    }
}
