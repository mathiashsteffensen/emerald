use crate::compiler::bytecode::{Bytecode, Opcode};

pub struct Frame {
    pub bytecode: Bytecode,
    pub cp: u64, // Code pointer, always points to index of next Opcode to fetch
}

impl Frame {
    pub fn new(bytecode: Bytecode) -> Frame {
        Frame { bytecode, cp: 0 }
    }

    pub fn fetch(&mut self) -> Option<Opcode> {
        let op = self.bytecode.get(self.cp as usize);

        self.cp += 1;

        op.cloned()
    }
}
