use crate::ast::node::Expression;
use crate::lexer;
use crate::parser::{parse_expression_list, Parser};

pub fn exec(p: &mut Parser) -> Option<Vec<Expression>> {
    if matches!(&p.peek_token, lexer::token::Token::LeftParen(_data)) {
        p.next_token();

        if let Some(expressions) = parse_expression_list::exec(p) {
            p.next_token();

            Some(expressions)
        } else {
            None
        }
    } else {
        Some(Vec::new())
    }
}
