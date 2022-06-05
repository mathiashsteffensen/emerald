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

            let consequence = parse_statement_list::exec(p, |token| {
                // Parse until 'else' or 'end' token
                matches!(token, token::Token::End(_)) || matches!(token, token::Token::Else(_))
            });

            let alternative = match &p.current_token {
                token::Token::End(_) => None,
                // Otherwise it's an else and we need to parse an alternative
                _ => {
                    p.next_token();

                    let alternative = parse_statement_list::exec(p, |token| {
                        // Parse until 'end' token
                        matches!(token, token::Token::End(_))
                    });

                    Some(alternative)
                }
            };

            Some(Expression::IfExpression(IfExpressionData {
                token: data,
                condition: Box::new(condition),
                consequence,
                alternative,
            }))
        }
        None => None,
    }
}
