package task

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func TestTask_Wait(t1 *testing.T) {
	type fields struct {
		command   func() error
		finish    chan error
		panicChan chan interface{}
		running   bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				command: func() error {
					fmt.Println("test command no err")
					return nil
				},
				finish:    make(chan error, 1),
				panicChan: make(chan interface{}, 1),
			},
			args: args{
				context.Background(),
			},
		},
		{
			fields: fields{
				command: func() error {
					fmt.Println("test command has err")
					return errors.New("test err")
				},
				finish:    make(chan error, 1),
				panicChan: make(chan interface{}, 1),
			},
			args: args{
				context.Background(),
			},
			wantErr: true,
		},
		{
			fields: fields{
				command: func() error {
					fmt.Println("test command has panic")
					panic("test panic")
					return nil
				},
				finish:    make(chan error, 1),
				panicChan: make(chan interface{}, 1),
			},
			args: args{
				context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				command:   tt.fields.command,
				finish:    tt.fields.finish,
				panicChan: tt.fields.panicChan,
				running:   tt.fields.running,
			}
			if err := t.Wait(tt.args.ctx); (err != nil) != tt.wantErr {
				t1.Errorf("Wait() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println("wait result: ", err)
			}
		})
	}
	fmt.Println("routine num: ", runtime.NumGoroutine())
}
