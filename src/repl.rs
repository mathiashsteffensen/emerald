use std::ops::Add;
use std::sync::Arc;

use linefeed::{DefaultTerminal, Interface, ReadResult};

use crate::object::{EmeraldObject, ExecutionContext, UnderlyingValueType};
use crate::{compiler, debug, vm};

pub struct REPL {
    compiler: compiler::Compiler,
    reader: Interface<DefaultTerminal>,
    line: i64,
}

const PREFIX: &str = "iem(main):";
const HISTORY_FILE: &str = "/tmp/iem.hst";

impl REPL {
    pub fn new() -> REPL {
        REPL {
            compiler: compiler::Compiler::new(),
            reader: Interface::new("iem").unwrap(),
            line: 1,
        }
    }

    pub fn run(&mut self) {
        self.set_prompt();

        if let Err(e) = self.reader.load_history(HISTORY_FILE) {
            debug::log(format!(
                "Could not load history file {}: {}",
                HISTORY_FILE, e
            ));
        }

        while let ReadResult::Input(line) = self.reader.read_line().unwrap() {
            if line.as_str() == "quit" {
                println!("\nBye!\n");
                return;
            }

            if !line.trim().is_empty() {
                self.reader.add_history_unique(line.clone());
            }

            let result = self.interpret_line(line);

            if result.responds_to("inspect") {
                let stringified = result
                    .send(
                        result.clone(),
                        "inspect",
                        Arc::from(ExecutionContext::new(Arc::clone(&result))),
                        Vec::new(),
                    )
                    .unwrap();

                match &stringified.underlying_value {
                    UnderlyingValueType::String(str) => println!("=> {}", str),
                    _ => println!("{:?}#inspect did not return a string", result),
                }
            } else {
                println!("(Object does not support inspect)\n=>")
            }

            if let Err(e) = self.reader.save_history(HISTORY_FILE) {
                debug::log(format!(
                    "Could not save history file {}: {}",
                    HISTORY_FILE, e
                ));
            }

            self.line += 1;
            self.set_prompt();
        }
    }

    fn set_prompt(&self) {
        self.reader
            .set_prompt(&*PREFIX.to_string().add(&*format!("{:03}:0> ", self.line)))
            .unwrap();
    }

    fn interpret_line(&mut self, line: String) -> Arc<EmeraldObject> {
        self.compiler.compile_string("(iem)".to_string(), line);

        let mut vm = vm::VM::new(self.compiler.bytecode().clone(), Default::default());

        match vm.run() {
            Ok(()) => vm.last_popped_stack_object(),
            Err(e) => panic!("{}", e),
        }
    }
}
