#[allow(unused_imports)]
use emerald::compiler::bytecode::Stringable;

#[cfg(test)]
pub mod parser {
    pub fn parse(input: &str) -> emerald::ast::AST {
        let mut parser = emerald::parser::Parser::new(emerald::lexer::input::Input::new(
            "test.rb".to_string(),
            input.to_string(),
        ));

        let ast = parser.parse();

        if parser.errors.len() != 0 {
            for error in parser.errors.clone() {
                println!("parser_test error: {}", error)
            }
        }
        assert_eq!(parser.errors.len(), 0, "failed to parse {}", input);

        ast
    }

    #[allow(dead_code)]
    pub fn test_expression_stmt<F>(stmt: emerald::ast::node::Statement, cb: F)
    where
        F: Fn(emerald::ast::node::Expression),
    {
        match stmt {
            emerald::ast::node::Statement::ExpressionStatement(expr) => cb(expr),
            _ => assert_eq!(
                0, 1,
                "statement is not expression statement \ngot={:?}",
                stmt,
            ),
        }
    }

    #[allow(dead_code)]
    pub fn test_boolean_object(expression: emerald::ast::node::Expression, expected: bool) {
        if expected {
            match expression {
                emerald::ast::node::Expression::TrueLiteral(_) => {}
                _ => assert!(false, "expression was not true, got={:?}", expression,),
            }
        } else {
            match expression {
                emerald::ast::node::Expression::FalseLiteral(_) => {}
                _ => assert!(false, "expression was not false, got={:?}", expression,),
            }
        }
    }

    #[allow(dead_code)]
    pub fn test_integer_object(expression: emerald::ast::node::Expression, expected: i64) {
        match expression {
            emerald::ast::node::Expression::IntegerLiteral(_data, val) => assert_eq!(val, expected),
            _ => assert_eq!(
                0, 1,
                "expression is not integer literal got={:?}",
                expression
            ),
        }
    }

    #[allow(dead_code)]
    pub fn test_float_object(expression: emerald::ast::node::Expression, expected: f64) {
        match expression {
            emerald::ast::node::Expression::FloatLiteral(_data, val) => assert_eq!(val, expected),
            _ => assert_eq!(0, 1, "expression is not float literal got={:?}", expression),
        }
    }

    #[allow(dead_code)]
    pub fn test_string_object(expression: emerald::ast::node::Expression, expected: &str) {
        match expression {
            emerald::ast::node::Expression::StringLiteral(data) => {
                assert_eq!(data.literal, expected)
            }
            _ => assert_eq!(
                0, 1,
                "expression is not string literal got={:?}",
                expression
            ),
        }
    }

    #[allow(dead_code)]
    pub fn test_identifier_object(expression: emerald::ast::node::Expression, expected: &str) {
        match expression {
            emerald::ast::node::Expression::IdentifierExpression(data) => {
                assert_eq!(data.literal, expected)
            }
            _ => assert_eq!(
                0, 1,
                "expression is not identifier expression got={:?}",
                expression
            ),
        }
    }
}

#[cfg(test)]
pub mod compiler {
    use std::sync::Arc;

    use emerald::compiler::bytecode::{Bytecode, Stringable};
    use emerald::compiler::Compiler;
    use emerald::kernel;
    use emerald::object::{Block, EmeraldObject, UnderlyingValueType};

    use super::parser;

    #[allow(dead_code)]
    pub struct CompilerTestCase<'a> {
        pub input: &'a str,
        pub expected_constants: Vec<UnderlyingValueType>,
        pub expected_bytecode: Bytecode,
    }

    #[allow(dead_code)]
    pub fn run_compiler_tests(cases: Vec<CompilerTestCase>) {
        for case in cases {
            kernel::reset_consts();

            let mut c = compile(&case.input);

            assert_eq!(
                c.bytecode_mut().to_string(),
                case.expected_bytecode.to_string(),
                "Bytecode did not match"
            );

            let pool = kernel::CONSTANT_POOL.lock().unwrap();

            assert_eq!(
                pool.len(),
                case.expected_constants.len(),
                "Unexpected amount of constants"
            );
            drop(pool);

            for (i, constant) in case.expected_constants.iter().enumerate() {
                let actual = kernel::get_const(i).unwrap();

                match constant {
                    UnderlyingValueType::Integer(expected) => {
                        test_integer_object(*expected, actual)
                    }
                    UnderlyingValueType::String(expected) => {
                        test_string_object(expected.to_string(), actual)
                    }
                    UnderlyingValueType::Symbol(expected) => {
                        test_symbol_object(expected.to_string(), actual)
                    }
                    UnderlyingValueType::Proc(block) => {
                        test_proc_object(block, actual);
                    }
                    _ => assert_eq!(0, 1, "Unknown expected object type"),
                }
            }
        }
    }

    #[allow(dead_code)]
    pub fn test_integer_object(expected: i64, actual: Arc<EmeraldObject>) {
        match actual.underlying_value {
            UnderlyingValueType::Integer(val) => assert_eq!(
                expected, val,
                "expected integer to equal {}, but got {}",
                expected, val
            ),
            _ => assert_eq!(
                expected,
                -12382190231,
                "Object is not Integer, got {}",
                actual.class_name()
            ),
        }
    }

    #[allow(dead_code)]
    pub fn test_string_object(expected: String, actual: Arc<EmeraldObject>) {
        match &actual.underlying_value {
            UnderlyingValueType::String(val) => assert_eq!(expected, *val),
            _ => assert_eq!(0, 1, "Object is not String"),
        }
    }

    #[allow(dead_code)]
    pub fn test_symbol_object(expected: String, actual: Arc<EmeraldObject>) {
        match &actual.underlying_value {
            UnderlyingValueType::Symbol(val) => assert_eq!(expected, *val),
            _ => assert_eq!(0, 1, "Object is not Symbol"),
        }
    }

    #[allow(dead_code)]
    pub fn test_boolean_object(expected: bool, actual: Arc<EmeraldObject>) {
        if expected {
            match &actual.underlying_value {
                UnderlyingValueType::True => {}
                _ => assert!(false, "expression was not true, got={:?}", actual),
            }
        } else {
            match &actual.underlying_value {
                UnderlyingValueType::False => {}
                _ => assert!(false, "expression was not false, got={:?}", actual),
            }
        }
    }

    #[allow(dead_code)]
    pub fn test_proc_object(expected: &Block, actual: Arc<EmeraldObject>) {
        match &actual.underlying_value {
            UnderlyingValueType::Proc(actual) => {
                assert_eq!(
                    actual.bytecode.to_string(),
                    expected.bytecode.to_string(),
                    "Proc bytecode did not match"
                );
                assert_eq!(actual.arity, expected.arity);
                assert_eq!(actual.num_locals, expected.num_locals);
            }
            _ => assert!(
                false,
                "object was not proc, got={:?}",
                &actual.underlying_value
            ),
        }
    }

    #[allow(dead_code)]
    pub fn compile(input: &str) -> Compiler {
        let mut c = Compiler::new();

        let ast = parser::parse(input);

        c.compile(ast);

        c
    }
}

#[cfg(test)]
pub mod vm {
    use emerald::kernel;
    use emerald::object::UnderlyingValueType;

    use super::compiler;

    #[allow(dead_code)]
    pub struct VMTestCase<'a> {
        pub input: &'a str,
        pub expected: UnderlyingValueType,
    }

    #[allow(dead_code)]
    pub fn run_vm_tests(cases: Vec<VMTestCase>) {
        for case in cases {
            let result = kernel::execute("test.rb".to_string(), case.input.to_string());

            match result {
                Ok(actual) => match case.expected {
                    UnderlyingValueType::Integer(expected) => {
                        compiler::test_integer_object(expected, actual)
                    }
                    UnderlyingValueType::String(expected) => {
                        compiler::test_string_object(expected, actual)
                    }
                    UnderlyingValueType::Symbol(expected) => {
                        compiler::test_symbol_object(expected, actual)
                    }
                    UnderlyingValueType::True => compiler::test_boolean_object(true, actual),
                    UnderlyingValueType::False => compiler::test_boolean_object(false, actual),
                    UnderlyingValueType::Nil => match actual.underlying_value {
                        UnderlyingValueType::Nil => {}
                        _ => assert!(false, "was not nil, got={:?}", actual.underlying_value),
                    },
                    _ => assert_eq!(0, 1, "Unknown expected object type"),
                },
                Err(err) => assert_eq!(0, 1, "VM test failed with error {}", err),
            }

            kernel::EMERALD_VM.lock().unwrap().current_fiber.reset();
        }
    }
}
