use std::sync::Arc;

use crate::core;
use crate::object::EmeraldObject;

pub const NAME: &str = "Object";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::basic_object::NAME),
        Default::default(),
    ))
    .unwrap()
}

pub fn em_instance() -> Arc<EmeraldObject> {
    EmeraldObject::new_instance(NAME)
}
