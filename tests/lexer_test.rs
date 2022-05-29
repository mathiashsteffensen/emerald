use emerald::lexer;
use emerald::lexer::input;
use emerald::lexer::token::Token;
use emerald::lexer::token::TokenData;

fn new_token_data(literal: &str, line: i64, column: i64, pos: i64) -> TokenData {
    return TokenData {
        literal: literal.to_string(),
        line,
        column,
        pos,
    };
}

struct LexerTest {
    input: String,
    expected: Vec<Token>,
}

impl LexerTest {
    fn run(&self) {
        let mut lexer = lexer::Lexer::new(input::Input::new(
            "test.rb".to_string(),
            self.input.to_string(),
        ));

        for tok in &self.expected {
            assert_eq!(*tok, lexer.next_token())
        }
    }
}

#[test]
fn test_lex_identifiers() {
    let test = LexerTest {
        input: "
            five = bar;
            ten = bar;

            class Integer
                def add_num(y)
                    self y
                end
            end

            if return true else return false end
            nil module begin rescue ensure $: $LOAD_PATH
        "
        .to_string(),
        expected: Vec::from([
            Token::Newline(new_token_data("\n", 1, 1, 0)),
            Token::Ident(new_token_data("five", 2, 13, 13)),
            Token::Assign(new_token_data("=", 2, 18, 18)),
            Token::Ident(new_token_data("bar", 2, 20, 20)),
            Token::Semicolon(new_token_data(";", 2, 23, 23)),
            Token::Newline(new_token_data("\n", 2, 24, 24)),
            Token::Ident(new_token_data("ten", 3, 13, 37)),
            Token::Assign(new_token_data("=", 3, 17, 41)),
            Token::Ident(new_token_data("bar", 3, 19, 43)),
            Token::Semicolon(new_token_data(";", 3, 22, 46)),
            Token::Newline(new_token_data("\n", 3, 23, 47)),
            Token::Newline(new_token_data("\n", 4, 1, 48)),
            Token::Class(new_token_data("class", 5, 13, 61)),
            Token::Ident(new_token_data("Integer", 5, 19, 67)),
            Token::Newline(new_token_data("\n", 5, 26, 74)),
            Token::Def(new_token_data("def", 6, 17, 91)),
            Token::Ident(new_token_data("add_num", 6, 21, 95)),
            Token::LeftParen(new_token_data("(", 6, 28, 102)),
            Token::Ident(new_token_data("y", 6, 29, 103)),
            Token::RightParen(new_token_data(")", 6, 30, 104)),
            Token::Newline(new_token_data("\n", 6, 31, 105)),
            Token::Selff(new_token_data("self", 7, 21, 126)),
            Token::Ident(new_token_data("y", 7, 26, 131)),
            Token::Newline(new_token_data("\n", 7, 27, 132)),
            Token::End(new_token_data("end", 8, 17, 149)),
            Token::Newline(new_token_data("\n", 8, 20, 152)),
            Token::End(new_token_data("end", 9, 13, 165)),
            Token::Newline(new_token_data("\n", 9, 16, 168)),
            Token::Newline(new_token_data("\n", 10, 1, 169)),
            Token::If(new_token_data("if", 11, 13, 182)),
            Token::Return(new_token_data("return", 11, 16, 185)),
            Token::True(new_token_data("true", 11, 23, 192)),
            Token::Else(new_token_data("else", 11, 28, 197)),
            Token::Return(new_token_data("return", 11, 33, 202)),
            Token::False(new_token_data("false", 11, 40, 209)),
            Token::End(new_token_data("end", 11, 46, 215)),
            Token::Newline(new_token_data("\n", 11, 49, 218)),
            Token::Nil(new_token_data("nil", 12, 13, 231)),
            Token::Module(new_token_data("module", 12, 17, 235)),
            Token::Begin(new_token_data("begin", 12, 24, 242)),
            Token::Rescue(new_token_data("rescue", 12, 30, 248)),
            Token::Ensure(new_token_data("ensure", 12, 37, 255)),
            Token::GlobalIdent(new_token_data("$:", 12, 44, 262)),
            Token::GlobalIdent(new_token_data("$LOAD_PATH", 12, 47, 265)),
        ]),
    };

    test.run()
}

#[test]
fn test_lex_method_calls() {
    let test = LexerTest {
        input: "
            @result = add(five, ten);

            @result.method

            print(result)

            [].each { |w| do_stuff(w) }
            [].each do |w|
                do_stuff(w)
            end
        "
        .to_string(),
        expected: Vec::from([
            Token::Newline(new_token_data("\n", 1, 1, 0)),
            Token::InstanceVar(new_token_data("@result", 2, 13, 13)),
            Token::Assign(new_token_data("=", 2, 21, 21)),
            Token::Ident(new_token_data("add", 2, 23, 23)),
            Token::LeftParen(new_token_data("(", 2, 26, 26)),
            Token::Ident(new_token_data("five", 2, 27, 27)),
            Token::Comma(new_token_data(",", 2, 31, 31)),
            Token::Ident(new_token_data("ten", 2, 33, 33)),
            Token::RightParen(new_token_data(")", 2, 36, 36)),
            Token::Semicolon(new_token_data(";", 2, 37, 37)),
            Token::Newline(new_token_data("\n", 2, 38, 38)),
            Token::Newline(new_token_data("\n", 3, 1, 39)),
            Token::InstanceVar(new_token_data("@result", 4, 13, 52)),
            Token::Dot(new_token_data(".", 4, 20, 59)),
            Token::Ident(new_token_data("method", 4, 21, 60)),
            Token::Newline(new_token_data("\n", 4, 27, 66)),
            Token::Newline(new_token_data("\n", 5, 1, 67)),
            Token::Ident(new_token_data("print", 6, 13, 80)),
            Token::LeftParen(new_token_data("(", 6, 18, 85)),
            Token::Ident(new_token_data("result", 6, 19, 86)),
            Token::RightParen(new_token_data(")", 6, 25, 92)),
            Token::Newline(new_token_data("\n", 6, 26, 93)),
            Token::Newline(new_token_data("\n", 7, 1, 94)),
            Token::LeftBracket(new_token_data("[", 8, 13, 107)),
            Token::RightBracket(new_token_data("]", 8, 14, 108)),
            Token::Dot(new_token_data(".", 8, 15, 109)),
            Token::Ident(new_token_data("each", 8, 16, 110)),
            Token::LeftBrace(new_token_data("{", 8, 21, 115)),
            Token::BitOr(new_token_data("|", 8, 23, 117)),
            Token::Ident(new_token_data("w", 8, 24, 118)),
            Token::BitOr(new_token_data("|", 8, 25, 119)),
            Token::Ident(new_token_data("do_stuff", 8, 27, 121)),
            Token::LeftParen(new_token_data("(", 8, 35, 129)),
            Token::Ident(new_token_data("w", 8, 36, 130)),
            Token::RightParen(new_token_data(")", 8, 37, 131)),
            Token::RightBrace(new_token_data("}", 8, 39, 133)),
            Token::Newline(new_token_data("\n", 8, 40, 134)),
            Token::LeftBracket(new_token_data("[", 9, 13, 147)),
            Token::RightBracket(new_token_data("]", 9, 14, 148)),
            Token::Dot(new_token_data(".", 9, 15, 149)),
            Token::Ident(new_token_data("each", 9, 16, 150)),
            Token::Do(new_token_data("do", 9, 21, 155)),
            Token::BitOr(new_token_data("|", 9, 24, 158)),
            Token::Ident(new_token_data("w", 9, 25, 159)),
            Token::BitOr(new_token_data("|", 9, 26, 160)),
            Token::Newline(new_token_data("\n", 9, 27, 161)),
            Token::Ident(new_token_data("do_stuff", 10, 17, 178)),
            Token::LeftParen(new_token_data("(", 10, 25, 186)),
            Token::Ident(new_token_data("w", 10, 26, 187)),
            Token::RightParen(new_token_data(")", 10, 27, 188)),
            Token::Newline(new_token_data("\n", 10, 28, 189)),
            Token::End(new_token_data("end", 11, 13, 202)),
        ]),
    };

    test.run()
}

#[test]
fn test_lex_operators() {
    let test = LexerTest {
        input: "# This is a comment
            !-/*; # This is a comment
            <> <<;
            == != >= <= & && | || &&= ||=
        "
        .to_string(),
        expected: Vec::from([
            Token::Newline(new_token_data("\n", 1, 20, 19)),
            Token::Bang(new_token_data("!", 2, 13, 32)),
            Token::Minus(new_token_data("-", 2, 14, 33)),
            Token::Slash(new_token_data("/", 2, 15, 34)),
            Token::Asterisk(new_token_data("*", 2, 16, 35)),
            Token::Semicolon(new_token_data(";", 2, 17, 36)),
            Token::Newline(new_token_data("\n", 2, 38, 57)),
            Token::LessThan(new_token_data("<", 3, 13, 70)),
            Token::GreaterThan(new_token_data(">", 3, 14, 71)),
            Token::Append(new_token_data("<<", 3, 16, 73)),
            Token::Semicolon(new_token_data(";", 3, 18, 75)),
            Token::Newline(new_token_data("\n", 3, 19, 76)),
            Token::Equals(new_token_data("==", 4, 13, 89)),
            Token::NotEquals(new_token_data("!=", 4, 16, 92)),
            Token::GreaterThanOrEquals(new_token_data(">=", 4, 19, 95)),
            Token::LessThanOrEquals(new_token_data("<=", 4, 22, 98)),
            Token::BitAnd(new_token_data("&", 4, 25, 101)),
            Token::BooleanAnd(new_token_data("&&", 4, 27, 103)),
            Token::BitOr(new_token_data("|", 4, 30, 106)),
            Token::BooleanOr(new_token_data("||", 4, 32, 108)),
            Token::BooleanAndAssign(new_token_data("&&=", 4, 35, 111)),
            Token::BooleanOrAssign(new_token_data("||=", 4, 39, 115)),
        ]),
    };

    test.run()
}

#[test]
fn test_lex_literals() {
    let test = LexerTest {
        input: "
            10 10_000 10.50
            \"foobar\" \"foo bar\"
        "
        .to_string(),
        expected: Vec::from([
            Token::Newline(new_token_data("\n", 1, 1, 0)),
            Token::Int(new_token_data("10", 2, 13, 13)),
            Token::Int(new_token_data("10_000", 2, 16, 16)),
            Token::Float(new_token_data("10.50", 2, 23, 23)),
            Token::Newline(new_token_data("\n", 2, 28, 28)),
            Token::String(new_token_data("foobar", 3, 13, 41)),
            Token::String(new_token_data("foo bar", 3, 22, 50)),
            Token::Newline(new_token_data("\n", 3, 31, 59)),
            Token::Eof(new_token_data("\u{0}", 4, 9, 68)),
        ]),
    };

    test.run()
}
