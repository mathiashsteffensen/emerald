use std::sync::Arc;

use crate::ast::node;
use crate::{ast, lexer, parser};

use crate::object::EmeraldObject;

use crate::core;

use crate::compiler::bytecode::Opcode::{
    OpAdd, OpDiv, OpFalse, OpGetGlobal, OpGreaterThan, OpGreaterThanOrEq, OpLessThan,
    OpLessThanOrEq, OpMul, OpNil, OpPop, OpPush, OpReturn, OpReturnValue, OpSend, OpSetGlobal,
    OpSub, OpTrue,
};
use crate::compiler::bytecode::{Bytecode, ConstantIndex, Opcode};
use crate::lexer::token;

pub mod bytecode;
mod symbol_table;

pub struct Compiler {
    pub bytecode: Bytecode,
    pub constant_pool: Vec<Arc<EmeraldObject>>,
    pub symbol_table: symbol_table::SymbolTable,
}

impl Compiler {
    pub fn new() -> Compiler {
        Compiler {
            bytecode: Vec::new(),
            constant_pool: Vec::with_capacity(u16::MAX as usize),
            symbol_table: symbol_table::SymbolTable::new(),
        }
    }

    pub fn compile_string(&mut self, file: String, input: String) {
        let mut parser = parser::Parser::new(lexer::input::Input::new(file, input));

        let ast = parser.parse();

        self.compile(ast);
    }

    pub fn compile(&mut self, node: ast::AST) {
        core::all::init();

        for stmt in node.statements {
            self.compile_statement(stmt);
        }
    }

    fn compile_statement(&mut self, node: node::Statement) {
        match node {
            node::Statement::ExpressionStatement(expr) => {
                self.compile_expression(expr);
                self.emit(OpPop);
            }
            node::Statement::ReturnStatement(_data, value) => match value {
                Some(expr) => {
                    self.compile_expression(expr);
                    self.emit(OpReturnValue);
                }
                None => {
                    self.emit(OpReturn);
                }
            },
        }
    }

    fn compile_expression(&mut self, node: node::Expression) {
        match node {
            node::Expression::AssignmentExpression(var, _data, val) => {
                self.compile_assignment_expression(*var, *val)
            }
            node::Expression::IdentifierExpression(data) => {
                self.compile_identifier_expression(data)
            }
            node::Expression::MethodCall(data) => self.compile_method_call(data),
            node::Expression::InfixExpression(left, data, right) => {
                self.compile_infix_expression(*left, data.literal, *right)
            }
            node::Expression::IntegerLiteral(_data, val) => self.compile_integer_literal(val),
            node::Expression::StringLiteral(data) => self.compile_string_literal(data.literal),
            node::Expression::TrueLiteral(_data) => self.emit(OpTrue),
            node::Expression::FalseLiteral(_data) => self.emit(OpFalse),
            node::Expression::NilLiteral(_data) => self.emit(OpNil),
            _ => panic!(
                "Compiler#compile_expression - no match arm defined for {:?}",
                node
            ),
        }
    }

    fn compile_identifier_expression(&mut self, data: token::TokenData) {
        if let Some(sym) = self.symbol_table.resolve(&data.literal) {
            self.emit(OpGetGlobal { index: sym.index })
        };
    }

    fn compile_assignment_expression(&mut self, var: node::Expression, val: node::Expression) {
        match &var {
            node::Expression::IdentifierExpression(data) => {
                let symbol = if let Some(sym) = self.symbol_table.resolve(&data.literal) {
                    sym
                } else {
                    self.symbol_table.define(&data.literal)
                };

                self.compile_expression(val);

                self.emit(OpSetGlobal {
                    index: symbol.index,
                });
            }
            _ => panic!(
                "target of assignment expression was not identifier, got={:?}",
                var
            ),
        }
    }

    fn compile_method_call(&mut self, data: node::MethodCallData) {
        if let Some(receiver) = data.receiver {
            self.compile_expression(*receiver);
        };

        match *data.ident {
            node::Expression::IdentifierExpression(data) => {
                let symbol = core::symbol::em_instance(data.literal);
                let index = self.push_constant(symbol);

                self.emit(OpSend { index })
            }
            _ => panic!("Method call ident was, well, not an ident ..."),
        }
    }

    fn compile_infix_expression(
        &mut self,
        left: node::Expression,
        op: String,
        right: node::Expression,
    ) {
        self.compile_expression(left);
        self.compile_expression(right);

        match op.as_str() {
            "+" => self.emit(OpAdd),
            "-" => self.emit(OpSub),
            "*" => self.emit(OpMul),
            "/" => self.emit(OpDiv),
            ">" => self.emit(OpGreaterThan),
            ">=" => self.emit(OpGreaterThanOrEq),
            "<" => self.emit(OpLessThan),
            "<=" => self.emit(OpLessThanOrEq),
            _ => panic!("Unknown operator {:?}", op),
        }
    }

    fn compile_string_literal(&mut self, val: String) {
        let obj = core::string::em_instance(val);

        self.emit_constant(obj)
    }

    fn compile_integer_literal(&mut self, val: i64) {
        let obj = core::integer::em_instance(val);

        self.emit_constant(obj)
    }

    fn emit(&mut self, op: Opcode) {
        self.bytecode.push(op);
    }

    fn push_constant(&mut self, constant: Arc<EmeraldObject>) -> ConstantIndex {
        self.constant_pool.push(constant);

        (self.constant_pool.len() - 1) as ConstantIndex
    }

    fn emit_constant(&mut self, constant: Arc<EmeraldObject>) {
        let index = self.push_constant(constant);

        self.emit(OpPush { index })
    }
}
