use std::rc::Rc;

use crate::compiler::bytecode::{Bytecode, Opcode};
use crate::object::ExecutionContext;
use crate::{compiler, lexer, object, parser};

const STACK_SIZE: u16 = 2048;
const GLOBALS_SIZE: u16 = u16::MAX;
// const MAX_FRAMES: u16 = STACK_SIZE / 2;

pub struct VM {
    execution_context: Rc<ExecutionContext>,
    constants: Vec<Rc<object::EmeraldObject>>,
    globals: Vec<Rc<object::EmeraldObject>>,
    stack: Vec<Rc<object::EmeraldObject>>,
    pub sp: u16, // Always points to the next value. Top of stack is stack[sp-1]
    bytecode: Bytecode,
    cp: i64,
}

impl VM {
    pub fn new(c: &compiler::Compiler) -> VM {
        VM {
            execution_context: Rc::from(ExecutionContext::new(
                c.get_built_in("String"),
                c.built_ins.clone(),
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
                Opcode::OpAdd => self.execute_infix_operator("+"),
                Opcode::OpSub => self.execute_infix_operator("-"),
                Opcode::OpMul => self.execute_infix_operator("*"),
                _ => return Err(format!("Opcode not yet implemented {:?}", op)),
            }
        }

        Ok(())
    }

    fn execute_infix_operator(&mut self, op: &str) {
        let right = self.pop();
        let left = self.pop();

        let method = right.method(op);

        let result = match method {
            Some(method) => method(
                Rc::from(ExecutionContext::with_outer(left, self.get_ec())),
                Vec::from([right]),
            ),
            None => panic!("TODO! handle NoMethodError for {} on {:?}", op, right),
        };

        self.push(result);
    }

    fn fetch(&mut self) -> Option<Opcode> {
        let op = self.bytecode.get(self.cp as usize);

        self.cp += 1;

        if let Some(raw_op) = op {
            Some(raw_op.clone())
        } else {
            None
        }
    }

    fn push(&mut self, obj: Rc<object::EmeraldObject>) {
        self.stack.insert(self.sp as usize, obj);
        self.sp += 1;
    }

    fn pop(&mut self) -> Rc<object::EmeraldObject> {
        let obj = self.stack_top();

        self.sp -= 1;

        obj
    }

    // fetches the object at the top of the stack
    pub fn stack_top(&self) -> Rc<object::EmeraldObject> {
        Rc::clone(self.stack.get((self.sp - 1) as usize).unwrap())
    }

    // fetches the last object that was popped off the stack
    pub fn last_popped_stack_object(&mut self) -> Rc<object::EmeraldObject> {
        Rc::clone(self.stack.get(self.sp as usize).unwrap())
    }

    pub fn get_ec(&self) -> Rc<ExecutionContext> {
        Rc::clone(&self.execution_context)
    }
}
