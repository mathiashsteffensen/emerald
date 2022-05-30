use std::{env, fs};

use emerald;

#[cfg(not(tarpaulin_include))]
fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() == 1 {
        let mut repl = emerald::repl::REPL::new();

        repl.run();
    } else {
        let file_name = args.get(1).unwrap();
        let content = fs::read_to_string(file_name).expect("Failed to read file");

        emerald::vm::VM::interpret(file_name.to_string(), content).expect("Interpreter failed");
    }
}
