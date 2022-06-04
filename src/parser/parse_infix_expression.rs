use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_expression, Parser};

pub fn exec(p: &mut Parser, data: token::TokenData, left: Expression) -> Option<Expression> {
    let precedence = p.peek_precedence();

    p.next_token();
    p.next_token();

    let right = parse_expression::exec(p, precedence);

    match right {
        Some(right_expr) => Some(Expression::InfixExpression(
            Box::new(left),
            data,
            Box::new(right_expr),
        )),
        None => None,
    }
}
