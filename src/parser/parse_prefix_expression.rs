use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser, data: token::TokenData) -> Option<Expression> {
    p.next_token();

    let right = parse_expression::exec(p, precedence::PREFIX);

    match right {
        Some(expr) => Some(Expression::PrefixExpression(data, Box::new(expr))),
        None => None,
    }
}
