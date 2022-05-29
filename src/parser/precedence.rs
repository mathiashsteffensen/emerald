use crate::lexer;

pub const LOWEST: i16 = 0;
pub const EQUALS: i16 = 1;
pub const LESS_GREATER: i16 = 2;
pub const SUM: i16 = 3;
pub const PRODUCT: i16 = 4;
pub const PREFIX: i16 = 5;
pub const CALL: i16 = 6;
pub const ACCESSOR: i16 = 7;

pub fn precedence_for(token: lexer::token::Token) -> i16 {
    match token {
        lexer::token::Token::Equals(_data) => EQUALS,
        lexer::token::Token::NotEquals(_data) => EQUALS,
        lexer::token::Token::LessThan(_data) => LESS_GREATER,
        lexer::token::Token::LessThanOrEquals(_data) => LESS_GREATER,
        lexer::token::Token::GreaterThan(_data) => LESS_GREATER,
        lexer::token::Token::GreaterThanOrEquals(_data) => LESS_GREATER,
        lexer::token::Token::Plus(_data) => SUM,
        lexer::token::Token::Minus(_data) => SUM,
        lexer::token::Token::Slash(_data) => PRODUCT,
        lexer::token::Token::Asterisk(_data) => PRODUCT,
        _ => LOWEST,
    }
}
