use std::io;
use std::io::{BufRead, Stdin};
use std::rc::Rc;

use linefeed::{Interface, ReadResult};

use crate::object::{EmeraldObject, ExecutionContext, UnderlyingValueType};
use crate::{core, vm};

pub struct REPL {
    stdin: Stdin,
}

const PREFIX: &str = "(iem)>>> ";

impl REPL {
    pub fn new() -> REPL {
        REPL { stdin: io::stdin() }
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

            if result.responds_to("to_s") {
                let stringified = result
                    .send(
                        "to_s",
                        Rc::from(ExecutionContext::new(Rc::clone(&result), core::all::map())),
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

    fn next(&self) -> String {
        let mut line = String::new();

        self.stdin.lock().read_line(&mut line).unwrap();

        line
    }

    fn interpret_line(&self, line: String) -> Rc<EmeraldObject> {
        let result = vm::VM::interpret("(iem)".to_string(), line);

        match result {
            Ok((_, mut vm)) => vm.last_popped_stack_object(),
            Err(e) => panic!("{}", e),
        }
    }
}
