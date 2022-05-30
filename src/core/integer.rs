use std::collections::HashMap;
use std::rc::Rc;

use crate::core;
use crate::object;
use crate::object::{BuiltInMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};

pub const NAME: &str = "Integer";

pub fn em_class() -> Rc<EmeraldObject> {
    let method_set = HashMap::from([
        ("+".to_string(), Rc::from(em_integer_add as BuiltInMethod)),
        ("-".to_string(), Rc::from(em_integer_sub as BuiltInMethod)),
        ("*".to_string(), Rc::from(em_integer_mul as BuiltInMethod)),
        (
            "to_s".to_string(),
            Rc::from(em_integer_to_s as BuiltInMethod),
        ),
    ]);

    object::class::new(NAME, method_set)
}

pub fn em_instance(class: Rc<EmeraldObject>, val: i64) -> Rc<EmeraldObject> {
    Rc::from(EmeraldObject {
        class: Some(class),
        q_super: None,
        built_in_method_set: Default::default(),
        underlying_value: UnderlyingValueType::Integer(val),
    })
}

fn em_integer_to_s(ctx: Rc<ExecutionContext>, _args: Vec<Rc<EmeraldObject>>) -> Rc<EmeraldObject> {
    extract_underlying_value(
        ctx.borrow_self(),
        { |val| core::string::em_instance(ctx.get_const("String"), val.to_string()) },
        format!("Calling Integer#to_s on not an integer?!?!?!"),
    )
}

fn em_integer_add(ctx: Rc<ExecutionContext>, args: Vec<Rc<EmeraldObject>>) -> Rc<EmeraldObject> {
    integer_infix_operator(ctx, args, "+", |l, r| l + r)
}

fn em_integer_sub(ctx: Rc<ExecutionContext>, args: Vec<Rc<EmeraldObject>>) -> Rc<EmeraldObject> {
    integer_infix_operator(ctx, args, "-", |l, r| l - r)
}

fn em_integer_mul(ctx: Rc<ExecutionContext>, args: Vec<Rc<EmeraldObject>>) -> Rc<EmeraldObject> {
    integer_infix_operator(ctx, args, "*", |l, r| l * r)
}

fn integer_infix_operator<F>(
    ctx: Rc<ExecutionContext>,
    args: Vec<Rc<EmeraldObject>>,
    op: &str,
    cb: F,
) -> Rc<EmeraldObject>
where
    F: Fn(i64, i64) -> i64,
{
    extract_underlying_value(
        ctx.borrow_self(),
        {
            |left| {
                extract_underlying_value(
                    Rc::clone(args.get(0).unwrap()),
                    { |right| em_instance(Rc::clone(ctx.q_self.borrow_class()), cb(left, right)) },
                    format!("Integer#{} was not passed an integer", op),
                )
            }
        },
        format!("Calling Integer#{} on not an integer?!?!?!", op),
    )
}

fn extract_underlying_value<F, TReturns>(obj: Rc<EmeraldObject>, cb: F, reject: String) -> TReturns
where
    F: Fn(i64) -> TReturns,
{
    match obj.underlying_value {
        UnderlyingValueType::Integer(val) => cb(val),
        _ => panic!("{}", reject),
    }
}
