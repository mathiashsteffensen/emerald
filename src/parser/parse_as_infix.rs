use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_infix_expression, parse_method_call_with_receiver, Parser};

pub fn exec(p: &mut Parser, left: Expression) -> Option<Expression> {
    let tok = p.peek_token.clone();

    let expr = match tok {
        token::Token::Plus(data) => parse_infix_expression::exec(p, data, left),
        token::Token::Minus(data) => parse_infix_expression::exec(p, data, left),
        token::Token::Asterisk(data) => parse_infix_expression::exec(p, data, left),
        token::Token::Slash(data) => parse_infix_expression::exec(p, data, left),
        token::Token::LessThan(data) => parse_infix_expression::exec(p, data, left),
        token::Token::LessThanOrEquals(data) => parse_infix_expression::exec(p, data, left),
        token::Token::GreaterThan(data) => parse_infix_expression::exec(p, data, left),
        token::Token::GreaterThanOrEquals(data) => parse_infix_expression::exec(p, data, left),
        token::Token::Equals(data) => parse_infix_expression::exec(p, data, left),
        token::Token::NotEquals(data) => parse_infix_expression::exec(p, data, left),
        token::Token::Dot(_data) => parse_method_call_with_receiver::exec(p, left),
        _ => None,
    };

    expr
}
