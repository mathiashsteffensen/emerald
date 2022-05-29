use std::string::String;

pub mod node;

pub struct AST {
    pub statements: Vec<node::Statement>,
}

impl node::Node for AST {
    fn token_literal(&self) -> String {
        if self.statements.len() > 0 {
            self.statements[0].token_literal()
        } else {
            "".to_string()
        }
    }

    fn to_string(&self) -> String {
        let mut out = String::new();

        for statement in &self.statements {
            out.push_str(&*statement.to_string())
        }

        out
    }
}

#[cfg(test)]
mod tests {
    use crate::ast;
    use crate::ast::node;
    use crate::ast::node::Node;
    use crate::lexer::token::TokenData;

    fn new_ast() -> ast::AST {
        let int = node::Expression::IntegerLiteral(
            TokenData {
                literal: "5".to_string(),
                pos: 0,
                line: 0,
                column: 0,
            },
            5,
        );

        let ret = node::Statement::ReturnStatement(
            TokenData {
                literal: "return".to_string(),
                pos: 0,
                line: 0,
                column: 0,
            },
            Some(int),
        );

        ast::AST {
            statements: Vec::from([ret]),
        }
    }

    #[test]
    fn token_literal() {
        assert_eq!(new_ast().token_literal(), "return");
    }

    #[test]
    fn to_string() {
        assert_eq!(new_ast().to_string(), "return 5;");
    }
}
