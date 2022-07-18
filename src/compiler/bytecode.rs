use std::string::String;

pub type ConstantIndex = u16;
pub type SymbolIndex = u16;
pub type JumpOffset = u16;

pub trait Stringable {
    fn to_string(&self) -> String;
}

#[derive(PartialEq, Debug, Clone)]
pub enum Opcode {
    // OpPush pushes a constant from the constant pool onto the stack
    OpPush {
        index: ConstantIndex,
    },
    // OpPop pops the topmost value from the stack
    OpPop,

    // Push EM_TRUE onto the stack
    OpTrue,
    // Push EM_FALSE onto the stack
    OpFalse,
    // Push EM_NIL onto the stack
    OpNil,

    // OpGetGlobal resolves a global variable reference
    OpGetGlobal {
        index: SymbolIndex,
    },
    // OpSetGlobal creates a global variable reference
    OpSetGlobal {
        index: SymbolIndex,
    },

    // OpGetLocal resolves a local variable reference
    OpGetLocal {
        index: SymbolIndex,
    },
    // OpSetLocal creates a local variable reference
    OpSetLocal {
        index: SymbolIndex,
    },

    // Jump 'offset' forward if element at top of stack is not truthy
    OpJumpNotTruthy {
        offset: JumpOffset,
    },
    // Jump 'offset' forward no matter what
    OpJump {
        offset: JumpOffset,
    },

    // Infix operators
    OpAdd,
    OpSub,
    OpMul,
    OpDiv,
    OpEqual,
    OpNotEqual,
    OpGreaterThan,
    OpGreaterThanOrEq,
    OpLessThan,
    OpLessThanOrEq,

    // Defines a new method
    OpDefineMethod {
        proc_index: ConstantIndex,
        name_index: ConstantIndex,
    },

    // Sends a method call, name of method to call is at the index
    OpSend {
        index: ConstantIndex,
        num_args: u8,
    },

    // Set execution context
    OpSetEC,

    // Push 'self' from the current execution context onto the stack
    OpPushSelf,

    OpReturn,
    OpReturnValue,
}

impl Stringable for Opcode {
    fn to_string(&self) -> String {
        format!("{:?}", self)
    }
}

pub type Bytecode = Vec<Opcode>;

trait MyClone {
    fn clone(&self) -> Self;
}

impl Stringable for Bytecode {
    fn to_string(&self) -> String {
        let mut out = String::new();

        for (i, op) in self.iter().enumerate() {
            out.push_str(&*format!("{:04} ", i));
            out.push_str(&*op.to_string());
            out.push_str("\n")
        }

        out
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::mem::size_of;

    #[test]
    fn test_opcode_is_6_bytes() {
        // An Opcode should be 6 bytes; anything bigger and we've mis-defined some
        // variant
        assert_eq!(size_of::<Opcode>(), 6);
    }
}
