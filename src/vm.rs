use std::sync::Arc;

use crate::compiler::bytecode::{Bytecode, ConstantIndex, JumpOffset, Opcode, SymbolIndex};
use crate::core::nil_class::EM_NIL;
use crate::core::object::EM_MAIN_OBJ;
use crate::object::{Block, EmeraldMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};
use crate::vm::frame::Frame;
use crate::{compiler, core, kernel, lexer, parser};

mod frame;

const STACK_SIZE: u16 = 2048;
const MAX_FRAMES: u16 = 1024;

pub struct VM {
    execution_context: Arc<ExecutionContext>,
    locals: Vec<Arc<EmeraldObject>>,
    stack: Vec<Arc<EmeraldObject>>,
    pub sp: u16, // Always points to the next value. Top of stack is stack[sp-1]
    frames: Vec<Frame>,
    fp: u16, // Points to current frame
}

impl VM {
    pub fn new(bytecode: Bytecode, locals: Vec<Arc<EmeraldObject>>) -> VM {
        let base_frame = Frame::new(Block::new(bytecode, 0, locals.len() as u16));
        let mut frames = Vec::with_capacity(MAX_FRAMES as usize);
        frames.push(base_frame);

        VM {
            execution_context: Arc::from(ExecutionContext::new(Arc::clone(&EM_MAIN_OBJ))),
            stack: Vec::with_capacity(STACK_SIZE as usize),
            locals,
            sp: 0,
            frames,
            fp: 0,
        }
    }

    pub fn interpret(
        file_name: String,
        content: String,
    ) -> Result<(compiler::Compiler, VM), String> {
        let input = lexer::input::Input::new(file_name, content);
        let mut parser = parser::Parser::new(input);
        let ast = parser.parse();

        if parser.errors.len() != 0 {
            return Err(parser.errors.get(0).cloned().unwrap());
        }

        let mut compiler = compiler::Compiler::new();
        compiler.compile(ast);

        let mut vm = VM::new(compiler.bytecode().clone(), Default::default());
        if let Err(e) = vm.run() {
            return Err(e);
        }

        Ok((compiler, vm))
    }

    pub fn set_bytecode_and_locals(&mut self, bytecode: Bytecode, locals: Vec<Arc<EmeraldObject>>) {
        let base_frame = Frame::new(Block::new(bytecode, 0, locals.len() as u16));
        let mut frames = Vec::with_capacity(MAX_FRAMES as usize);
        frames.push(base_frame);
        self.frames = frames
    }

    pub fn run(&mut self) -> Result<(), String> {
        self.run_with_done_callback::<Box<dyn Fn(&Opcode, &mut VM) -> bool>>(&mut None)
    }

    fn run_until_return(&mut self) -> Result<(), String> {
        self.run_with_done_callback(&mut Some(|op: &Opcode, vm: &mut VM| match op {
            Opcode::OpReturnValue => {
                vm.execute_op_return_value();
                vm.pop_frame();
                true
            }
            Opcode::OpReturn => {
                vm.execute_op_return();
                vm.pop_frame();
                true
            }
            _ => false,
        }))
    }

    fn run_with_done_callback<F>(&mut self, done: &mut Option<F>) -> Result<(), String>
    where
        F: FnMut(&Opcode, &mut VM) -> bool,
    {
        while let Some(op) = self.current_frame().fetch() {
            if let Some(done_callback) = done {
                if done_callback(&op, self) {
                    return Ok(());
                }
            }

            match op {
                Opcode::OpPush { index } => {
                    let obj = kernel::get_const(index as usize).unwrap();

                    self.push(obj)
                }
                Opcode::OpPushSelf => self.push(self.get_ec().borrow_self()),
                Opcode::OpPop => {
                    self.pop();
                }
                Opcode::OpReturnValue => {
                    self.execute_op_return_value();

                    if self.is_base_frame() {
                        return Ok(());
                    }
                }
                Opcode::OpReturn => {
                    self.execute_op_return();

                    if self.is_base_frame() {
                        return Ok(());
                    }
                }
                Opcode::OpDefineMethod {
                    name_index,
                    proc_index,
                } => self.execute_op_define_method(name_index, proc_index),
                Opcode::OpSend { index, num_args } => self.execute_op_send(index, num_args),
                Opcode::OpJumpNotTruthy { offset } => self.execute_op_jump_not_truthy(offset),
                Opcode::OpJump { offset } => self.execute_jump(offset),
                Opcode::OpSetGlobal { index } => self.execute_op_set_global(index),
                Opcode::OpGetGlobal { index } => self.execute_op_get_global(index),
                Opcode::OpGetLocal { index } => {
                    self.push(self.locals.get(index as usize).cloned().unwrap())
                }
                Opcode::OpTrue => self.push(Arc::clone(&core::true_class::EM_TRUE)),
                Opcode::OpFalse => self.push(Arc::clone(&core::false_class::EM_FALSE)),
                Opcode::OpNil => self.push(Arc::clone(&core::nil_class::EM_NIL)),
                Opcode::OpAdd => self.execute_infix_operator("+"),
                Opcode::OpSub => self.execute_infix_operator("-"),
                Opcode::OpMul => self.execute_infix_operator("*"),
                Opcode::OpDiv => self.execute_infix_operator("/"),
                Opcode::OpGreaterThan => self.execute_infix_operator(">"),
                Opcode::OpGreaterThanOrEq => self.execute_infix_operator(">="),
                Opcode::OpLessThan => self.execute_infix_operator("<"),
                Opcode::OpLessThanOrEq => self.execute_infix_operator("<="),
                Opcode::OpEqual => self.execute_infix_operator("=="),
                _ => return Err(format!("Opcode not yet implemented {:?}", op)),
            }
        }

        Ok(())
    }

    fn execute_op_define_method(&mut self, name_index: ConstantIndex, proc_index: ConstantIndex) {
        let name = kernel::get_const(name_index as usize).unwrap();
        let proc = kernel::get_const(proc_index as usize).unwrap();

        match &name.underlying_value {
            UnderlyingValueType::Symbol(name_str) => match &proc.underlying_value {
                UnderlyingValueType::Proc(block) => {
                    let ec = self.get_ec();

                    ec.q_self.define_method(name_str.clone(), block.clone());

                    self.push(name);
                }
                _ => unreachable!(),
            },
            _ => unreachable!(),
        }
    }

    fn execute_op_jump_not_truthy(&mut self, offset: JumpOffset) {
        let val = self.pop();

        if !core::em_is_truthy(&val) {
            self.execute_jump(offset)
        }
    }

    fn execute_jump(&mut self, offset: JumpOffset) {
        self.current_frame().cp += offset as u64
    }

    fn execute_op_set_global(&mut self, index: SymbolIndex) {
        let val = self.stack_top();

        kernel::set_global(index as usize, val)
    }

    fn execute_op_get_global(&mut self, index: SymbolIndex) {
        let val =
            kernel::get_global(index as usize).unwrap_or(Arc::clone(&core::nil_class::EM_NIL));

        self.push(val);
    }

    fn execute_infix_operator(&mut self, op: &str) {
        let right = self.pop();
        let left = self.pop();

        self.execute_method_call(left, op, Vec::from([right]));
    }

    fn execute_op_return(&mut self) {
        self.push(Arc::clone(&EM_NIL));
        self.pop();
    }

    fn execute_op_return_value(&mut self) {
        self.pop();
    }

    fn execute_op_send(&mut self, index: ConstantIndex, num_args: u8) {
        let receiver = self.pop();
        let mut args = Vec::with_capacity(num_args as usize);

        for _ in 0..num_args {
            args.push(self.pop())
        }

        args.reverse();

        let method = kernel::get_const(index as usize).unwrap();

        match &method.underlying_value {
            UnderlyingValueType::Symbol(method_name) => {
                self.execute_method_call(receiver, method_name.as_str(), args)
            }
            _ => panic!(
                "Method name was not symbol, got={:?}",
                method.underlying_value
            ),
        }
    }

    fn execute_method_call(
        &mut self,
        receiver: Arc<EmeraldObject>,
        method_name: &str,
        args: Vec<Arc<EmeraldObject>>,
    ) {
        let method = receiver.method(method_name);

        let result = match method {
            Some(method) => match method {
                EmeraldMethod::BuiltIn(method) => method(
                    Arc::from(ExecutionContext::with_outer(receiver, self.get_ec())),
                    args,
                ),
                EmeraldMethod::Compiled(block) => {
                    self.push_frame(Frame::new(block));
                    self.locals = args;
                    self.run_until_return().unwrap();
                    self.last_popped_stack_object()
                }
            },
            None => panic!(
                "TODO! handle NoMethodError for {} on {}",
                method_name,
                receiver
                    .send(receiver.clone(), "to_s", self.get_ec(), Default::default())
                    .unwrap_or(receiver)
                    .underlying_value
            ),
        };

        self.push(result);
    }

    fn push(&mut self, obj: Arc<EmeraldObject>) {
        self.stack.insert(self.sp as usize, obj);
        self.sp += 1;
    }

    fn pop(&mut self) -> Arc<EmeraldObject> {
        if self.sp == 0 {
            return Arc::clone(&EM_NIL);
        }

        let obj = self.stack_top();

        self.sp -= 1;

        obj
    }

    fn current_frame(&mut self) -> &mut Frame {
        &mut self.frames[self.fp as usize]
    }

    fn is_base_frame(&self) -> bool {
        self.fp == 0
    }

    fn pop_frame(&mut self) {
        self.frames.pop();
        self.fp -= 1
    }

    fn push_frame(&mut self, frame: Frame) {
        self.frames.push(frame);
        self.fp += 1
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
