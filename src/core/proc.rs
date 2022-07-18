use std::sync::Arc;

use crate::core;
use crate::object::{Block, EmeraldObject, UnderlyingValueType};

const NAME: &str = "Proc";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(
        NAME,
        core::em_get_class(core::object::NAME),
        Default::default(),
    ))
    .unwrap()
}

pub fn em_instance(block: Block) -> Arc<EmeraldObject> {
    EmeraldObject::new_instance_with_underlying_value(NAME, UnderlyingValueType::Proc(block))
}
