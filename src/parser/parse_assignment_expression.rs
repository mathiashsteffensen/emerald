use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser, ident: Expression, data: token::TokenData) -> Option<Expression> {
    p.next_token();
    p.next_token();

    let val = parse_expression::exec(p, precedence::LOWEST);

    match val {
        Some(expr) => Some(Expression::AssignmentExpression(
            Box::new(ident),
            data,
            Box::new(expr),
        )),
        None => None,
    }
}
