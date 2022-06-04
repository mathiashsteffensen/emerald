use crate::ast::node::Statement;
use crate::lexer::token;
use crate::parser::{parse_statement, Parser, EOF_LITERAL};

pub fn exec<F>(p: &mut Parser, is_end_check: F) -> Vec<Statement>
where
    F: Fn(&token::Token) -> bool,
{
    let mut statements: Vec<Statement> = Vec::new();
    let mut is_end = false;

    while !is_end {
        match &p.current_token {
            token::Token::Newline(_data) => {}
            token::Token::Semicolon(_data) => {}
            token::Token::Eof(data) => {
                let mut cloned_data = data.clone();

                cloned_data.literal = EOF_LITERAL.to_string();

                p.add_syntax_error(cloned_data, "end");
                break;
            }
            _ => {
                let statement = parse_statement::exec(p);
                if let Some(stmt) = statement {
                    statements.push(stmt)
                }
            }
        }

        p.next_token();

        is_end = is_end_check(&p.current_token);
    }

    statements
}
