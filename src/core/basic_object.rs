use std::sync::Arc;

use crate::core;
use crate::object::EmeraldObject;

pub const NAME: &str = "BasicObject";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(NAME, None, Default::default())).unwrap()
}

pub fn em_instance() -> Arc<EmeraldObject> {
    EmeraldObject::new_instance(NAME)
}
