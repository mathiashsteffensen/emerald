use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{
    parse_class_literal, parse_float_literal, parse_grouped_expression,
    parse_identifier_expression, parse_if_expression, parse_integer_literal, parse_method_literal,
    parse_prefix_expression, Parser,
};

pub fn exec(p: &mut Parser) -> Option<Expression> {
    let tok = p.current_token.clone();

    match tok {
        token::Token::Bang(data) => parse_prefix_expression::exec(p, data),
        token::Token::Minus(data) => parse_prefix_expression::exec(p, data),
        token::Token::Ident(data) => parse_identifier_expression::exec(p, data),
        token::Token::Def(data) => parse_method_literal::exec(p, data),
        token::Token::If(data) => parse_if_expression::exec(p, data),
        token::Token::Class(data) => parse_class_literal::exec(p, data),
        token::Token::Int(data) => parse_integer_literal::exec(p, data),
        token::Token::Float(data) => parse_float_literal::exec(p, data),
        token::Token::String(data) => p.parse_string_literal(data),
        token::Token::True(data) => p.parse_true_literal(data),
        token::Token::False(data) => p.parse_false_literal(data),
        token::Token::Nil(data) => p.parse_nil_literal(data),
        token::Token::LeftParen(_data) => parse_grouped_expression::exec(p),
        _ => {
            p.add_error(&*format!(
                "No prefix parse function found for token {:?}",
                p.current_token,
            ));

            None
        }
    }
}
