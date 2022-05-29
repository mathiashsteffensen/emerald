use std::str;

mod char_classifiers;

pub mod input;
pub mod token;

pub struct Lexer {
    input: input::Input,
    line: i64,
    column: i64,
    position: i64, // current position in input (points to current char)
    next_position: i64,
    current_char: char, // current char under examination
}

impl Lexer {
    pub fn new(input: input::Input) -> Lexer {
        let mut lexer = Lexer {
            input,
            line: 1,
            column: 0,
            position: 0,
            next_position: 0,
            current_char: char::from(0),
        };

        lexer.read_char();

        return lexer;
    }

    pub fn next_token(&mut self) -> token::Token {
        self.eat_whitespace();

        let mut tok: token::Token = self.new_illegal_token();

        match self.current_char {
            '\u{0}' => tok = token::Token::Eof(self.new_token_data_with_current_char()),
            '#' => {
                while self.current_char != '\n' {
                    self.read_char();
                }
                tok = self.next_token()
            }
            '\n' => {
                tok = token::Token::Newline(self.new_token_data_with_current_char());
                self.line += 1;
                self.column = 0;
            }
            ';' => tok = token::Token::Semicolon(self.new_token_data_with_current_char()),
            '(' => tok = token::Token::LeftParen(self.new_token_data_with_current_char()),
            ')' => tok = token::Token::RightParen(self.new_token_data_with_current_char()),
            '[' => tok = token::Token::LeftBracket(self.new_token_data_with_current_char()),
            ']' => tok = token::Token::RightBracket(self.new_token_data_with_current_char()),
            '{' => tok = token::Token::LeftBrace(self.new_token_data_with_current_char()),
            '}' => tok = token::Token::RightBrace(self.new_token_data_with_current_char()),
            ',' => tok = token::Token::Comma(self.new_token_data_with_current_char()),
            '.' => tok = token::Token::Dot(self.new_token_data_with_current_char()),
            '=' => {
                if self.peek_char() == '=' {
                    tok = token::Token::Equals(self.new_token_data_with_current_and_peek_char());
                } else {
                    tok = token::Token::Assign(self.new_token_data_with_current_char());
                }
            }
            '!' => {
                if self.peek_char() == '=' {
                    tok = token::Token::NotEquals(self.new_token_data_with_current_and_peek_char());
                } else {
                    tok = token::Token::Bang(self.new_token_data_with_current_char());
                }
            }
            '&' => {
                if self.peek_char() == '&' {
                    let mut tok_data = self.new_token_data_with_current_and_peek_char();

                    if self.peek_char() == '=' {
                        self.read_char();
                        tok_data.literal.push(self.current_char);

                        tok = token::Token::BooleanAndAssign(tok_data);
                    } else {
                        tok = token::Token::BooleanAnd(tok_data);
                    }
                } else {
                    tok = token::Token::BitAnd(self.new_token_data_with_current_char());
                }
            }
            '|' => {
                if self.peek_char() == '|' {
                    let mut tok_data = self.new_token_data_with_current_and_peek_char();

                    if self.peek_char() == '=' {
                        self.read_char();
                        tok_data.literal.push(self.current_char);

                        tok = token::Token::BooleanOrAssign(tok_data);
                    } else {
                        tok = token::Token::BooleanOr(tok_data);
                    }
                } else {
                    tok = token::Token::BitOr(self.new_token_data_with_current_char());
                }
            }
            '+' => tok = token::Token::Plus(self.new_token_data_with_current_char()),
            '-' => tok = token::Token::Minus(self.new_token_data_with_current_char()),
            '/' => tok = token::Token::Slash(self.new_token_data_with_current_char()),
            '*' => tok = token::Token::Asterisk(self.new_token_data_with_current_char()),
            '<' => {
                if self.peek_char() == '<' {
                    tok = token::Token::Append(self.new_token_data_with_current_and_peek_char());
                } else if self.peek_char() == '=' {
                    tok = token::Token::LessThanOrEquals(
                        self.new_token_data_with_current_and_peek_char(),
                    );
                } else {
                    tok = token::Token::LessThan(self.new_token_data_with_current_char());
                }
            }
            '>' => {
                if self.peek_char() == '=' {
                    tok = token::Token::GreaterThanOrEquals(
                        self.new_token_data_with_current_and_peek_char(),
                    );
                } else {
                    tok = token::Token::GreaterThan(self.new_token_data_with_current_char());
                }
            }
            '$' => {
                return self.lex_global_ident();
            }
            '@' => {
                return self.lex_instance_var();
            }
            '"' => {
                return self.lex_string();
            }
            _ => {
                if char_classifiers::is_letter(self.current_char) {
                    return self.lex_identifier();
                } else if char_classifiers::is_digit(self.current_char) {
                    return self.lex_numeric();
                }
            }
        }

        self.read_char();

        tok
    }

    fn lex_identifier(&mut self) -> token::Token {
        let mut tok_data = self.new_token_data("".to_string());

        tok_data.literal = self.read_identifier();

        match tok_data.literal.as_str() {
            "class" => token::Token::Class(tok_data),
            "module" => token::Token::Module(tok_data),
            "def" => token::Token::Def(tok_data),
            "do" => token::Token::Do(tok_data),
            "begin" => token::Token::Begin(tok_data),
            "rescue" => token::Token::Rescue(tok_data),
            "ensure" => token::Token::Ensure(tok_data),
            "if" => token::Token::If(tok_data),
            "else" => token::Token::Else(tok_data),
            "end" => token::Token::End(tok_data),
            "self" => token::Token::Selff(tok_data),
            "return" => token::Token::Return(tok_data),
            "true" => token::Token::True(tok_data),
            "false" => token::Token::False(tok_data),
            "nil" => token::Token::Nil(tok_data),
            _ => token::Token::Ident(tok_data),
        }
    }

    fn lex_numeric(&mut self) -> token::Token {
        let mut tok_data = self.new_token_data("".to_string());

        tok_data.literal = self.read_number();

        if self.current_char == '.' {
            self.read_char();

            let decimals = self.read_number();

            tok_data.literal.push('.');
            tok_data.literal.push_str(&*decimals);

            token::Token::Float(tok_data)
        } else {
            token::Token::Int(tok_data)
        }
    }

    fn lex_instance_var(&mut self) -> token::Token {
        let mut tok_data = self.new_token_data_with_current_char();

        self.read_char();

        tok_data.literal.push_str(&*self.read_identifier());

        token::Token::InstanceVar(tok_data)
    }

    fn lex_global_ident(&mut self) -> token::Token {
        let mut tok_data = self.new_token_data_with_current_char();

        self.read_char();

        if self.current_char == ':' {
            tok_data.literal.push(self.current_char);
            self.read_char()
        } else {
            tok_data.literal.push_str(&*self.read_identifier());
        }

        token::Token::GlobalIdent(tok_data)
    }

    fn lex_string(&mut self) -> token::Token {
        let mut tok_data = self.new_token_data("".to_string());

        let position = self.next_position;

        self.read_char();

        while self.current_char != '"' {
            self.read_char();
        }

        let string_bytes =
            &self.input.content.as_bytes()[position as usize..self.position as usize];

        tok_data
            .literal
            .push_str(str::from_utf8(string_bytes).unwrap());

        self.read_char();

        token::Token::String(tok_data)
    }

    fn read_char(&mut self) {
        self.current_char = self.peek_char();

        self.position = self.next_position;
        self.next_position += 1;
        self.column += 1;
    }

    fn peek_char(&self) -> char {
        if self.next_position >= self.input.content.len().try_into().unwrap() {
            char::from(0)
        } else {
            self.input
                .content
                .chars()
                .nth(self.next_position as usize)
                .unwrap()
        }
    }

    fn eat_whitespace(&mut self) {
        while self.current_char == ' ' || self.current_char == '\t' || self.current_char == '\r' {
            self.read_char()
        }
    }

    fn read_identifier(&mut self) -> String {
        let position = self.position;

        while char_classifiers::is_letter(self.current_char) {
            self.read_char()
        }

        if self.current_char == '?' || self.current_char == '!' {
            self.read_char()
        }

        let ident_bytes = &self.input.content.as_bytes()[position as usize..self.position as usize];

        String::from(str::from_utf8(ident_bytes).unwrap())
    }

    fn read_number(&mut self) -> String {
        let position = self.position;

        while char_classifiers::is_digit(self.current_char) || self.current_char == '_' {
            self.read_char();
        }

        let ident_bytes = &self.input.content.as_bytes()[position as usize..self.position as usize];

        String::from(str::from_utf8(ident_bytes).unwrap())
    }

    fn new_illegal_token(&self) -> token::Token {
        token::Token::Illegal(self.new_token_data(self.current_char.to_string()))
    }

    fn new_token_data(&self, literal: String) -> token::TokenData {
        token::TokenData {
            literal,
            line: self.line,
            column: self.column,
            pos: self.position,
        }
    }

    fn new_token_data_with_current_char(&self) -> token::TokenData {
        self.new_token_data(self.current_char.to_string())
    }

    fn new_token_data_with_current_and_peek_char(&mut self) -> token::TokenData {
        let mut tok_data = self.new_token_data_with_current_char();

        self.read_char();

        tok_data.literal.push(self.current_char);

        tok_data
    }
}
