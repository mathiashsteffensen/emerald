use crate::ast::node::Statement;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser) -> Option<Statement> {
    let expression = parse_expression::exec(p, precedence::LOWEST);

    while p.peek_token_is_newline() || p.peek_token_is_semicolon() {
        p.next_token()
    }

    match expression {
        Some(expr) => Some(Statement::ExpressionStatement(expr)),
        None => None,
    }
}
