package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	listenAddr string
	hookFile   string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8000", "address to listen to")
	flag.StringVar(&hookFile, "hook-file", "", "path to the hook file to use")
}

func main() {
	flag.Parse()

	if len(hookFile) == 0 {
		fmt.Println("a hook file (-hook-file) is required")
		os.Exit(1)
	}
	hooks, err := ReadHookFile(hookFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l, err := listen(listenAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mux := http.NewServeMux()

	for _, hook := range hooks {
		commander := NewCommander(hook.Command)
		go commander.Run()
		defer commander.Close()
		handler := NewHandler(hook.Hook, hook.Auth, commander)
		mux.Handle(fmt.Sprintf("/%s", hook.Hook), handler)
	}
	go http.Serve(l, mux)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sig

	l.Close()
}

func listen(addr string) (net.Listener, error) {
	if len(addr) > 0 && addr[0] == '/' {
		return net.Listen("unix", addr)
	}
	return net.Listen("tcp", addr)
}
