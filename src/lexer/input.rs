pub struct Input {
    pub file_name: String,
    pub content: String,
}

impl Input {
    pub fn new(file_name: String, content: String) -> Input {
        return Input { file_name, content };
    }
}
