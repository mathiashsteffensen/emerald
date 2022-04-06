package lexer

type Input struct {
	fileName string
	content  string
}

func NewInput(fileName string, content string) *Input {
	return &Input{fileName: fileName, content: content}
}
