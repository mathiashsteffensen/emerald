use crate::ast;
use crate::ast::node::Expression;
use crate::lexer;
use crate::lexer::input::Input;
use crate::lexer::token;

mod parse_as_infix;
mod parse_as_prefix;
mod parse_assignment_expression;
mod parse_ast;
mod parse_class_literal;
mod parse_expression;
mod parse_expression_list;
mod parse_expression_statement;
mod parse_float_literal;
mod parse_grouped_expression;
mod parse_identifier_expression;
mod parse_if_expression;
mod parse_infix_expression;
mod parse_integer_literal;
mod parse_method_call_with_receiver;
mod parse_method_call_without_receiver;
mod parse_method_literal;
mod parse_paren_delimited_expression_list;
mod parse_prefix_expression;
mod parse_return_statement;
mod parse_statement;
mod parse_statement_list;
mod precedence;

const EOF_LITERAL: &str = "EOF";

pub struct Parser {
    input: Input,
    lexer: lexer::Lexer,
    pub errors: Vec<String>,
    current_token: token::Token,
    peek_token: token::Token,
}

impl Parser {
    pub fn new(input: Input) -> Parser {
        let mut lexer = lexer::Lexer::new(input.clone());

        let cur_token = lexer.next_token();
        let peek_token = lexer.next_token();

        Parser {
            lexer,
            input,
            errors: Vec::new(),
            current_token: cur_token,
            peek_token,
        }
    }

    pub fn parse(&mut self) -> ast::AST {
        parse_ast::exec(self)
    }

    fn parse_true_literal(&self, data: token::TokenData) -> Option<Expression> {
        Some(Expression::TrueLiteral(data))
    }

    fn parse_false_literal(&self, data: token::TokenData) -> Option<Expression> {
        Some(Expression::FalseLiteral(data))
    }

    fn parse_nil_literal(&self, data: token::TokenData) -> Option<Expression> {
        Some(Expression::NilLiteral(data))
    }

    fn parse_string_literal(&self, data: token::TokenData) -> Option<Expression> {
        Some(Expression::StringLiteral(data))
    }

    fn peek_precedence(&mut self) -> i16 {
        precedence::precedence_for(self.peek_token.clone())
    }

    fn add_error(&mut self, msg: &str) {
        self.errors.push(msg.to_string())
    }

    fn add_syntax_error(&mut self, got_data: token::TokenData, expected: &str) {
        self.add_error(&*format!(
            "syntax error at {}:{}:{}: expected '{}', found '{}'",
            self.input.file_name, got_data.line, got_data.column, expected, got_data.literal,
        ));
    }

    fn next_token(&mut self) {
        self.current_token = self.peek_token.clone();
        self.peek_token = self.lexer.next_token();
    }

    fn peek_token_is_semicolon(&mut self) -> bool {
        match &self.peek_token {
            token::Token::Semicolon(_data) => true,
            _ => false,
        }
    }

    fn peek_token_is_newline(&mut self) -> bool {
        match &self.peek_token {
            token::Token::Newline(_data) => true,
            _ => false,
        }
    }

    fn peek_token_is_eof(&mut self) -> bool {
        match &self.peek_token {
            token::Token::Eof(_data) => true,
            _ => false,
        }
    }

    fn next_if_semicolon_or_newline(&mut self) {
        if self.peek_token_is_newline() || self.peek_token_is_semicolon() {
            self.next_token()
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn add_error() {
        let mut parser = Parser::new(Input::new("test.rb".to_string(), "return 5;".to_string()));

        parser.add_error("this is an error");

        assert_eq!(parser.errors.len(), 1)
    }
}
