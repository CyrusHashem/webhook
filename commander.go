package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Commander struct {
	command  string
	notifier Notifier
	queue    chan struct{}
	done     chan struct{}
}

func NewCommander(command string, notifier Notifier) *Commander {
	return &Commander{
		command:  command,
		notifier: notifier,
		queue:    make(chan struct{}, 1),
		done:     make(chan struct{}),
	}
}

func (c *Commander) Exec() {
	select {
	case c.queue <- struct{}{}:
	default:
	}
}

func (c *Commander) Close() {
	close(c.queue)
	<-c.done
}

func (c *Commander) Run() {
	for {
		_, ok := <-c.queue

		if !ok {
			close(c.done)
			return
		}
		c.exec()
	}
}

func (c *Commander) exec() {
	fmt.Println("exec:", c.command)
	cmd := exec.Command("sh", "-c", c.command)
	buf := bytes.NewBuffer(nil)
	w := io.MultiWriter(buf, os.Stdout)
	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Run(); err != nil {
		fmt.Printf("failed to exec: %s: %s\n", c.command, err)
		c.notifier.Notify(NewNotification(c.command, err.Error(), buf.String()))
	}
}
