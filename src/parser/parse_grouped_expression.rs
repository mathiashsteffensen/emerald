use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser) -> Option<Expression> {
    p.next_token();

    let expr = parse_expression::exec(p, precedence::LOWEST);

    match p.peek_token.clone() {
        token::Token::RightParen(_data) => {
            p.next_token();
            expr
        }
        _ => {
            p.add_syntax_error(p.peek_token.data(), ")");

            None
        }
    }
}
