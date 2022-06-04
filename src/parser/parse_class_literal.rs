use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_statement_list, Parser};

pub fn exec(p: &mut Parser, data: token::TokenData) -> Option<Expression> {
    p.next_token();

    if let token::Token::Ident(ident_data) = p.current_token.clone() {
        let name = Expression::IdentifierExpression(ident_data);

        p.next_token();

        let body = parse_statement_list::exec(p, |token| matches!(token, token::Token::End(_data)));

        Some(Expression::ClassLiteral(data, Box::new(name), body))
    } else {
        None
    }
}
