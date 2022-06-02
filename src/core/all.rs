use crate::core;

pub fn init() {
    // Object hierarchy base
    core::basic_object::em_init_class();
    core::object::em_init_class();

    // Primitives
    core::integer::em_init_class();
    core::string::em_init_class();
    core::symbol::em_init_class();
}
