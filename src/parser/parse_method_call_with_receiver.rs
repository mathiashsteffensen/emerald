use crate::ast::node;
use crate::ast::node::{Block, Expression};
use crate::lexer::token;
use crate::parser::{parse_identifier_expression, parse_paren_delimited_expression_list, Parser};

pub fn exec(p: &mut Parser, receiver: Expression) -> Option<Expression> {
    p.next_token();

    let ident = match &p.peek_token {
        token::Token::Ident(data) => {
            if let Some(expr) = parse_identifier_expression::exec(p, data.clone()) {
                expr
            } else {
                return None;
            }
        }
        _ => {
            p.add_syntax_error(p.peek_token.data(), "ident");
            return None;
        }
    };

    p.next_token();

    let args = if let Some(expr) = parse_paren_delimited_expression_list::exec(p) {
        expr
    } else {
        return None;
    };

    p.next_token();

    Some(Expression::MethodCall(node::MethodCallData {
        receiver: Some(Box::new(receiver)),
        ident: Box::new(ident),
        block: Block::new(Vec::new(), Vec::new()),
        args,
    }))
}
