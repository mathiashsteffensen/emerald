package object

type YieldFunc func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue
