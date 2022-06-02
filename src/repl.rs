use std::sync::Arc;

use linefeed::{Interface, ReadResult};

use crate::object::{EmeraldObject, ExecutionContext, UnderlyingValueType};
use crate::vm;

pub struct REPL {}

const PREFIX: &str = "(iem)>>> ";

impl REPL {
    pub fn new() -> REPL {
        REPL {}
    }

    pub fn run(&mut self) {
        let reader = Interface::new("iem").unwrap();

        if let Err(e) = reader.set_prompt(PREFIX) {
            panic!("Error initializing REPL: {}", e);
        };

        while let ReadResult::Input(line) = reader.read_line().unwrap() {
            if line.as_str() == "quit" {
                println!("\nBye!\n");
                return;
            }

            let result = self.interpret_line(line);

            if result.responds_to("inspect") {
                let stringified = result
                    .send(
                        "inspect",
                        Arc::from(ExecutionContext::new(Arc::clone(&result))),
                        Vec::new(),
                    )
                    .unwrap();

                match &stringified.underlying_value {
                    UnderlyingValueType::String(str) => println!("{}", str),
                    _ => println!("{:?}#to_s did not return a string", result),
                }
            }
        }
    }

    fn interpret_line(&self, line: String) -> Arc<EmeraldObject> {
        let result = vm::VM::interpret("(iem)".to_string(), line);

        match result {
            Ok((_, mut vm)) => vm.last_popped_stack_object(),
            Err(e) => panic!("{}", e),
        }
    }
}
