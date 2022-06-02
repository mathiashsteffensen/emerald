use lazy_static::lazy_static;
use std::sync::Arc;

use crate::core;
use crate::object::{EmeraldObject, UnderlyingValueType};

const NAME: &str = "TrueClass";

lazy_static! {
    static ref EM_TRUE: Arc<EmeraldObject> = em_instance();
}

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::object::NAME),
        Default::default(),
    ))
    .unwrap()
}

pub fn em_instance() -> Arc<EmeraldObject> {
    Arc::from(EmeraldObject {
        class: Some(core::em_get_class(NAME).unwrap()),
        q_super: None,
        built_in_method_set: Default::default(),
        underlying_value: UnderlyingValueType::True,
    })
}
