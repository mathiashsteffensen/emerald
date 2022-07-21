use crate::compiler::bytecode::Bytecode;
use crate::compiler::symbol_table::SymbolTable;
use crate::compiler::Compiler;

pub(crate) struct CompilationScope {
    pub(crate) bytecode: Bytecode,
}

pub(crate) fn new() -> CompilationScope {
    CompilationScope {
        bytecode: Vec::new(),
    }
}

pub(crate) fn enter(c: &mut Compiler) {
    let scope = new();

    c.scopes.push(scope);
    c.scope_index += 1;

    c.symbol_table = SymbolTable::with_outer(c.symbol_table.clone());
}

pub(crate) fn leave(c: &mut Compiler) -> (Bytecode, u16) {
    let bytecode = c.scopes.pop().unwrap().bytecode;
    let num_definitions = c.symbol_table.num_definitions;

    c.scope_index -= 1;

    c.symbol_table = *c.symbol_table.outer.as_ref().unwrap().clone();

    (bytecode, num_definitions)
}
