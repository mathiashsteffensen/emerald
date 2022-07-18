use emerald;
use emerald::compiler::bytecode::Opcode::{
    OpAdd, OpDefineMethod, OpDiv, OpFalse, OpGetGlobal, OpGetLocal, OpGreaterThan,
    OpGreaterThanOrEq, OpJump, OpJumpNotTruthy, OpLessThan, OpLessThanOrEq, OpMul, OpNil, OpPop,
    OpPush, OpPushSelf, OpReturn, OpReturnValue, OpSend, OpSetGlobal, OpSub, OpTrue,
};
use emerald::object::{Block, UnderlyingValueType};

mod helpers_test;
use helpers_test::compiler::*;

#[test]
fn test_compile_return_statement() {
    let tests = Vec::from([
        CompilerTestCase {
            input: "if true
                return 5
            end",
            expected_constants: Vec::from([UnderlyingValueType::Integer(5)]),
            expected_bytecode: Vec::from([
                OpTrue,
                OpJumpNotTruthy { offset: 3 },
                OpPush { index: 0 },
                OpReturnValue,
                OpJump { offset: 1 },
                OpNil,
                OpPop,
            ]),
        },
        CompilerTestCase {
            input: "if true
                return
            end",
            expected_constants: Vec::from([]),
            expected_bytecode: Vec::from([
                OpTrue,
                OpJumpNotTruthy { offset: 2 },
                OpReturn,
                OpJump { offset: 1 },
                OpNil,
                OpPop,
            ]),
        },
    ]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_method_literal() {
    let tests = Vec::from([
        CompilerTestCase {
            input: "def num
                2
            end",
            expected_constants: Vec::from([
                UnderlyingValueType::Symbol("num".to_string()),
                UnderlyingValueType::Integer(2),
                UnderlyingValueType::Proc(Block::new(
                    Vec::from([OpPush { index: 1 }, OpReturnValue]),
                    0,
                    0,
                )),
            ]),
            expected_bytecode: Vec::from([
                OpDefineMethod {
                    name_index: 0,
                    proc_index: 2,
                },
                OpPop,
            ]),
        },
        CompilerTestCase {
            input: "def my_method(arg, other)
                puts(arg, other)
            end

            my_method(3, 4)",
            expected_constants: Vec::from([
                UnderlyingValueType::Symbol("my_method".to_string()),
                UnderlyingValueType::Symbol("puts".to_string()),
                UnderlyingValueType::Proc(Block::new(
                    Vec::from([
                        OpGetLocal { index: 0 },
                        OpGetLocal { index: 1 },
                        OpPushSelf,
                        OpSend {
                            index: 1,
                            num_args: 2,
                        },
                        OpReturnValue,
                    ]),
                    2,
                    0,
                )),
                UnderlyingValueType::Integer(3),
                UnderlyingValueType::Integer(4),
                UnderlyingValueType::Symbol("my_method".to_string()),
            ]),
            expected_bytecode: Vec::from([
                OpDefineMethod {
                    name_index: 0,
                    proc_index: 2,
                },
                OpPop,
                OpPush { index: 3 },
                OpPush { index: 4 },
                OpPushSelf,
                OpSend {
                    index: 5,
                    num_args: 2,
                },
                OpPop,
            ]),
        },
    ]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_method_call() {
    let tests = Vec::from([CompilerTestCase {
        input: "2.to_s",
        expected_constants: Vec::from([
            UnderlyingValueType::Integer(2),
            UnderlyingValueType::Symbol("to_s".to_string()),
        ]),
        expected_bytecode: Vec::from([
            OpPush { index: 0 },
            OpSend {
                index: 1,
                num_args: 0,
            },
            OpPop,
        ]),
    }]);

    run_compiler_tests(tests)
}

#[test]
fn test_compile_global_assignments() {
    let tests = Vec::from([CompilerTestCase {
        input: "var = 5; var + 5",
        expected_constants: Vec::from([
            UnderlyingValueType::Integer(5),
            UnderlyingValueType::Integer(5),
        ]),
        expected_bytecode: Vec::from([
            OpPush { index: 0 },
            OpSetGlobal { index: 0 },
            OpPop,
            OpGetGlobal { index: 0 },
            OpPush { index: 1 },
            OpAdd,
            OpPop,
        ]),
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
            input: "1 > 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([
                OpPush { index: 0 },
                OpPush { index: 1 },
                OpGreaterThan,
                OpPop,
            ]),
        },
        CompilerTestCase {
            input: "1 >= 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([
                OpPush { index: 0 },
                OpPush { index: 1 },
                OpGreaterThanOrEq,
                OpPop,
            ]),
        },
        CompilerTestCase {
            input: "1 < 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([
                OpPush { index: 0 },
                OpPush { index: 1 },
                OpLessThan,
                OpPop,
            ]),
        },
        CompilerTestCase {
            input: "1 <= 2",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::Integer(2),
            ]),
            expected_bytecode: Vec::from([
                OpPush { index: 0 },
                OpPush { index: 1 },
                OpLessThanOrEq,
                OpPop,
            ]),
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

#[test]
fn test_compile_if_expression() {
    let tests = Vec::from([
        CompilerTestCase {
            input: "if true
                5
            end",
            expected_constants: Vec::from([UnderlyingValueType::Integer(5)]),
            expected_bytecode: Vec::from([
                OpTrue,
                OpJumpNotTruthy { offset: 2 },
                OpPush { index: 0 },
                OpJump { offset: 1 },
                OpNil,
                OpPop,
            ]),
        },
        CompilerTestCase {
            input: "if true
                5
                4
                3
                2
                1
            else
                \"Hello\"
            end",
            expected_constants: Vec::from([
                UnderlyingValueType::Integer(5),
                UnderlyingValueType::Integer(4),
                UnderlyingValueType::Integer(3),
                UnderlyingValueType::Integer(2),
                UnderlyingValueType::Integer(1),
                UnderlyingValueType::String("Hello".to_string()),
            ]),
            expected_bytecode: Vec::from([
                OpTrue,
                OpJumpNotTruthy { offset: 10 },
                OpPush { index: 0 },
                OpPop,
                OpPush { index: 1 },
                OpPop,
                OpPush { index: 2 },
                OpPop,
                OpPush { index: 3 },
                OpPop,
                OpPush { index: 4 },
                OpJump { offset: 1 },
                OpPush { index: 5 },
                OpPop,
            ]),
        },
    ]);

    run_compiler_tests(tests)
}
