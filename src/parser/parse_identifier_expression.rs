use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{
    parse_assignment_expression, parse_method_call_with_receiver,
    parse_method_call_without_receiver, Parser,
};

pub fn exec(p: &mut Parser, data: token::TokenData) -> Option<Expression> {
    let ident = Expression::IdentifierExpression(data);

    match p.peek_token.clone() {
        token::Token::Assign(data) => parse_assignment_expression::exec(p, ident, data),
        token::Token::LeftParen(_data) => parse_method_call_without_receiver::exec(p, ident),
        token::Token::Dot(_data) => parse_method_call_with_receiver::exec(p, ident),
        _ => Some(ident),
    }
}
