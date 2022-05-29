use std::string::String;

pub type ConstantIndex = u16;
pub type SymbolIndex = u16;

pub trait Stringable {
    fn to_string(&self) -> String;
}

#[derive(Debug)]
pub enum Opcode {
    // OpPush pushes a constant from the constant pool onto the stack
    OpPush { index: ConstantIndex },
    // OpPop pops the topmost value from the stack
    OpPop,

    // OpGetGlobal resolves a global variable reference
    OpGetGlobal { index: SymbolIndex },
    // OpSetGlobal creates a global variable reference
    OpSetGlobal { index: SymbolIndex },

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

    OpReturn,
    OpReturnValue,
}

impl Stringable for Opcode {
    fn to_string(&self) -> String {
        format!("{:?}", self)
    }
}

pub type Bytecode = Vec<Opcode>;

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
    fn test_opcode_is_32_bits() {
        // An Opcode should be 32 bits; anything bigger and we've mis-defined some
        // variant
        assert_eq!(size_of::<Opcode>(), 4);
    }
}
