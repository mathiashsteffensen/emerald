use std::sync::Arc;

use crate::ast::node;
use crate::{ast, kernel, lexer, parser};

use crate::object::EmeraldObject;

use crate::core;

use crate::compiler::bytecode::Opcode::{
    OpAdd, OpDiv, OpEqual, OpFalse, OpGetGlobal, OpGetLocal, OpGreaterThan, OpGreaterThanOrEq,
    OpLessThan, OpLessThanOrEq, OpMul, OpNil, OpPop, OpPush, OpPushSelf, OpReturn, OpReturnValue,
    OpSend, OpSetGlobal, OpSetLocal, OpSub, OpTrue,
};
use crate::compiler::bytecode::{Bytecode, ConstantIndex, Opcode};
use crate::compiler::scope::CompilationScope;
use crate::lexer::token;

pub mod bytecode;
mod compile_block;
mod compile_if_expression;
mod compile_method_literal;
pub(crate) mod scope;
pub(crate) mod symbol_table;

pub struct Compiler {
    pub(crate) symbol_table: symbol_table::SymbolTable,
    pub(crate) scopes: Vec<CompilationScope>,
    pub(crate) scope_index: u8,
}

impl Compiler {
    pub fn new() -> Compiler {
        Compiler {
            symbol_table: symbol_table::SymbolTable::new(),
            scopes: Vec::from([scope::new()]),
            scope_index: 0,
        }
    }

    pub fn compile_string(&mut self, file: String, input: String) {
        let mut parser = parser::Parser::new(lexer::input::Input::new(file, input));

        let ast = parser.parse();

        self.compile(ast);
    }

    pub fn compile(&mut self, node: ast::AST) {
        core::all::init();

        for stmt in &node.statements {
            self.compile_statement(stmt.clone());
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
            node::Expression::MethodLiteral(data) => compile_method_literal::exec(self, data),
            node::Expression::InfixExpression(left, data, right) => {
                self.compile_infix_expression(*left, data.literal, *right)
            }
            node::Expression::IfExpression(data) => compile_if_expression::exec(self, data),
            node::Expression::IntegerLiteral(_data, val) => self.compile_integer_literal(val),
            node::Expression::StringLiteral(data) => self.compile_string_literal(data.literal),
            node::Expression::TrueLiteral(_data) => {
                self.emit(OpTrue);
            }
            node::Expression::FalseLiteral(_data) => {
                self.emit(OpFalse);
            }
            node::Expression::NilLiteral(_data) => {
                self.emit(OpNil);
            }
            _ => panic!(
                "Compiler#compile_expression - no match arm defined for {:?}",
                node
            ),
        }
    }

    fn compile_identifier_expression(&mut self, data: token::TokenData) {
        if let Some(sym) = self.symbol_table.resolve(&data.literal) {
            let op = match sym.scope {
                symbol_table::SymbolScope::Global => OpGetGlobal { index: sym.index },
                symbol_table::SymbolScope::Local => OpGetLocal { index: sym.index },
            };
            self.emit(op);
        } else {
            let index = self.push_constant(core::symbol::em_instance(data.literal));

            self.emit(OpPushSelf);
            self.emit(OpSend { index, num_args: 0 });
        };
    }

    fn compile_assignment_expression(&mut self, var: node::Expression, val: node::Expression) {
        match &var {
            node::Expression::IdentifierExpression(data) => {
                let sym = if let Some(sym) = self.symbol_table.resolve(&data.literal) {
                    sym
                } else {
                    self.symbol_table.define(&data.literal)
                };

                self.compile_expression(val);

                let op = match sym.scope {
                    symbol_table::SymbolScope::Global => OpSetGlobal { index: sym.index },
                    symbol_table::SymbolScope::Local => OpSetLocal { index: sym.index },
                };

                self.emit(op);
            }
            _ => panic!(
                "target of assignment expression was not identifier, got={:?}",
                var
            ),
        }
    }

    fn compile_method_call(&mut self, data: node::MethodCallData) {
        let args = &data.args;
        let num_args = args.len() as u8;
        for arg in args {
            self.compile_expression(arg.clone())
        }

        if let Some(receiver) = data.receiver {
            self.compile_expression(*receiver);
        } else {
            self.emit(OpPushSelf);
        };

        match *data.ident {
            node::Expression::IdentifierExpression(data) => {
                let symbol = core::symbol::em_instance(data.literal);
                let index = self.push_constant(symbol);

                self.emit(OpSend { index, num_args });
            }
            _ => unreachable!(),
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
            "==" => self.emit(OpEqual),
            _ => panic!("Unknown operator {:?}", op),
        };
    }

    fn compile_string_literal(&mut self, val: String) {
        let obj = core::string::em_instance(val);

        self.emit_constant(obj)
    }

    fn compile_integer_literal(&mut self, val: i64) {
        let obj = core::integer::em_instance(val);

        self.emit_constant(obj)
    }

    fn remove_last_if_op_pop(&mut self) {
        if self.check_last_op(|op| matches!(op, Opcode::OpPop)) {
            self.bytecode_mut().pop();
        }
    }

    fn check_last_op<F>(&mut self, checker: F) -> bool
    where
        F: Fn(&Opcode) -> bool,
    {
        checker(self.bytecode_mut().last().unwrap())
    }

    fn change_op(&mut self, index: usize, new: Opcode) {
        self.bytecode_mut()[index] = new
    }

    fn emit(&mut self, op: Opcode) -> usize {
        self.bytecode_mut().push(op);

        self.bytecode_mut().len() - 1
    }

    fn push_constant(&mut self, constant: Arc<EmeraldObject>) -> ConstantIndex {
        kernel::push_const(constant) as ConstantIndex
    }

    fn emit_constant(&mut self, constant: Arc<EmeraldObject>) {
        let index = self.push_constant(constant);

        self.emit(OpPush { index });
    }

    pub fn bytecode_mut(&mut self) -> &mut Bytecode {
        let scope = &mut self.scopes[self.scope_index as usize];

        &mut scope.bytecode
    }

    pub fn bytecode(&self) -> &Bytecode {
        let scope = &self.scopes[self.scope_index as usize];

        &scope.bytecode
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_compile_string() {
        let mut compiler = Compiler::new();

        compiler.compile_string("test.rb".to_string(), "2".to_string());

        assert_eq!(compiler.bytecode_mut()[0], OpPush { index: 0 })
    }
}
