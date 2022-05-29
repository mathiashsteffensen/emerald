use crate::lexer::token::TokenData;
use std::string::String;

pub trait Node {
    fn token_literal(&self) -> String;
    fn to_string(&self) -> String;
}

#[derive(PartialEq, Debug, Clone)]
pub enum Statement {
    ReturnStatement(TokenData, Option<Expression>),
    ExpressionStatement(Expression),
}

impl Node for Statement {
    fn token_literal(&self) -> String {
        match self {
            Statement::ReturnStatement(data, _expr) => data.literal.to_string(),
            Statement::ExpressionStatement(expr) => expr.token_literal(),
        }
    }

    fn to_string(&self) -> String {
        match self {
            Statement::ReturnStatement(data, expr) => {
                let mut out = data.literal.to_string();

                match expr {
                    Some(expr) => {
                        out.push_str(" ");
                        out.push_str(expr.to_string().as_str());
                    }
                    None => {}
                }

                out.push_str(";");

                out
            }
            Statement::ExpressionStatement(expr) => expr.to_string(),
        }
    }
}

pub type ExpressionList = Vec<Expression>;
pub type StatementList = Vec<Statement>;

#[derive(PartialEq, Debug, Clone)]
pub struct Block {
    args: ExpressionList,
    body: StatementList,
}

impl Block {
    pub fn new(args: ExpressionList, body: StatementList) -> Block {
        Block { args, body }
    }
}

#[derive(PartialEq, Debug, Clone)]
pub struct MethodCallData {
    pub ident: Box<Expression>,
    pub args: ExpressionList,
    pub block: Block,
}

#[derive(PartialEq, Debug, Clone)]
pub enum Expression {
    // Expressions
    InfixExpression(Box<Expression>, TokenData, Box<Expression>),
    PrefixExpression(TokenData, Box<Expression>),
    IdentifierExpression(TokenData),
    AssignmentExpression(Box<Expression>, TokenData, Box<Expression>),
    MethodCall(MethodCallData),

    // Literals, which are also expressions
    IntegerLiteral(TokenData, i64),
    FloatLiteral(TokenData, f64),
    StringLiteral(TokenData),
    NilLiteral(TokenData),
    MethodLiteral(TokenData, Box<Expression>, Vec<Expression>, Vec<Statement>),
    ClassLiteral(TokenData, Box<Expression>, Vec<Statement>),
}

impl Node for Expression {
    fn token_literal(&self) -> String {
        match self {
            Expression::InfixExpression(_left, data, _right) => data.literal.to_string(),
            Expression::PrefixExpression(data, _expr) => data.literal.to_string(),
            Expression::IdentifierExpression(data) => data.literal.to_string(),
            Expression::AssignmentExpression(_name, data, _val) => data.literal.to_string(),
            Expression::MethodCall(data) => data.ident.token_literal(),
            Expression::IntegerLiteral(data, _val) => data.literal.to_string(),
            Expression::FloatLiteral(data, _val) => data.literal.to_string(),
            Expression::StringLiteral(data) => data.literal.to_string(),
            Expression::NilLiteral(data) => data.literal.to_string(),
            Expression::MethodLiteral(data, _name, _args, _body) => data.literal.to_string(),
            Expression::ClassLiteral(data, _name, _body) => data.literal.to_string(),
        }
    }

    fn to_string(&self) -> String {
        match self {
            Expression::InfixExpression(left, data, right) => {
                let mut out = "(".to_string();

                out.push_str(left.to_string().as_str());
                out.push_str(" ");
                out.push_str(data.literal.as_str());
                out.push_str(" ");
                out.push_str(right.to_string().as_str());
                out.push_str(")");

                out
            }
            Expression::PrefixExpression(data, expr) => {
                let mut out = "(".to_string();

                out.push_str(data.literal.as_str());
                out.push_str(expr.to_string().as_str());
                out.push_str(")");

                out
            }
            Expression::IdentifierExpression(data) => data.literal.to_string(),
            Expression::AssignmentExpression(name, data, val) => {
                let mut out = name.to_string();

                out.push_str(" ");
                out.push_str(data.literal.as_str());
                out.push_str(" ");
                out.push_str(val.to_string().as_str());
                out.push_str(";");

                out
            }
            Expression::MethodCall(data) => {
                let mut out = data.ident.to_string();

                if *&data.args.len() != 0 as usize {
                    out.push_str("(");
                    for arg in &data.args {
                        out.push_str(&*arg.to_string())
                    }
                    out.push_str(")")
                }

                out.push_str("\n");

                if *&data.block.body.len() != 0 as usize {
                    out.push_str("do");

                    if *&data.block.args.len() != 0 as usize {
                        out.push_str("|");
                        for arg in &data.block.args {
                            out.push_str(&*arg.to_string())
                        }
                        out.push_str("|")
                    }

                    out.push_str("\n");

                    for stmt in &data.block.body {
                        out.push_str(&*stmt.to_string());
                        out.push_str("\n");
                    }

                    out.push_str("end");
                }

                out
            }
            Expression::IntegerLiteral(data, _val) => data.literal.to_string(),
            Expression::FloatLiteral(data, _val) => data.literal.to_string(),
            Expression::StringLiteral(data) => data.literal.to_string(),
            Expression::NilLiteral(data) => data.literal.to_string(),
            Expression::MethodLiteral(_data, name, args, body) => {
                let mut out = "def ".to_string();

                out.push_str(name.to_string().as_str());
                out.push_str("(");

                let mut arg_strings: Vec<String> = Vec::new();
                for arg in args {
                    arg_strings.push(arg.to_string())
                }

                out.push_str(arg_strings.join(", ").as_str());
                out.push_str(")\n");

                for stmt in body {
                    out.push_str(stmt.to_string().as_str());
                    out.push_str("\n");
                }

                out.push_str("end");

                out
            }
            Expression::ClassLiteral(data, name, body) => {
                let mut out = data.literal.to_string();

                out.push_str(" ");
                out.push_str(&*name.to_string());
                out.push_str("\n");

                for stmt in body {
                    out.push_str(&*stmt.to_string());
                    out.push_str("\n");
                }

                out.push_str("end");

                out
            }
        }
    }
}
