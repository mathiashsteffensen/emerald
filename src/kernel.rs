use lazy_static::lazy_static;
use log::debug;
use std::sync::{Arc, Mutex};

use crate::compiler::Compiler;
use crate::object::EmeraldObject;
use crate::vm::VM;
use crate::{lexer, parser};

const CONSTANT_POOL_SIZE: u16 = u16::MAX;
const GLOBALS_SIZE: u16 = u16::MAX;

lazy_static! {
    pub static ref CONSTANT_POOL: Mutex<Vec<Arc<EmeraldObject>>> =
        Mutex::new(Vec::with_capacity(CONSTANT_POOL_SIZE as usize));
    pub static ref GLOBALS: Mutex<Vec<Arc<EmeraldObject>>> =
        Mutex::new(Vec::with_capacity(GLOBALS_SIZE as usize));
    pub static ref EMERALD_VM: Mutex<VM> = Mutex::new(VM::new());
}

pub fn execute(file_name: String, content: String) -> Result<Arc<EmeraldObject>, String> {
    let input = lexer::input::Input::new(file_name, content);
    let mut parser = parser::Parser::new(input);
    let ast = parser.parse();

    if parser.errors.len() != 0 {
        return Err(parser.errors.get(0).cloned().unwrap());
    }

    let mut c = Compiler::new();
    c.compile(ast);

    let mut vm = EMERALD_VM.lock().unwrap();
    vm.set_bytecode(c.bytecode().clone());
    vm.run()?;

    Ok(vm.last_popped_stack_object())
}

pub fn execute_method_call(
    receiver: Arc<EmeraldObject>,
    method_name: &str,
    args: Vec<Arc<EmeraldObject>>,
) -> Arc<EmeraldObject> {
    debug!("Trying to acquire lock on main VM");
    let mut vm = EMERALD_VM.lock().unwrap();
    debug!("Lock acquired");
    let num_args = *&args.len() as u8;
    for arg in args {
        vm.push(arg)
    }
    vm.execute_method_call(receiver, method_name, num_args);
    vm.current_fiber.pop()
}

pub fn push_const(obj: Arc<EmeraldObject>) -> usize {
    let mut pool = CONSTANT_POOL
        .lock()
        .expect("failed to acquire lock on constant pool");
    pool.push(obj);

    pool.len() - 1
}

pub fn get_const(index: usize) -> Option<Arc<EmeraldObject>> {
    CONSTANT_POOL.lock().unwrap().get(index).cloned()
}

pub fn reset_consts() {
    *CONSTANT_POOL.lock().unwrap() = Vec::new()
}

pub fn set_global(index: usize, obj: Arc<EmeraldObject>) {
    let mut globals = GLOBALS.lock().expect("failed to acquire lock on globals");
    globals.insert(index, obj)
}

pub fn get_global(index: usize) -> Option<Arc<EmeraldObject>> {
    GLOBALS.lock().unwrap().get(index).cloned()
}

pub fn reset_globals() {
    *GLOBALS.lock().unwrap() = Vec::new()
}
