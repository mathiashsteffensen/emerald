use crate::ast::node::Expression;
use crate::parser::{parse_as_infix, parse_as_prefix, Parser};

pub fn exec(p: &mut Parser, precedence: i16) -> Option<Expression> {
    let mut left = parse_as_prefix::exec(p);

    while !p.peek_token_is_semicolon()
        && !p.peek_token_is_newline()
        && precedence < p.peek_precedence()
    {
        left = match left {
            Some(ref left_expr) => {
                let right = parse_as_infix::exec(p, left_expr.clone());

                match right {
                    None => break,
                    _ => right,
                }
            }
            None => break,
        }
    }

    left.clone()
}
