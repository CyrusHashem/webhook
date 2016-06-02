package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Commander struct {
	command string
	queue   chan struct{}
	done    chan struct{}
}

func NewCommander(command string) *Commander {
	return &Commander{
		command: command,
		queue:   make(chan struct{}, 1),
		done:    make(chan struct{}),
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
		fmt.Println("exec:", c.command)

		if err := c.exec(); err != nil {
			fmt.Printf("failed to exec: %s: %s\n", c.command, err)
		}
	}
}

func (c *Commander) exec() error {
	out, err := exec.Command("sh", "-c", c.command).CombinedOutput()
	os.Stdout.Write(out)
	return err
}
