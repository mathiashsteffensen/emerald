use std::collections::HashMap;

#[derive(PartialEq, Debug, Clone)]
pub enum SymbolScope {
    Global,
    Local,
    //     Free,
}

#[derive(PartialEq, Debug, Clone)]
pub struct Symbol {
    pub name: String,
    pub scope: SymbolScope,
    pub index: u16,
}

#[derive(PartialEq, Debug, Clone)]
pub(crate) struct SymbolTable {
    pub(crate) outer: Option<Box<SymbolTable>>,
    // free_symbols: Vec<Symbol>,
    store: HashMap<String, Symbol>,
    num_definitions: u16,
}

impl SymbolTable {
    pub fn new() -> SymbolTable {
        SymbolTable {
            outer: None,
            // free_symbols: Vec::new(),
            store: HashMap::new(),
            num_definitions: 0,
        }
    }

    pub fn with_outer(outer: SymbolTable) -> SymbolTable {
        let mut table = SymbolTable::new();

        table.outer = Some(Box::new(outer));

        table
    }

    pub fn define(&mut self, name: &String) -> Symbol {
        let scope = match self.outer {
            Some(_) => SymbolScope::Local,
            None => SymbolScope::Global,
        };

        let symbol = Symbol {
            scope,
            name: name.clone(),
            index: self.num_definitions,
        };

        self.store.insert(name.clone(), symbol.clone());
        self.num_definitions += 1;

        symbol
    }

    pub fn resolve(&self, name: &String) -> Option<Symbol> {
        if let Some(sym) = self.store.get(name.as_str()) {
            Some(sym.clone())
        } else {
            self.resolve_outer(name)
        }
    }

    fn resolve_outer(&self, name: &String) -> Option<Symbol> {
        if let Some(outer) = &self.outer {
            outer.resolve(name)
        } else {
            None
        }
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

        let a = global.define(&"a".to_string());
        assert_eq!(a, *expected.get("a").unwrap());

        let b = global.define(&"b".to_string());
        assert_eq!(b, *expected.get("b").unwrap());

        let mut local = SymbolTable::with_outer(global);

        let c = local.define(&"c".to_string());
        assert_eq!(c, *expected.get("c").unwrap());
        assert_eq!(
            local.resolve(&"c".to_string()).unwrap(),
            *expected.get("c").unwrap()
        );
    }
}
