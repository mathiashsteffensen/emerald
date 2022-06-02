use std::sync::Arc;

use crate::compiler::bytecode::{Bytecode, Opcode};
use crate::object::{EmeraldObject, ExecutionContext, UnderlyingValueType};
use crate::{compiler, core, lexer, parser};

const STACK_SIZE: u16 = 2048;
const GLOBALS_SIZE: u16 = u16::MAX;
// const MAX_FRAMES: u16 = STACK_SIZE / 2;

pub struct VM {
    execution_context: Arc<ExecutionContext>,
    constants: Vec<Arc<EmeraldObject>>,
    globals: Vec<Arc<EmeraldObject>>,
    stack: Vec<Arc<EmeraldObject>>,
    pub sp: u16, // Always points to the next value. Top of stack is stack[sp-1]
    bytecode: Bytecode,
    cp: i64, // Code pointer, always points to index of next Opcode to fetch
}

impl VM {
    pub fn new(c: &compiler::Compiler) -> VM {
        VM {
            execution_context: Arc::from(ExecutionContext::new(
                core::em_get_class("String").unwrap(),
            )),
            constants: c.constant_pool.clone(),
            globals: Vec::with_capacity(GLOBALS_SIZE as usize),
            stack: Vec::with_capacity(STACK_SIZE as usize),
            sp: 0,
            bytecode: c.bytecode.clone(),
            cp: 0,
        }
    }

    pub fn interpret(
        file_name: String,
        content: String,
    ) -> Result<(compiler::Compiler, VM), String> {
        let input = lexer::input::Input::new(file_name, content);
        let mut parser = parser::Parser::new(input);
        let ast = parser.parse_ast();

        if parser.errors.len() != 0 {
            return Err(parser.errors.get(0).cloned().unwrap());
        }

        let mut compiler = compiler::Compiler::new();
        compiler.compile(ast);

        let mut vm = VM::new(&compiler);
        if let Err(e) = vm.run() {
            return Err(e);
        }

        Ok((compiler, vm))
    }

    pub fn run(&mut self) -> Result<(), String> {
        while let Some(op) = self.fetch() {
            match op {
                Opcode::OpPush { index } => {
                    let obj = self.constants.get(index as usize).cloned().unwrap();

                    self.push(obj)
                }
                Opcode::OpPop => {
                    self.pop();
                }
                Opcode::OpTrue => self.push(Arc::clone(&core::true_class::EM_TRUE)),
                Opcode::OpFalse => self.push(Arc::clone(&core::false_class::EM_FALSE)),
                Opcode::OpSend { index } => {
                    let receiver = self.pop();
                    let method = self.constants.get(index as usize).cloned().unwrap();

                    match &method.underlying_value {
                        UnderlyingValueType::Symbol(method_name) => {
                            self.execute_method_call(receiver, method_name.as_str(), Vec::new())
                        }
                        _ => panic!(
                            "Method name was not symbol, got={:?}",
                            method.underlying_value
                        ),
                    }
                }
                Opcode::OpAdd => self.execute_infix_operator("+"),
                Opcode::OpSub => self.execute_infix_operator("-"),
                Opcode::OpMul => self.execute_infix_operator("*"),
                Opcode::OpDiv => self.execute_infix_operator("/"),
                _ => return Err(format!("Opcode not yet implemented {:?}", op)),
            }
        }

        Ok(())
    }

    fn execute_infix_operator(&mut self, op: &str) {
        let right = self.pop();
        let left = self.pop();

        self.execute_method_call(left, op, Vec::from([right]));
    }

    fn execute_method_call(
        &mut self,
        receiver: Arc<EmeraldObject>,
        method_name: &str,
        args: Vec<Arc<EmeraldObject>>,
    ) {
        let method = receiver.method(method_name);

        let result = match method {
            Some(method) => method(
                Arc::from(ExecutionContext::with_outer(receiver, self.get_ec())),
                args,
            ),
            None => panic!(
                "TODO! handle NoMethodError for {} on {:?}",
                method_name, receiver
            ),
        };

        self.push(result);
    }

    fn fetch(&mut self) -> Option<Opcode> {
        let op = self.bytecode.get(self.cp as usize);

        self.cp += 1;

        op.cloned()
    }

    fn push(&mut self, obj: Arc<EmeraldObject>) {
        self.stack.insert(self.sp as usize, obj);
        self.sp += 1;
    }

    fn pop(&mut self) -> Arc<EmeraldObject> {
        let obj = self.stack_top();

        self.sp -= 1;

        obj
    }

    // fetches the object at the top of the stack
    pub fn stack_top(&self) -> Arc<EmeraldObject> {
        Arc::clone(self.stack.get((self.sp - 1) as usize).unwrap())
    }

    // fetches the last object that was popped off the stack
    pub fn last_popped_stack_object(&mut self) -> Arc<EmeraldObject> {
        Arc::clone(self.stack.get(self.sp as usize).unwrap())
    }

    pub fn get_ec(&self) -> Arc<ExecutionContext> {
        Arc::clone(&self.execution_context)
    }
}
