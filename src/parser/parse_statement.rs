use crate::ast::node::Statement;
use crate::lexer::token::Token;
use crate::parser::{parse_expression_statement, parse_return_statement, Parser};

pub fn exec(p: &mut Parser) -> Option<Statement> {
    match p.current_token.clone() {
        Token::Return(data) => parse_return_statement::exec(p, data),
        _ => parse_expression_statement::exec(p),
    }
}
