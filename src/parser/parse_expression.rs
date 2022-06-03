use crate::ast::node::Expression;
use crate::parser::Parser;

pub fn exec(p: &mut Parser, precedence: i16) -> Option<Expression> {
    let mut left = p.parse_as_prefix();

    while !p.peek_token_is_semicolon()
        && !p.peek_token_is_newline()
        && precedence < p.peek_precedence()
    {
        left = match left {
            Some(ref left_expr) => {
                let right = p.parse_as_infix(left_expr.clone());

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
