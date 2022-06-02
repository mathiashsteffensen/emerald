use std::collections::HashMap;
use std::sync::Arc;

use crate::object::{EmeraldObject, UnderlyingValueType};
use crate::{core, object};

pub const NAME: &str = "BasicObject";

pub fn em_init_class() {
    core::em_define_class(EmeraldObject::new_class(NAME, None, Default::default())).unwrap()
}

pub fn em_instance() -> Arc<EmeraldObject> {
    Arc::from(EmeraldObject {
        class: Some(core::em_get_class(NAME).unwrap()),
        q_super: None,
        built_in_method_set: Default::default(),
        underlying_value: UnderlyingValueType::None,
    })
}
