use lazy_static::lazy_static;
use std::env;
use std::sync::Mutex;
use std::time::SystemTime;

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
    if is_debug() {
        println!("{}", msg)
    }
}

pub(crate) fn time<CB, CBReturns, Formatter>(mut cb: CB, formatter: Formatter) -> CBReturns
where
    CB: FnMut() -> CBReturns,
    Formatter: Fn(u128) -> String,
{
    if !is_debug() {
        return cb();
    }

    let start = SystemTime::now();

    let result = cb();

    match start.elapsed() {
        Ok(elapsed) => log(formatter(elapsed.as_millis())),
        Err(e) => log(format!("error timing {:?}", e)),
    }

    result
}

fn env_var_is_true(var: &str) -> bool {
    let val = env::var(var).unwrap_or("".to_string());

    TRUE_ENV_VAR_VALUES.lock().unwrap().contains(&val)
}

fn is_debug() -> bool {
    *DEBUG.lock().unwrap()
}
