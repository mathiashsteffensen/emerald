use crate::ast::node::{Expression, MethodLiteralData};
use crate::compiler;
use crate::compiler::bytecode::Opcode::OpDefineMethod;
use crate::compiler::Compiler;
use crate::core;

pub fn exec(c: &mut Compiler, data: MethodLiteralData) {
    let name = match *data.name {
        Expression::IdentifierExpression(data) => core::symbol::em_instance(data.literal),
        _ => unreachable!(),
    };
    let name_index = c.push_constant(name);

    let block = compiler::compile_block::exec(c, data.block);
    let proc = core::proc::em_instance(block);
    let proc_index = c.push_constant(proc);

    c.emit(OpDefineMethod {
        proc_index,
        name_index,
    });
}
