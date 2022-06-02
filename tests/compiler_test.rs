use emerald;
use emerald::compiler::bytecode::Opcode::{
    OpAdd, OpDiv, OpFalse, OpMul, OpNil, OpPop, OpPush, OpSend, OpSub, OpTrue,
};
use emerald::object::UnderlyingValueType;

mod helpers_test;
use helpers_test::compiler::*;

#[test]
fn test_compile_method_call() {
    let tests = Vec::from([CompilerTestCase {
        input: "2.to_s",
        expected_constants: Vec::from([
            UnderlyingValueType::Integer(2),
            UnderlyingValueType::Symbol("to_s".to_string()),
        ]),
        expected_bytecode: Vec::from([OpPush { index: 0 }, OpSend { index: 1 }, OpPop]),
    }]);

    run_compiler_tests(tests)
}

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
        CompilerTestCase {
            input: "1 * 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPush { index: 1 }, OpMul, OpPop]),
        },
        CompilerTestCase {
            input: "1 / 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPush { index: 1 }, OpDiv, OpPop]),
        },
        CompilerTestCase {
            input: "\"Hello \" + \"World!\"",
            expected_constants: Vec::from([
                UnderlyingValueType::String("Hello ".to_string()),
                UnderlyingValueType::String("World!".to_string()),
            ]),
            expected_bytecode: Vec::from([OpPush { index: 0 }, OpPush { index: 1 }, OpAdd, OpPop]),
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

#[test]
fn test_compile_true_literal() {
    let tests = Vec::from([CompilerTestCase {
        input: "true",
        expected_constants: Vec::from([]),
        expected_bytecode: Vec::from([OpTrue, OpPop]),
    }]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_false_literal() {
    let tests = Vec::from([CompilerTestCase {
        input: "false",
        expected_constants: Vec::from([]),
        expected_bytecode: Vec::from([OpFalse, OpPop]),
    }]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_nil_literal() {
    let tests = Vec::from([CompilerTestCase {
        input: "nil",
        expected_constants: Vec::from([]),
        expected_bytecode: Vec::from([OpNil, OpPop]),
    }]);

    run_compiler_tests(tests)
}
