package watcher

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatch(t *testing.T) {
	t.Run("edit", func(t *testing.T) {
		t.Parallel()
		var (
			ch     chan struct{}
			called bool
		)
		ch = make(chan struct{})
		f, _ := ioutil.TempFile(".", "*")
		defer os.Remove(f.Name())

		ioutil.WriteFile(f.Name(), []byte(`foo`), os.ModePerm)

		w := File{f.Name()}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go w.Watch(ctx, func() error {
			called = true
			ch <- struct{}{}
			return nil
		})
		time.Sleep(time.Second)
		ioutil.WriteFile(f.Name(), []byte(`bar`), os.ModePerm)
		<-ch
		assert.True(t, called)
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()
		var (
			ch     chan struct{}
			called bool
		)
		ch = make(chan struct{})
		f, _ := ioutil.TempFile(".", "*")

		ioutil.WriteFile(f.Name(), []byte(`foo`), os.ModePerm)

		w := File{f.Name()}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			w.Watch(ctx, func() error {
				called = true
				ch <- struct{}{}
				return nil
			})
			ch <- struct{}{}
		}()
		time.Sleep(time.Second)
		os.Remove(f.Name())
		<-ch
		assert.False(t, called)
	})

	t.Run("reload failed", func(t *testing.T) {
		t.Parallel()
		var (
			ch     chan struct{}
			called bool
		)
		ch = make(chan struct{})
		f, _ := ioutil.TempFile(".", "*")

		ioutil.WriteFile(f.Name(), []byte(`foo`), os.ModePerm)
		defer os.Remove(f.Name())

		w := File{f.Name()}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			w.Watch(ctx, func() error {
				return errors.New("foo")
			})
			called = true
			ch <- struct{}{}
		}()
		time.Sleep(time.Second)
		ioutil.WriteFile(f.Name(), []byte(`bar`), os.ModePerm)
		<-ch
		assert.True(t, called)
	})

}
