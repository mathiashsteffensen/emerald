use crate::ast::node::Expression;
use crate::lexer;
use crate::parser::{parse_expression, precedence, Parser};

pub fn exec(p: &mut Parser) -> Option<Vec<Expression>> {
    let mut args: Vec<Expression> = Vec::new();

    p.next_token();

    if let Some(expr) = parse_expression::exec(p, precedence::LOWEST) {
        args.push(expr);
    } else {
        return None;
    };

    let mut is_comma = matches!(&p.peek_token, lexer::token::Token::Comma(_data));

    while is_comma {
        p.next_token();
        p.next_token();

        if let Some(expr) = parse_expression::exec(p, precedence::LOWEST) {
            args.push(expr);
        } else {
            return None;
        };

        if p.peek_token_is_newline() {
            p.next_token();
        }

        is_comma = matches!(&p.peek_token, lexer::token::Token::Comma(_data));
    }

    Some(args)
}
