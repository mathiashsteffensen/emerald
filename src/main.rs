use emerald;

#[cfg(not(tarpaulin_include))]
fn main() {
    let mut lexer = emerald::lexer::Lexer::new(emerald::lexer::input::Input::new(
        "main.rb".to_string(),
        "puts(\"Hello World\")".to_string(),
    ));

    println!("{:?}", lexer.next_token())
}
