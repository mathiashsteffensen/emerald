use std::collections::HashMap;
use std::rc::Rc;

use crate::ast;
use crate::ast::node;

use crate::object::EmeraldObject;

use crate::core;

use crate::compiler::bytecode::Opcode::{
    OpAdd, OpMul, OpPop, OpPush, OpReturn, OpReturnValue, OpSub,
};
use crate::compiler::bytecode::{Bytecode, ConstantIndex, Opcode};

pub mod bytecode;
mod symbol_table;

pub struct Compiler {
    pub bytecode: Bytecode,
    pub constant_pool: Vec<Rc<EmeraldObject>>,
    pub built_ins: HashMap<String, Rc<EmeraldObject>>,
}

impl Compiler {
    pub fn new() -> Compiler {
        Compiler {
            bytecode: Vec::new(),
            constant_pool: Vec::with_capacity(u16::MAX as usize),
            built_ins: core::all::map(),
        }
    }

    pub fn compile(&mut self, node: ast::AST) {
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
            node::Expression::InfixExpression(left, data, right) => {
                self.compile_infix_expression(*left, data.literal, *right)
            }
            node::Expression::IntegerLiteral(_data, val) => self.compile_integer_literal(val),
            node::Expression::StringLiteral(data) => self.compile_string_literal(data.literal),
            _ => {}
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
            _ => panic!("Unknown operator {:?}", op),
        }
    }

    fn compile_string_literal(&mut self, val: String) {
        let obj = core::string::em_instance(self.get_built_in(core::string::NAME), val);

        self.emit_constant(obj)
    }

    fn compile_integer_literal(&mut self, val: i64) {
        let obj = core::integer::em_instance(self.get_built_in(core::integer::NAME), val);

        self.emit_constant(obj)
    }

    fn emit(&mut self, op: Opcode) {
        self.bytecode.push(op);
    }

    fn emit_constant(&mut self, constant: Rc<EmeraldObject>) {
        self.constant_pool.push(constant);

        self.emit(OpPush {
            index: (self.constant_pool.len() - 1) as ConstantIndex,
        })
    }

    pub fn get_built_in(&self, name: &str) -> Rc<EmeraldObject> {
        Rc::clone(self.built_ins.get(name).unwrap())
    }
}
