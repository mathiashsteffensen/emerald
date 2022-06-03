use crate::ast;
use crate::ast::node;
use crate::ast::node::Statement;
use crate::ast::node::{Block, Expression};
use crate::lexer;
use crate::lexer::input::Input;
use crate::lexer::token;

mod parse_ast;
mod parse_expression;
mod parse_expression_statement;
mod parse_if_expression;
mod parse_return_statement;
mod parse_statement;
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

    fn parse_as_infix(&mut self, left: Expression) -> Option<Expression> {
        let tok = self.peek_token.clone();

        let expr = match tok {
            token::Token::Plus(data) => self.parse_infix_expression(data, left),
            token::Token::Minus(data) => self.parse_infix_expression(data, left),
            token::Token::Asterisk(data) => self.parse_infix_expression(data, left),
            token::Token::Slash(data) => self.parse_infix_expression(data, left),
            token::Token::LessThan(data) => self.parse_infix_expression(data, left),
            token::Token::LessThanOrEquals(data) => self.parse_infix_expression(data, left),
            token::Token::GreaterThan(data) => self.parse_infix_expression(data, left),
            token::Token::GreaterThanOrEquals(data) => self.parse_infix_expression(data, left),
            token::Token::Equals(data) => self.parse_infix_expression(data, left),
            token::Token::NotEquals(data) => self.parse_infix_expression(data, left),
            token::Token::Dot(_data) => self.parse_method_call_with_receiver(left),
            _ => None,
        };

        expr
    }

    fn parse_as_prefix(&mut self) -> Option<Expression> {
        let tok = self.current_token.clone();

        match tok {
            token::Token::Bang(data) => self.parse_prefix_expression(data),
            token::Token::Minus(data) => self.parse_prefix_expression(data),
            token::Token::Ident(data) => self.parse_identifier_expression(data),
            token::Token::Def(data) => self.parse_method_literal(data),
            token::Token::If(data) => parse_if_expression::exec(self, data),
            token::Token::Class(data) => self.parse_class_literal(data),
            token::Token::Int(data) => self.parse_integer_literal(data),
            token::Token::Float(data) => self.parse_float_literal(data),
            token::Token::String(data) => self.parse_string_literal(data),
            token::Token::True(data) => self.parse_true_literal(data),
            token::Token::False(data) => self.parse_false_literal(data),
            token::Token::Nil(data) => self.parse_nil_literal(data),
            token::Token::LeftParen(_data) => self.parse_grouped_expression(),
            _ => {
                self.add_error(&*format!(
                    "No prefix parse function found for token {:?}",
                    self.current_token
                ));

                None
            }
        }
    }

    fn parse_prefix_expression(&mut self, data: token::TokenData) -> Option<Expression> {
        self.next_token();

        let right = parse_expression::exec(self, precedence::PREFIX);

        match right {
            Some(expr) => Some(Expression::PrefixExpression(data, Box::new(expr))),
            None => None,
        }
    }

    fn parse_grouped_expression(&mut self) -> Option<Expression> {
        self.next_token();

        let expr = parse_expression::exec(self, precedence::LOWEST);

        match self.peek_token.clone() {
            token::Token::RightParen(_data) => {
                self.next_token();
                expr
            }
            _ => {
                self.add_syntax_error(self.peek_token.data(), ")");

                None
            }
        }
    }

    fn parse_identifier_expression(&mut self, data: token::TokenData) -> Option<Expression> {
        let ident = Expression::IdentifierExpression(data);

        match self.peek_token.clone() {
            token::Token::Assign(data) => self.parse_assignment_expression(ident, data),
            token::Token::LeftParen(_data) => self.parse_method_call_without_receiver(ident),
            token::Token::Dot(_data) => self.parse_method_call_with_receiver(ident),
            _ => Some(ident),
        }
    }

    fn parse_assignment_expression(
        &mut self,
        ident: Expression,
        data: token::TokenData,
    ) -> Option<Expression> {
        self.next_token();
        self.next_token();

        let val = parse_expression::exec(self, precedence::LOWEST);

        match val {
            Some(expr) => Some(Expression::AssignmentExpression(
                Box::new(ident),
                data,
                Box::new(expr),
            )),
            None => None,
        }
    }

    fn parse_method_call_without_receiver(&mut self, ident: Expression) -> Option<Expression> {
        let args = if let Some(expr) = self.parse_paren_delimited_expression_list() {
            expr
        } else {
            return None;
        };

        self.next_token();
        self.next_token();

        Some(Expression::MethodCall(node::MethodCallData {
            receiver: None,
            ident: Box::new(ident),
            block: Block::new(Vec::new(), Vec::new()),
            args,
        }))
    }

    fn parse_method_call_with_receiver(&mut self, receiver: Expression) -> Option<Expression> {
        self.next_token();

        let ident = match &self.peek_token {
            token::Token::Ident(data) => {
                if let Some(expr) = self.parse_identifier_expression(data.clone()) {
                    expr
                } else {
                    return None;
                }
            }
            _ => {
                self.add_syntax_error(self.peek_token.data(), "ident");
                return None;
            }
        };

        self.next_token();

        let args = if let Some(expr) = self.parse_paren_delimited_expression_list() {
            expr
        } else {
            return None;
        };

        self.next_token();
        self.next_token();

        Some(Expression::MethodCall(node::MethodCallData {
            receiver: Some(Box::new(receiver)),
            ident: Box::new(ident),
            block: Block::new(Vec::new(), Vec::new()),
            args,
        }))
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

    fn parse_integer_literal(&mut self, data: token::TokenData) -> Option<Expression> {
        let parse_result = data.literal.replace("_", "").parse::<i64>();
        match parse_result {
            Ok(val) => Some(Expression::IntegerLiteral(data, val)),
            Err(_err) => {
                self.add_error(&*format!(
                    "failed to parse integer value '{}'",
                    data.literal
                ));

                None
            }
        }
    }

    fn parse_float_literal(&mut self, data: token::TokenData) -> Option<Expression> {
        let parse_result = data.literal.parse::<f64>();
        match parse_result {
            Ok(val) => Some(Expression::FloatLiteral(data, val)),
            Err(_err) => {
                self.add_error(&*format!("failed to parse float value '{}'", data.literal));

                None
            }
        }
    }

    fn parse_method_literal(&mut self, data: token::TokenData) -> Option<Expression> {
        self.next_token();

        if let token::Token::Ident(ident_data) = self.current_token.clone() {
            let name = Expression::IdentifierExpression(ident_data);

            let args = if let Some(expr) = self.parse_paren_delimited_expression_list() {
                expr
            } else {
                return None;
            };

            self.next_token();
            self.next_token();

            let body = self.parse_do_end_block();

            Some(Expression::MethodLiteral(data, Box::new(name), args, body))
        } else {
            None
        }
    }

    fn parse_class_literal(&mut self, data: token::TokenData) -> Option<Expression> {
        self.next_token();

        if let token::Token::Ident(ident_data) = self.current_token.clone() {
            let name = Expression::IdentifierExpression(ident_data);

            self.next_token();

            let body = self.parse_do_end_block();

            Some(Expression::ClassLiteral(data, Box::new(name), body))
        } else {
            None
        }
    }

    fn parse_do_end_block(&mut self) -> Vec<Statement> {
        let mut statements: Vec<Statement> = Vec::new();
        let mut is_end = false;

        while !is_end {
            match &self.current_token {
                token::Token::End(_data) => is_end = true,
                token::Token::Newline(_data) => {}
                token::Token::Semicolon(_data) => {}
                token::Token::Eof(data) => {
                    let mut cloned_data = data.clone();

                    cloned_data.literal = EOF_LITERAL.to_string();

                    self.add_syntax_error(cloned_data, "end");
                    break;
                }
                _ => {
                    let statement = parse_statement::exec(self);
                    if let Some(stmt) = statement {
                        statements.push(stmt)
                    }
                }
            }

            self.next_token();
        }

        statements
    }

    fn parse_infix_expression(
        &mut self,
        data: token::TokenData,
        left: Expression,
    ) -> Option<Expression> {
        let precedence = self.peek_precedence();

        self.next_token();
        self.next_token();

        let right = parse_expression::exec(self, precedence);

        match right {
            Some(right_expr) => Some(Expression::InfixExpression(
                Box::new(left),
                data,
                Box::new(right_expr),
            )),
            None => None,
        }
    }

    fn parse_paren_delimited_expression_list(&mut self) -> Option<Vec<Expression>> {
        if matches!(&self.peek_token, lexer::token::Token::LeftParen(_data)) {
            self.next_token();

            if let Some(expressions) = self.parse_expression_list() {
                self.next_token();

                Some(expressions)
            } else {
                None
            }
        } else {
            Some(Vec::new())
        }
    }

    fn parse_expression_list(&mut self) -> Option<Vec<Expression>> {
        let mut args: Vec<Expression> = Vec::new();

        self.next_token();

        if let Some(expr) = parse_expression::exec(self, precedence::LOWEST) {
            args.push(expr);
        } else {
            return None;
        };

        let mut is_comma = matches!(&self.peek_token, lexer::token::Token::Comma(_data));

        while is_comma {
            self.next_token();
            self.next_token();

            if let Some(expr) = parse_expression::exec(self, precedence::LOWEST) {
                args.push(expr);
            } else {
                return None;
            };

            if self.peek_token_is_newline() {
                self.next_token();
            }

            is_comma = matches!(&self.peek_token, lexer::token::Token::Comma(_data));
        }

        Some(args)
    }

    fn peek_precedence(&mut self) -> i16 {
        precedence::precedence_for(self.peek_token.clone())
    }

    fn add_error(&mut self, msg: &str) {
        self.errors.push(msg.to_string())
    }

    fn add_syntax_error(&mut self, got_data: token::TokenData, expected: &str) {
        self.add_error(&*format!(
            "syntax error at line:{}:{}: expected '{}', found '{}'",
            got_data.line, got_data.column, expected, got_data.literal,
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
