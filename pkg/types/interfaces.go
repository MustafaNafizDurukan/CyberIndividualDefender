package types

type IChecker interface {
	Init()
	Descriptor() *Descriptor
	Check() error
}
