use emerald::compiler::bytecode::Stringable;

pub mod parser {
    pub fn parse(input: &str) -> emerald::ast::AST {
        let mut parser = emerald::parser::Parser::new(emerald::lexer::input::Input::new(
            "test.rb".to_string(),
            input.to_string(),
        ));

        let ast = parser.parse_ast();

        if parser.errors.len() != 0 {
            for error in parser.errors.clone() {
                println!("parser_test error: {}", error)
            }
        }
        assert_eq!(parser.errors.len(), 0, "failed to parse {}", input);

        ast
    }

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

    pub fn test_float_object(expression: emerald::ast::node::Expression, expected: f64) {
        match expression {
            emerald::ast::node::Expression::FloatLiteral(_data, val) => assert_eq!(val, expected),
            _ => assert_eq!(0, 1, "expression is not float literal got={:?}", expression),
        }
    }

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

pub mod compiler {
    use std::rc::Rc;

    use emerald::compiler::bytecode::{Bytecode, Stringable};
    use emerald::compiler::Compiler;
    use emerald::object::{EmeraldObject, UnderlyingValueType};

    use super::parser;

    pub struct CompilerTestCase<'a> {
        pub input: &'a str,
        pub expected_constants: Vec<UnderlyingValueType>,
        pub expected_bytecode: Bytecode,
    }

    pub fn run_compiler_tests(cases: Vec<CompilerTestCase>) {
        for case in cases {
            let c = compile(&case.input);

            assert_eq!(
                c.bytecode.to_string(),
                case.expected_bytecode.to_string(),
                "Bytecode did not match"
            );

            assert_eq!(
                c.constant_pool.len(),
                case.expected_constants.len(),
                "Unexpected amount of constants"
            );

            for (i, constant) in case.expected_constants.iter().enumerate() {
                let actual = Rc::clone(c.constant_pool.get(i).unwrap());

                match constant {
                    UnderlyingValueType::Integer(expected) => {
                        test_integer_object(*expected, actual)
                    }
                    UnderlyingValueType::String(expected) => {
                        test_string_object(expected.to_string(), actual)
                    }
                    _ => assert_eq!(0, 1, "Unknown expected object type"),
                }
            }
        }
    }

    pub fn test_integer_object(expected: i64, actual: Rc<EmeraldObject>) {
        match actual.underlying_value {
            UnderlyingValueType::Integer(val) => assert_eq!(expected, val),
            _ => assert_eq!(0, 1, "Object is not Integer"),
        }
    }

    pub fn test_string_object(expected: String, actual: Rc<EmeraldObject>) {
        match &actual.underlying_value {
            UnderlyingValueType::String(val) => assert_eq!(expected, *val),
            _ => assert_eq!(0, 1, "Object is not String"),
        }
    }

    pub fn compile(input: &str) -> Compiler {
        let mut c = Compiler::new();

        let ast = parser::parse(input);

        c.compile(ast);

        c
    }
}
