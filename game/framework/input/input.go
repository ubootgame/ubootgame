package input

import "github.com/samber/do"

type Input interface {
	Cursor() *Cursor
}

type input struct {
	cursor *Cursor
}

func NewInput(_ *do.Injector) (Input, error) {
	return &input{cursor: NewCursor()}, nil
}

func (input *input) Cursor() *Cursor {
	return input.cursor
}

func (input *input) Shutdown() error {
	println("input stopped")
	return nil
}
