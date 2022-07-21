use crate::core::nil_class::EM_NIL;
use crate::object::EmeraldObject;
use crate::vm::frame::Frame;
use std::sync::Arc;

const STACK_SIZE: u16 = 2048;
const MAX_FRAMES: u16 = 1024;

// A fiber is essentially an abstract thread not managed by the OS but by emerald instead
// This is to allow for concurrent execution in the future, and will allow us to implement Ruby's Fiber class
// and usage of a fiber scheduler.
// TODO: Current memory overhead of a Fiber is 65KB, this should be improved when implementing concurrency
pub struct Fiber {
    // All fibers have their own stack allocated.
    // This stack allocates a Vec with STACK_SIZE capacity.
    // An EmeraldObject occupies 8 bytes so this allocates 16KB.
    stack: Vec<Arc<EmeraldObject>>,
    // Always points to the next value. Top of stack is stack[sp-1]
    pub(crate) sp: u16,
    // All fibers also have their own call frames.
    // Fibers allocate a Vec with MAX_FRAMES capacity.
    // A frame takes up 48 bytes so this allocates roughly 49KB.
    pub(crate) frames: Vec<Frame>,
    fp: i32, // Points to current frame
}

// The Fiber implementation encapsulates all logic for operating on the stack, so the VM just calls
// these methods on the currently executing fiber
impl Fiber {
    pub(crate) fn new() -> Fiber {
        Fiber {
            stack: Vec::with_capacity(STACK_SIZE as usize),
            sp: 0,
            frames: Vec::with_capacity(MAX_FRAMES as usize),
            fp: -1,
        }
    }

    // Resets the frames and stack of the current fiber
    // Useful in testing to ensure isolation between test runs.
    pub fn reset(&mut self) {
        self.stack = Vec::with_capacity(STACK_SIZE as usize);
        self.sp = 0;
        self.frames = Vec::with_capacity(MAX_FRAMES as usize);
        self.fp = -1;
    }

    // fetches the object at the top of the stack
    pub(crate) fn stack_top(&self) -> Arc<EmeraldObject> {
        Arc::clone(self.stack.get((self.sp - 1) as usize).unwrap())
    }

    // fetches the last object that was popped off the stack
    pub(crate) fn last_popped_stack_object(&mut self) -> Arc<EmeraldObject> {
        Arc::clone(self.stack.get(self.sp as usize).unwrap())
    }

    // Push an object onto the stack
    pub(crate) fn push(&mut self, obj: Arc<EmeraldObject>) {
        self.stack.insert(self.sp as usize, obj);
        self.sp += 1;
    }

    // Pop an object from the stack
    pub(crate) fn pop(&mut self) -> Arc<EmeraldObject> {
        if self.sp == 0 {
            return EM_NIL.clone();
        }

        // TODO: This check probably has performance implication but it always leaves
        // TODO: the bottom of the stack which is nice for debugging and testing the VM
        let obj = if self.sp != 0 {
            self.stack_top()
        } else {
            self.stack.pop().unwrap()
        };

        self.sp -= 1;

        obj
    }

    pub(crate) fn current_frame(&mut self) -> &mut Frame {
        &mut self.frames[self.fp as usize]
    }

    pub(crate) fn is_base_frame(&self) -> bool {
        self.fp == 0
    }

    pub(crate) fn pop_frame(&mut self) {
        let frame = self.frames.pop().unwrap();
        self.fp -= 1;
        self.sp = frame.base_sp
    }

    pub(crate) fn push_frame(&mut self, frame: Frame) {
        self.frames.push(frame);
        self.fp += 1
    }

    // Local variable values are also stored on the stack
    // so methods for operating on them are also contained here

    // Creates a local variable binding on the stack
    pub(crate) fn insert_local(&mut self, index: usize, local: Arc<EmeraldObject>) {
        let local_index = self.locals_index(index);
        self.stack.insert(local_index, local);
        self.sp += 1; // Insert pushes all objects in the stack so we increment the stack pointer
    }

    // Resolves a local variable binding
    pub(crate) fn get_local(&mut self, index: usize) -> Arc<EmeraldObject> {
        let local_index = self.locals_index(index);
        self.stack.get(local_index).unwrap().clone()
    }

    fn locals_index(&mut self, index: usize) -> usize {
        self.current_frame().base_sp as usize + index
    }
}
