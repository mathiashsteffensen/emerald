mod helpers_test;

use emerald::object::UnderlyingValueType;
use helpers_test::vm::*;

#[test]
fn test_integer_literal() {
    let tests = Vec::from([
        VMTestCase {
            input: "12",
            expected: UnderlyingValueType::Integer(12),
        },
        VMTestCase {
            input: "1_2",
            expected: UnderlyingValueType::Integer(12),
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
    ]);

    run_vm_tests(tests);
}
