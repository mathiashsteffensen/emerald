use crate::compiler::bytecode::Opcode;
use crate::object::Block;

pub struct Frame {
    pub block: Block,
    pub cp: u64, // Code pointer, always points to index of next Opcode to fetch
}

impl Frame {
    pub fn new(block: Block) -> Frame {
        Frame { block, cp: 0 }
    }

    pub fn fetch(&mut self) -> Option<Opcode> {
        let op = self.block.bytecode.get(self.cp as usize);

        self.cp += 1;

        op.cloned()
    }
}
