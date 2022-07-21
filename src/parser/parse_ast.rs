use crate::ast::AST;
use crate::lexer::token::Token;
use crate::parser::{parse_statement, Parser};

pub fn exec(p: &mut Parser) -> AST {
    let mut ast = AST {
        statements: Vec::new(),
    };

    let mut is_eof = false;

    while !is_eof {
        if let Token::Eof(_data) = &p.current_token {
            is_eof = true;
        } else {
            let statement = parse_statement::exec(p);
            if let Some(stmt) = statement {
                ast.statements.push(stmt)
            }
        }

        p.next_token();
    }

    ast
}
