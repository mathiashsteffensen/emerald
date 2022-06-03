use lazy_static::lazy_static;
use std::env;
use std::sync::Mutex;

lazy_static! {
    static ref TRUE_ENV_VAR_VALUES: Mutex<Vec<String>> = Mutex::new(Vec::from([
        "1".to_string(),
        "true".to_string(),
        "on".to_string()
    ]));
    static ref DEBUG: Mutex<bool> =
        Mutex::new(env_var_is_true("RUBY_DEBUG") || env_var_is_true("EMERALD_DEBUG"));
}

pub(crate) fn log(msg: String) {
    if *DEBUG.lock().unwrap() {
        println!("{}", msg)
    }
}

fn env_var_is_true(var: &str) -> bool {
    let val = env::var(var).unwrap_or("".to_string());

    TRUE_ENV_VAR_VALUES.lock().unwrap().contains(&val)
}
