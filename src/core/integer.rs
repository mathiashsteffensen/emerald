use std::collections::HashMap;
use std::sync::Arc;

use crate::core;
use crate::object::{BuiltInMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};

pub const NAME: &str = "Integer";

pub fn em_init_class() {
    let method_set = HashMap::from([
        ("+".to_string(), em_integer_add as BuiltInMethod),
        ("-".to_string(), em_integer_sub as BuiltInMethod),
        ("*".to_string(), em_integer_mul as BuiltInMethod),
        ("/".to_string(), em_integer_div as BuiltInMethod),
        (">".to_string(), em_integer_greater_than as BuiltInMethod),
        (
            ">=".to_string(),
            em_integer_greater_than_or_eq as BuiltInMethod,
        ),
        ("<".to_string(), em_integer_less_than as BuiltInMethod),
        (
            "<=".to_string(),
            em_integer_less_than_or_eq as BuiltInMethod,
        ),
        ("==".to_string(), em_integer_eq as BuiltInMethod),
        ("to_s".to_string(), em_integer_to_s as BuiltInMethod),
        ("inspect".to_string(), em_integer_to_s as BuiltInMethod),
    ]);

    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::object::NAME),
        method_set,
    ))
    .unwrap()
}

pub fn em_instance(val: i64) -> Arc<EmeraldObject> {
    Arc::from(EmeraldObject::new_instance_with_underlying_value(
        NAME,
        UnderlyingValueType::Integer(val),
    ))
}

fn em_integer_to_s(
    ctx: Arc<ExecutionContext>,
    _args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    extract_underlying_value(
        ctx.borrow_self(),
        |val| core::string::em_instance(val.to_string()),
        format!("Calling Integer#to_s on not an integer?!?!?!"),
    )
}

fn em_integer_add(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    integer_infix_operator(ctx, args, "+", |l, r| l + r)
}

fn em_integer_sub(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    integer_infix_operator(ctx, args, "-", |l, r| l - r)
}

fn em_integer_mul(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    integer_infix_operator(ctx, args, "*", |l, r| l * r)
}

fn em_integer_div(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    integer_infix_operator(ctx, args, "/", |l, r| l / r)
}

fn em_integer_greater_than(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    boolean_infix_operator(ctx, args, ">", |l, r| l > r)
}

fn em_integer_greater_than_or_eq(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    boolean_infix_operator(ctx, args, ">=", |l, r| l >= r)
}

fn em_integer_less_than(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    boolean_infix_operator(ctx, args, "<", |l, r| l < r)
}

fn em_integer_less_than_or_eq(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    boolean_infix_operator(ctx, args, "<=", |l, r| l <= r)
}

fn em_integer_eq(ctx: Arc<ExecutionContext>, args: Vec<Arc<EmeraldObject>>) -> Arc<EmeraldObject> {
    boolean_infix_operator(ctx, args, "==", |l, r| l == r)
}

fn integer_infix_operator<F>(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
    op: &str,
    cb: F,
) -> Arc<EmeraldObject>
where
    F: Fn(i64, i64) -> i64,
{
    infix_operator(ctx, args, op, |left, right| em_instance(cb(left, right)))
}

fn boolean_infix_operator<F>(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
    op: &str,
    cb: F,
) -> Arc<EmeraldObject>
where
    F: Fn(i64, i64) -> bool,
{
    infix_operator(ctx, args, op, |left, right| {
        if cb(left, right) {
            Arc::clone(&core::true_class::EM_TRUE)
        } else {
            Arc::clone(&core::false_class::EM_FALSE)
        }
    })
}

fn infix_operator<F>(
    ctx: Arc<ExecutionContext>,
    args: Vec<Arc<EmeraldObject>>,
    op: &str,
    cb: F,
) -> Arc<EmeraldObject>
where
    F: Fn(i64, i64) -> Arc<EmeraldObject>,
{
    extract_underlying_value(
        ctx.borrow_self(),
        {
            |left| {
                extract_underlying_value(
                    Arc::clone(args.get(0).unwrap()),
                    |right| cb(left, right),
                    format!("Integer#{} was not passed an integer", op),
                )
            }
        },
        format!("Calling Integer#{} on not an integer?!?!?!", op),
    )
}

fn extract_underlying_value<F, TReturns>(obj: Arc<EmeraldObject>, cb: F, reject: String) -> TReturns
where
    F: Fn(i64) -> TReturns,
{
    match obj.underlying_value {
        UnderlyingValueType::Integer(val) => cb(val),
        _ => panic!("{}", reject),
    }
}
