extern crate core;

mod helpers_test;

use emerald::object::UnderlyingValueType;
use helpers_test::vm::*;

#[test]
fn test_literals() {
    let tests = Vec::from([
        VMTestCase {
            input: "12",
            expected: UnderlyingValueType::Integer(12),
        },
        VMTestCase {
            input: "1_2",
            expected: UnderlyingValueType::Integer(12),
        },
        VMTestCase {
            input: "true",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "false",
            expected: UnderlyingValueType::False,
        },
        VMTestCase {
            input: "nil",
            expected: UnderlyingValueType::Nil,
        },
    ]);

    run_vm_tests(tests);
}

#[test]
fn test_infix_operations() {
    let tests = Vec::from([
        VMTestCase {
            input: "12 - 8",
            expected: UnderlyingValueType::Integer(4),
        },
        VMTestCase {
            input: "54 + 86",
            expected: UnderlyingValueType::Integer(140),
        },
        VMTestCase {
            input: "54 * 86",
            expected: UnderlyingValueType::Integer(4644),
        },
        VMTestCase {
            input: "15 / 3",
            expected: UnderlyingValueType::Integer(5),
        },
        VMTestCase {
            input: "3 > 3",
            expected: UnderlyingValueType::False,
        },
        VMTestCase {
            input: "4 > 3",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "3 >= 3",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "3 < 3",
            expected: UnderlyingValueType::False,
        },
        VMTestCase {
            input: "2 < 3",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "3 <= 3",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "3 <= 3",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "3 == 3",
            expected: UnderlyingValueType::True,
        },
        VMTestCase {
            input: "3 == 4",
            expected: UnderlyingValueType::False,
        },
    ]);

    run_vm_tests(tests);
}

#[test]
fn test_method_literals() {
    let tests = Vec::from([VMTestCase {
        input: "def get_num
            2
        end",
        expected: UnderlyingValueType::Symbol("get_num".to_string()),
    }]);

    run_vm_tests(tests);
}

#[test]
fn test_method_calls() {
    let tests = Vec::from([
        // VMTestCase {
        //     input: "2.to_s",
        //     expected: UnderlyingValueType::String("2".to_string()),
        // },
        // VMTestCase {
        //     input: "def num
        //         2
        //     end
        //
        //    num",
        //     expected: UnderlyingValueType::Integer(2),
        // },
        // VMTestCase {
        //     input: "def add(x, y)
        //         x + y
        //     end
        //
        //    add(5, 4)",
        //     expected: UnderlyingValueType::Integer(9),
        // },
        VMTestCase {
            input: "def add(x, y)
                if x == 5
                    return 0
                end

                x == 5
            end

           add(5, 4)",
            expected: UnderlyingValueType::Integer(0),
        },
    ]);

    run_vm_tests(tests);
}

#[test]
fn test_global_assignments() {
    let tests = Vec::from([VMTestCase {
        input: "var = 5; var + 5",
        expected: UnderlyingValueType::Integer(10),
    }]);

    run_vm_tests(tests);
}

#[test]
fn test_if_expressions() {
    let tests = Vec::from([
        VMTestCase {
            input: "if true
                5
            end",
            expected: UnderlyingValueType::Integer(5),
        },
        VMTestCase {
            input: "if true
                5
            else
                25
            end",
            expected: UnderlyingValueType::Integer(5),
        },
        VMTestCase {
            input: "if false
                5
            end",
            expected: UnderlyingValueType::Nil,
        },
        VMTestCase {
            input: "if false
                5
            else
                25
            end",
            expected: UnderlyingValueType::Integer(25),
        },
    ]);

    run_vm_tests(tests);
}

#[test]
fn test_return_statements() {
    let tests = Vec::from([
        VMTestCase {
            input: "if true
                return 5
            end
            10",
            expected: UnderlyingValueType::Integer(5),
        },
        VMTestCase {
            input: "if false
                return 5
            end
            10",
            expected: UnderlyingValueType::Integer(10),
        },
        VMTestCase {
            input: "if false
                return
            end
            25",
            expected: UnderlyingValueType::Integer(25),
        },
        VMTestCase {
            input: "if true
                return
            end
            25",
            expected: UnderlyingValueType::Nil,
        },
    ]);

    run_vm_tests(tests);
}
