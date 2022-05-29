use std::rc::Rc;

use emerald;
use emerald::compiler::bytecode::Bytecode;
use emerald::compiler::bytecode::Opcode::{OpAdd, OpPop, OpPush, OpSub};
use emerald::compiler::Compiler;
use emerald::object::{EmeraldObject, UnderlyingValueType};

mod helpers_test;
use helpers_test::compiler::*;

#[test]
fn test_compile_infix_expression() {
    let tests = Vec::from([
        CompilerTestCase {
            input: "1 + 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPush { index: 1 }, OpAdd, OpPop]),
        },
        CompilerTestCase {
            input: "1 - 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPush { index: 1 }, OpSub, OpPop]),
        },
    ]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_integer_literal() {
    let tests = Vec::from([
        CompilerTestCase {
            input: "1",
            expected_constants: Vec::from([UnderlyingValueType::Integer(1)]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPop]),
        },
        CompilerTestCase {
            input: "123_968",
            expected_constants: Vec::from([UnderlyingValueType::Integer(123_968)]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPop]),
        },
    ]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_string_literal() {
    let tests = Vec::from([CompilerTestCase {
        input: "\"Hello World\"",
        expected_constants: Vec::from([UnderlyingValueType::String("Hello World".to_string())]),
        expected_bytecode: Vec::from([OpPush { index: 0 }, OpPop]),
    }]);

    run_compiler_tests(tests)
}
