#[derive(PartialEq, Debug, Clone)]
pub struct TokenData {
    pub literal: String,
    pub line: i64,
    pub column: i64,
    pub pos: i64,
}

#[derive(PartialEq, Debug, Clone)]
pub enum Token {
    Illegal(TokenData),
    Eof(TokenData),

    // Identifiers
    Ident(TokenData),
    GlobalIdent(TokenData),
    InstanceVar(TokenData),

    // Literals
    Int(TokenData),
    Float(TokenData),
    String(TokenData),

    // Operators
    Assign(TokenData),
    Plus(TokenData),
    Minus(TokenData),
    Bang(TokenData),
    Asterisk(TokenData),
    Slash(TokenData),
    LessThan(TokenData),
    GreaterThan(TokenData),
    Equals(TokenData),
    NotEquals(TokenData),
    LessThanOrEquals(TokenData),
    GreaterThanOrEquals(TokenData),
    Append(TokenData),
    Dot(TokenData),
    BitAnd(TokenData),
    BitOr(TokenData),
    BooleanAnd(TokenData),
    BooleanAndAssign(TokenData),
    BooleanOr(TokenData),
    BooleanOrAssign(TokenData),

    // Delimiters
    Comma(TokenData),
    Semicolon(TokenData),
    Colon(TokenData),
    Arrow(TokenData),
    LeftParen(TokenData),
    RightParen(TokenData),
    LeftBrace(TokenData),
    RightBrace(TokenData),
    LeftBracket(TokenData),
    RightBracket(TokenData),
    Newline(TokenData),

    // Keywords
    Class(TokenData),
    Module(TokenData),
    Def(TokenData),
    Do(TokenData),
    Begin(TokenData),
    Rescue(TokenData),
    Ensure(TokenData),
    End(TokenData),
    True(TokenData),
    False(TokenData),
    Selff(TokenData),
    If(TokenData),
    Else(TokenData),
    Return(TokenData),
    Nil(TokenData),
}

impl Token {
    pub fn data(&self) -> TokenData {
        match self {
            Token::Illegal(data) => data.clone(),
            Token::Eof(data) => data.clone(),
            Token::Ident(data) => data.clone(),
            Token::GlobalIdent(data) => data.clone(),
            Token::InstanceVar(data) => data.clone(),
            Token::Int(data) => data.clone(),
            Token::Float(data) => data.clone(),
            Token::String(data) => data.clone(),
            Token::Assign(data) => data.clone(),
            Token::Plus(data) => data.clone(),
            Token::Minus(data) => data.clone(),
            Token::Bang(data) => data.clone(),
            Token::Asterisk(data) => data.clone(),
            Token::Slash(data) => data.clone(),
            Token::LessThan(data) => data.clone(),
            Token::GreaterThan(data) => data.clone(),
            Token::Equals(data) => data.clone(),
            Token::NotEquals(data) => data.clone(),
            Token::LessThanOrEquals(data) => data.clone(),
            Token::GreaterThanOrEquals(data) => data.clone(),
            Token::Append(data) => data.clone(),
            Token::Dot(data) => data.clone(),
            Token::BitAnd(data) => data.clone(),
            Token::BitOr(data) => data.clone(),
            Token::BooleanAnd(data) => data.clone(),
            Token::BooleanAndAssign(data) => data.clone(),
            Token::BooleanOr(data) => data.clone(),
            Token::BooleanOrAssign(data) => data.clone(),
            Token::Comma(data) => data.clone(),
            Token::Semicolon(data) => data.clone(),
            Token::Colon(data) => data.clone(),
            Token::Arrow(data) => data.clone(),
            Token::LeftParen(data) => data.clone(),
            Token::RightParen(data) => data.clone(),
            Token::LeftBrace(data) => data.clone(),
            Token::RightBrace(data) => data.clone(),
            Token::LeftBracket(data) => data.clone(),
            Token::RightBracket(data) => data.clone(),
            Token::Newline(data) => data.clone(),
            Token::Class(data) => data.clone(),
            Token::Module(data) => data.clone(),
            Token::Def(data) => data.clone(),
            Token::Do(data) => data.clone(),
            Token::Begin(data) => data.clone(),
            Token::Rescue(data) => data.clone(),
            Token::Ensure(data) => data.clone(),
            Token::End(data) => data.clone(),
            Token::True(data) => data.clone(),
            Token::False(data) => data.clone(),
            Token::Selff(data) => data.clone(),
            Token::If(data) => data.clone(),
            Token::Else(data) => data.clone(),
            Token::Return(data) => data.clone(),
            Token::Nil(data) => data.clone(),
        }
    }

    pub fn token_literal(&self) -> String {
        self.data().literal
    }
}
