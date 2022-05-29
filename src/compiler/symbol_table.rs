use std::collections::HashMap;

#[derive(PartialEq, Debug, Clone)]
enum SymbolScope {
    Global,
    Local,
    Free,
}

#[derive(PartialEq, Debug, Clone)]
struct Symbol {
    name: String,
    scope: SymbolScope,
    index: u16,
}

struct SymbolTable {
    outer: Option<Box<SymbolTable>>,
    free_symbols: Vec<Symbol>,
    store: HashMap<String, Symbol>,
    num_definitions: u16,
}

impl SymbolTable {
    fn new() -> SymbolTable {
        SymbolTable {
            outer: None,
            free_symbols: Vec::new(),
            store: HashMap::new(),
            num_definitions: 0,
        }
    }

    fn with_outer(outer: SymbolTable) -> SymbolTable {
        let mut table = SymbolTable::new();

        table.outer = Some(Box::new(outer));

        table
    }

    pub fn define(&mut self, name: String) -> Symbol {
        let scope = match self.outer {
            Some(_) => SymbolScope::Local,
            None => SymbolScope::Global,
        };

        let symbol = Symbol {
            scope,
            name: name.clone(),
            index: self.num_definitions,
        };

        self.store.insert(name, symbol.clone());
        self.num_definitions += 1;

        symbol
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn new_symbol(n: &str, i: u16, s: SymbolScope) -> Symbol {
        Symbol {
            name: n.to_string(),
            index: i,
            scope: s,
        }
    }

    fn new_local_symbol(n: &str, i: u16) -> Symbol {
        new_symbol(n, i, SymbolScope::Local)
    }

    fn new_global_symbol(n: &str, i: u16) -> Symbol {
        new_symbol(n, i, SymbolScope::Global)
    }

    #[test]
    fn test_define() {
        let expected: HashMap<String, Symbol> = HashMap::from([
            ("a".to_string(), new_global_symbol("a", 0)),
            ("b".to_string(), new_global_symbol("b", 1)),
            ("c".to_string(), new_local_symbol("c", 0)),
        ]);

        let mut global = SymbolTable::new();

        let a = global.define("a".to_string());
        assert_eq!(a, *expected.get("a").unwrap());

        let b = global.define("b".to_string());
        assert_eq!(b, *expected.get("b").unwrap());

        let mut local = SymbolTable::with_outer(global);

        let c = local.define("c".to_string());
        assert_eq!(c, *expected.get("c").unwrap())
    }
}
