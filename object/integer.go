package object

import (
	"fmt"
	"strconv"
)

var Integer *Class

type IntegerInstance struct {
	*Instance
	Value int64
}

func (i *IntegerInstance) Inspect() string {
	return strconv.Itoa(int(i.Value))
}

func NewInteger(val int64) EmeraldValue {
	return &IntegerInstance{Integer.New(), val}
}

var integerBuiltInMethodSet = BuiltInMethodSet{
	"+": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("+", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value + otherVal.Value)
	},
	"-": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("-", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value - otherVal.Value)
	},
	"*": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("*", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value * otherVal.Value)
	},
	"/": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("/", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value / otherVal.Value)
	},
	"<": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("<", args)
		if err != nil {
			return err
		}

		return nativeBoolToBooleanObject(target.(*IntegerInstance).Value < otherVal.Value)
	},
	">": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg(">", args)
		if err != nil {
			return err
		}

		return nativeBoolToBooleanObject(target.(*IntegerInstance).Value > otherVal.Value)
	},
	"==": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("==", args)
		if err != nil {
			return err
		}

		return nativeBoolToBooleanObject(target.(*IntegerInstance).Value == otherVal.Value)
	},
	"!=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return nativeBoolToBooleanObject(target.(*IntegerInstance).Value != otherVal.Value)
	},
	"<=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return nativeBoolToBooleanObject(target.(*IntegerInstance).Value <= otherVal.Value)
	},
	">=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return nativeBoolToBooleanObject(target.(*IntegerInstance).Value <= otherVal.Value)
	},
}

func requireOneIntegerArg(method string, args []EmeraldValue) (*IntegerInstance, EmeraldValue /* StandardError or nil */) {
	if len(args) != 1 {
		return nil, NewStandardError(fmt.Sprintf("Integer#%s expects single argument, got %d", method, len(args)))
	}

	otherVal, ok := args[0].(*IntegerInstance)

	if !ok {
		return nil, NewStandardError(fmt.Sprintf("Integer#%s can only be passed an integer", method))
	}

	return otherVal, nil
}
