use crate::ast::node::Statement;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser) -> Option<Statement> {
    let expression = parse_expression::exec(p, precedence::LOWEST);

    p.next_if_semicolon_or_newline();

    match expression {
        Some(expr) => Some(Statement::ExpressionStatement(expr)),
        None => None,
    }
}
