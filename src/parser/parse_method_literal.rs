use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_paren_delimited_expression_list, parse_statement_list, Parser};

pub fn exec(p: &mut Parser, data: token::TokenData) -> Option<Expression> {
    p.next_token();

    if let token::Token::Ident(ident_data) = p.current_token.clone() {
        let name = Expression::IdentifierExpression(ident_data);

        let args = if let Some(expr) = parse_paren_delimited_expression_list::exec(p) {
            expr
        } else {
            return None;
        };

        p.next_token();
        p.next_token();

        let body = parse_statement_list::exec(p, |token| matches!(token, token::Token::End(_data)));

        Some(Expression::MethodLiteral(data, Box::new(name), args, body))
    } else {
        None
    }
}
