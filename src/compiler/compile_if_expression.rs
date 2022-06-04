use crate::ast::node::IfExpressionData;
use crate::compiler::bytecode::JumpOffset;
use crate::compiler::bytecode::Opcode::{OpJump, OpJumpNotTruthy, OpNil};
use crate::compiler::Compiler;

pub fn exec(c: &mut Compiler, data: IfExpressionData) {
    c.compile_expression(*data.condition);

    // Emit with a fake offset
    // The offset can't be known until we have compiled the consequence
    let jump_not_truthy_pos = c.emit(OpJumpNotTruthy { offset: 0 });

    for stmt in data.consequence {
        c.compile_statement(stmt);
    }

    c.remove_last_if_op_pop();

    // Emit with a fake offset
    // The offset can't be known until we have compiled the alternative
    let jump_pos = c.emit(OpJump { offset: 0 });

    let real_jump_not_truthy_offset = (c.bytecode.len() - jump_not_truthy_pos - 1) as JumpOffset;
    c.change_op(
        jump_not_truthy_pos,
        OpJumpNotTruthy {
            offset: real_jump_not_truthy_offset,
        },
    );

    match data.alternative {
        Some(statements) => {
            for stmt in statements {
                c.compile_statement(stmt);
            }
        }
        None => {
            c.emit(OpNil);
        }
    }

    let real_jump_offset = (c.bytecode.len() - jump_pos - 1) as JumpOffset;
    c.change_op(
        jump_pos,
        OpJump {
            offset: real_jump_offset,
        },
    )
}
