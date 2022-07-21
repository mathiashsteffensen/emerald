use crate::compiler::bytecode::Opcode;
use crate::object::Block;

// Represents a call frame to be executed
pub struct Frame {
    // The block being called
    pub(crate) block: Block,
    // Code pointer, always points to index of next Opcode to fetch
    pub(crate) cp: u64,
    pub(crate) base_sp: u16,
}

impl Frame {
    pub(crate) fn new(block: Block, base_sp: u16) -> Frame {
        Frame {
            block,
            cp: 0,
            base_sp,
        }
    }

    pub(crate) fn fetch(&mut self) -> Option<Opcode> {
        let op = self.block.bytecode.get(self.cp as usize);

        self.cp += 1;

        op.cloned()
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use std::mem::size_of;

    #[test]
    fn test_frame_size() {
        assert_eq!(size_of::<Frame>(), 48);
    }
}
