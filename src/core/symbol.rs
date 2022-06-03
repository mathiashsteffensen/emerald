use crate::core;
use crate::object::{EmeraldObject, UnderlyingValueType};
use std::sync::Arc;

pub const NAME: &str = "Symbol";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::object::NAME),
        Default::default(),
    ))
    .unwrap()
}

pub fn em_instance(val: String) -> Arc<EmeraldObject> {
    Arc::from(EmeraldObject {
        class: Some(core::em_get_class(NAME).unwrap()),
        q_super: None,
        built_in_method_set: Default::default(),
        underlying_value: UnderlyingValueType::Symbol(val),
    })
}
