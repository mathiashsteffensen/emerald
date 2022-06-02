use lazy_static::lazy_static;
use std::sync::Arc;

use crate::core;
use crate::object::{EmeraldObject, UnderlyingValueType};

const NAME: &str = "FalseClass";

lazy_static! {
    pub static ref EM_FALSE: Arc<EmeraldObject> = em_instance();
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
    EmeraldObject::new_instance_with_underlying_value(NAME, UnderlyingValueType::False)
}
