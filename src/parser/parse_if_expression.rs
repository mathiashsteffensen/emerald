use crate::ast::node::{Expression, IfExpressionData};
use crate::lexer::token;
use crate::lexer::token::TokenData;
use crate::parser::{parse_expression, parse_statement_list, precedence, Parser};

pub fn exec(p: &mut Parser, data: TokenData) -> Option<Expression> {
    p.next_token();

    let condition = parse_expression::exec(p, precedence::LOWEST);
    match condition {
        Some(condition) => {
            p.next_token();
            let consequence =
                parse_statement_list::exec(p, |token| matches!(token, token::Token::End(_)));

            Some(Expression::IfExpression(IfExpressionData {
                condition: Box::new(condition),
                consequence,
                token: data,
                alternative: None,
            }))
        }
        None => None,
    }
}
