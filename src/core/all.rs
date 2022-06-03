use crate::core;

pub fn init() {
    // Object hierarchy base
    core::basic_object::em_init_class();
    core::object::em_init_class();

    // Primitives
    core::nil_class::em_init_class();
    core::true_class::em_init_class();
    core::false_class::em_init_class();
    core::integer::em_init_class();
    core::string::em_init_class();
    core::symbol::em_init_class();
}
