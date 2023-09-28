package core

import "emerald/object"

var File *object.Class

func InitFile() {
	File = DefineClass("File", IO)
}
