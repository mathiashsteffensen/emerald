use crate::ast::node;
use crate::ast::node::{Block, Expression};
use crate::parser::{parse_paren_delimited_expression_list, Parser};

pub fn exec(p: &mut Parser, ident: Expression) -> Option<Expression> {
    let args = if let Some(expr) = parse_paren_delimited_expression_list::exec(p) {
        expr
    } else {
        return None;
    };

    p.next_token();
    p.next_token();

    Some(Expression::MethodCall(node::MethodCallData {
        receiver: None,
        ident: Box::new(ident),
        block: Block::new(Vec::new(), Vec::new()),
        args,
    }))
}
