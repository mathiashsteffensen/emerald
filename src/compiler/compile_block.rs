use crate::ast::node;
use crate::compiler::bytecode::Opcode::{OpPop, OpReturnValue};
use crate::compiler::Compiler;
use crate::{compiler, object};

pub fn exec(c: &mut Compiler, node: node::Block) -> object::Block {
    compiler::scope::enter(c);

    for arg in &node.args {
        match arg {
            node::Expression::IdentifierExpression(data) => {
                c.symbol_table.define(&data.literal);
            }
            _ => unreachable!(),
        }
    }

    for stmt in node.body {
        c.compile_statement(stmt)
    }

    if c.check_last_op(|op| matches!(op, OpPop)) {
        let index = c.bytecode_mut().len() - 1;

        c.change_op(index, OpReturnValue);
    };

    let (bytecode, num_locals) = compiler::scope::leave(c);

    object::Block::new(bytecode, node.args.len() as u8, num_locals)
}
