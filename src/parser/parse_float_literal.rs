use crate::ast::node::Expression;
use crate::lexer::token;
use crate::parser::Parser;

pub fn exec(p: &mut Parser, data: token::TokenData) -> Option<Expression> {
    let parse_result = data.literal.parse::<f64>();
    match parse_result {
        Ok(val) => Some(Expression::FloatLiteral(data, val)),
        Err(_err) => {
            p.add_error(&*format!("failed to parse float value '{}'", data.literal));

            None
        }
    }
}
