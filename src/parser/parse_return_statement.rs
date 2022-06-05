use crate::ast::node::{Expression, Statement};
use crate::lexer::token::TokenData;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser, data: TokenData) -> Option<Statement> {
    let return_statement =
        { |expr: Option<Expression>| Some(Statement::ReturnStatement(data, expr)) };

    if p.peek_token_is_newline() || p.peek_token_is_semicolon() || p.peek_token_is_eof() {
        p.next_token();
        return return_statement(None);
    }

    p.next_token();

    let expression = parse_expression::exec(p, precedence::LOWEST);

    p.next_if_semicolon_or_newline();

    return_statement(expression)
}
