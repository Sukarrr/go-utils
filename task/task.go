package task

import (
	"context"
	"errors"
)

type Task struct {
	command   func() error
	finish    chan error
	panicChan chan interface{}
	running   bool
}

func (t *Task) NewTask(c func() error) *Task {
	return &Task{
		command: c,
		finish:  make(chan error, 1),
		running: false,
	}
}

func (t *Task) Go() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.panicChan <- r
			}
		}()
		t.finish <- t.command()
	}()
}

func (t *Task) Wait(ctx context.Context) error {
	if !t.running {
		t.Go()
	}
	select {
	case err := <-t.finish:
		return err
	case p := <-t.panicChan:
		switch x := p.(type) {
		case string:
			return errors.New(x)
		case error:
			return x
		default:
			return errors.New("unknown panic")
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}
