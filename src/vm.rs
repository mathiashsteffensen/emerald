use std::sync::Arc;

use crate::compiler::bytecode::{Bytecode, ConstantIndex, JumpOffset, Opcode, SymbolIndex};
use crate::core::nil_class::EM_NIL;
use crate::core::object::EM_MAIN_OBJ;
use crate::object::{Block, EmeraldMethod, EmeraldObject, ExecutionContext, UnderlyingValueType};
use crate::vm::fiber::Fiber;
use crate::vm::frame::Frame;
use crate::{core, kernel};

mod fiber;
mod frame;

pub struct VM {
    execution_context: Arc<ExecutionContext>,
    pub current_fiber: Fiber,
}

impl VM {
    pub fn new() -> VM {
        VM {
            execution_context: Arc::from(ExecutionContext::new(Arc::clone(&EM_MAIN_OBJ))),
            current_fiber: Fiber::new(),
        }
    }

    pub fn set_bytecode(&mut self, bytecode: Bytecode) {
        self.current_fiber
            .push_frame(Frame::new(Block::new(bytecode, 0, 0), 0))
    }

    pub fn run(&mut self) -> Result<(), String> {
        self.run_with_done_callback::<Box<dyn Fn(&Opcode, &mut VM) -> bool>>(&mut None)
    }

    fn run_until_return(&mut self) -> Result<(), String> {
        self.run_with_done_callback(&mut Some(|op: &Opcode, vm: &mut VM| match op {
            Opcode::OpReturnValue => {
                vm.execute_op_return_value();
                true
            }
            Opcode::OpReturn => {
                vm.execute_op_return();
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
                    let is_base = self.is_base_frame();
                    self.execute_op_return_value();

                    if is_base {
                        self.pop();
                        return Ok(());
                    }
                }
                Opcode::OpReturn => {
                    let is_base = self.is_base_frame();
                    self.execute_op_return();

                    if is_base {
                        self.pop();
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
                    let value = self.current_fiber.get_local(index as usize);
                    self.push(value)
                }
                Opcode::OpSetLocal { index } => self.execute_op_set_local(index),
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

    fn execute_op_set_local(&mut self, index: SymbolIndex) {
        let object = self.stack_top();

        self.current_fiber.insert_local(index as usize, object)
    }

    fn execute_infix_operator(&mut self, op: &str) {
        let right = self.pop();
        let left = self.pop();
        self.push(right); // TODO: This is dumb

        self.execute_method_call(left, op, 1);
    }

    fn execute_op_return(&mut self) {
        self.pop_frame();
        self.push(Arc::clone(&EM_NIL));
    }

    fn execute_op_return_value(&mut self) {
        let obj = self.pop();
        self.pop_frame();
        self.push(obj)
    }

    fn execute_op_send(&mut self, index: ConstantIndex, num_args: u8) {
        let receiver = self.pop();

        let method = kernel::get_const(index as usize).unwrap();

        match &method.underlying_value {
            UnderlyingValueType::Symbol(method_name) => {
                self.execute_method_call(receiver, method_name.as_str(), num_args)
            }
            _ => panic!(
                "Method name was not symbol, got={:?}",
                method.underlying_value
            ),
        }
    }

    pub fn execute_method_call(
        &mut self,
        receiver: Arc<EmeraldObject>,
        method_name: &str,
        num_args: u8,
    ) {
        let method = receiver.method(method_name);

        match method {
            Some(method) => match method {
                EmeraldMethod::BuiltIn(method) => {
                    let mut args = Vec::with_capacity(num_args as usize);

                    for _ in 0..num_args {
                        args.insert(0, self.pop())
                    }

                    let result = method(
                        Arc::from(ExecutionContext::with_outer(receiver, self.get_ec())),
                        args,
                    );

                    self.push(result)
                }
                EmeraldMethod::Compiled(block) => {
                    let base_pointer = if self.current_fiber.sp >= block.num_locals {
                        self.current_fiber.sp - block.num_locals
                    } else {
                        self.current_fiber.sp
                    };
                    self.push_frame(Frame::new(block, base_pointer));
                    self.run_until_return().unwrap();
                }
            },
            None => {
                panic!(
                    "TODO! handle NoMethodError for {} on {}",
                    method_name,
                    receiver.class_name(),
                )
            }
        };
    }

    pub(crate) fn push(&mut self, obj: Arc<EmeraldObject>) {
        self.current_fiber.push(obj)
    }

    fn pop(&mut self) -> Arc<EmeraldObject> {
        self.current_fiber.pop()
    }

    fn current_frame(&mut self) -> &mut Frame {
        let frame = self.current_fiber.current_frame();

        frame
    }

    fn is_base_frame(&self) -> bool {
        self.current_fiber.is_base_frame()
    }

    fn pop_frame(&mut self) {
        self.current_fiber.pop_frame()
    }

    fn push_frame(&mut self, frame: Frame) {
        self.current_fiber.push_frame(frame)
    }

    pub fn stack_top(&self) -> Arc<EmeraldObject> {
        self.current_fiber.stack_top()
    }

    pub fn last_popped_stack_object(&mut self) -> Arc<EmeraldObject> {
        self.current_fiber.last_popped_stack_object()
    }

    pub fn get_ec(&self) -> Arc<ExecutionContext> {
        Arc::clone(&self.execution_context)
    }
}
