use emerald;
use emerald::ast::node::Node;

mod helpers_test;
use helpers_test::parser::*;

#[test]
fn test_parse_return_statement() {
    let ast = parse("return 5;");

    assert_eq!(ast.statements.len(), 1);
    assert_eq!(ast.statements[0].token_literal(), "return");
    assert_eq!(ast.statements[0].to_string(), "return 5;");

    match ast.statements[0].clone() {
        emerald::ast::node::Statement::ReturnStatement(_data, expr) => match expr {
            Some(expr) => test_integer_object(expr, 5),
            None => assert_eq!(0, 1, "return statement value is None"),
        },
        _ => assert_eq!(
            0, 1,
            "statement is not return statement \ngot={:?}",
            ast.statements[0]
        ),
    }

    let ast = parse("return");

    assert_eq!(ast.statements.len(), 1);
    assert_eq!(ast.statements[0].token_literal(), "return");
    assert_eq!(ast.statements[0].to_string(), "return;");

    match ast.statements[0].clone() {
        emerald::ast::node::Statement::ReturnStatement(_data, expr) => match expr {
            None => {}
            Some(expr) => assert_eq!(
                0, 1,
                "return statement value is not None \ngot=Some({:?})",
                expr,
            ),
        },
        _ => assert_eq!(
            0, 1,
            "statement is not return statement \ngot={:?}",
            ast.statements[0]
        ),
    }
}

#[test]
fn test_parse_integer_expression_statement() {
    let ast = parse("1_5;");

    assert_eq!(ast.statements.len(), 1);
    assert_eq!(ast.statements[0].token_literal(), "1_5");
    assert_eq!(ast.statements[0].to_string(), "1_5");

    test_expression_stmt(ast.statements[0].clone(), {
        |expr| test_integer_object(expr, 15)
    });
}

#[test]
fn test_parse_float_expression_statement() {
    let ast = parse("15.78;");

    assert_eq!(ast.statements.len(), 1);
    assert_eq!(ast.statements[0].token_literal(), "15.78");
    assert_eq!(ast.statements[0].to_string(), "15.78");

    test_expression_stmt(ast.statements[0].clone(), {
        |expr| test_float_object(expr, 15.78)
    });
}

#[test]
fn test_parse_string_expression_statement() {
    let ast = parse("\"This is a string\";");

    assert_eq!(ast.statements.len(), 1);
    assert_eq!(ast.statements[0].token_literal(), "This is a string");
    assert_eq!(ast.statements[0].to_string(), "This is a string");

    test_expression_stmt(ast.statements[0].clone(), {
        |expr| test_string_object(expr, "This is a string")
    });
}

#[test]
fn test_parse_assignment_expression() {
    let ast = parse("var = 25");

    assert_eq!(ast.statements.len(), 1);
    assert_eq!(ast.statements[0].token_literal(), "=");
    assert_eq!(ast.statements[0].to_string(), "var = 25;");

    test_expression_stmt(ast.statements[0].clone(), {
        |expr| match expr {
            emerald::ast::node::Expression::AssignmentExpression(name, _data, value) => {
                test_identifier_object(*name, "var");
                test_integer_object(*value, 25);
            }
            _ => assert_eq!(
                0, 1,
                "expression is not assignment expression \ngot={:?}",
                ast.statements[0]
            ),
        }
    });
}

#[test]
fn test_parse_method_call() {
    let input = "puts(1, 6.56, \"Hello World!\")";

    let ast = parse(input);

    assert_eq!(ast.statements.len(), 1);

    test_expression_stmt(ast.statements[0].clone(), {
        |expr| match expr {
            emerald::ast::node::Expression::MethodCall(data) => {
                test_identifier_object(*data.ident, "puts");

                assert_eq!(data.args.len(), 3);

                test_integer_object(data.args[0].clone(), 1);
                test_float_object(data.args[1].clone(), 6.56);
                test_string_object(data.args[2].clone(), "Hello World!");
            }
            _ => assert_eq!(
                0, 1,
                "expression is not method call \ngot={:?}",
                ast.statements[0]
            ),
        }
    });
}

#[test]
fn test_method_literal() {
    struct Test {
        input: String,
        name: String,
        args: Vec<String>,
        num_stmts: i16,
    }

    let tests = Vec::from([
        Test {
            input: "def method(arg)
                        arg
                    end"
            .to_string(),
            name: "method".to_string(),
            args: Vec::from(["arg".to_string()]),
            num_stmts: 1,
        },
        Test {
            input: "def method
                       var

                        other_var
                    end"
            .to_string(),
            name: "method".to_string(),
            args: Vec::new(),
            num_stmts: 2,
        },
        Test {
            input: "def method(x, y)
                       x + y
                    end"
            .to_string(),
            name: "method".to_string(),
            args: Vec::from(["x".to_string(), "y".to_string()]),
            num_stmts: 1,
        },
        Test {
            input: "def method(x, y); x + y; end".to_string(),
            name: "method".to_string(),
            args: Vec::from(["x".to_string(), "y".to_string()]),
            num_stmts: 1,
        },
        Test {
            input: "def method(x, y)
                        x + y

                        def other_method
                            \"BOO\"
                        end

                        other_method
                    end"
            .to_string(),
            name: "method".to_string(),
            args: Vec::from(["x".to_string(), "y".to_string()]),
            num_stmts: 3,
        },
    ]);

    for test in tests {
        let ast = parse(&test.input);

        assert_eq!(ast.statements.len(), 1);

        test_expression_stmt(ast.statements[0].clone(), {
            |expr| match expr {
                emerald::ast::node::Expression::MethodLiteral(_data, name, args, body) => {
                    test_identifier_object(*name, &*test.name);

                    assert_eq!(args.len(), test.args.len());

                    for (i, arg) in args.iter().enumerate() {
                        test_identifier_object(arg.clone(), &*test.args[i])
                    }

                    assert_eq!(
                        body.len(),
                        test.num_stmts as usize,
                        "wrong num statements, got={:?}",
                        body
                    )
                }
                _ => assert_eq!(
                    0, 1,
                    "expression is not method literal \ngot={:?}",
                    ast.statements[0]
                ),
            }
        });
    }
}

#[test]
fn test_class_literal() {
    let input = "class MyClass
        def my_method
            do_stuff + 5
        end
    end";

    let ast = parse(input);

    assert_eq!(ast.statements.len(), 1);

    test_expression_stmt(ast.statements[0].clone(), {
        |expr| match expr {
            emerald::ast::node::Expression::ClassLiteral(_data, name, body) => {
                test_identifier_object(*name, "MyClass");
                assert_eq!(body.len(), 1);
            }
            _ => assert_eq!(
                0, 1,
                "expression is not class literal \ngot={:?}",
                ast.statements[0]
            ),
        }
    })
}

#[test]
fn test_operator_precedence_parsing() {
    let tests = Vec::from([
        ["-a * b", "((-a) * b)"],
        ["!-a", "(!(-a))"],
        ["a + b + c", "((a + b) + c)"],
        ["a + b - c", "((a + b) - c)"],
        ["a * b * c", "((a * b) * c)"],
        ["a * b / c", "((a * b) / c)"],
        ["a + b / c", "(a + (b / c))"],
        ["a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"],
        ["3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"],
        ["5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"],
        ["5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"],
        [
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        ],
        ["(5 + 5) * 2 * (5 + 5)", "(((5 + 5) * 2) * (5 + 5))"],
    ]);

    for test in tests {
        let ast = parse(test[0]);

        assert_eq!(ast.to_string(), test[1])
    }
}

#[test]
fn test_syntax_errors() {
    let tests = Vec::from([
        ["(5+4;", "syntax error at line:1:5: expected ')', found ';'"],
        [
            "var = (5+4;",
            "syntax error at line:1:11: expected ')', found ';'",
        ],
        [
            "def hello; 2;",
            "syntax error at line:1:14: expected 'end', found 'EOF'",
        ],
        [
            "def hello
               do_stuff

               do_other_stuff
            ",
            "syntax error at line:5:13: expected 'end', found 'EOF'",
        ],
    ]);

    for test in tests {
        let mut parser = emerald::parser::Parser::new(emerald::lexer::input::Input::new(
            "test.rb".to_string(),
            test[0].to_string(),
        ));

        parser.parse_ast();

        assert_eq!(parser.errors.len(), 1);
        assert_eq!(parser.errors[0], test[1]);
    }
}
