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
	smtpConf   SmtpConf
	mailFrom   string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8000", "address to listen to")
	flag.StringVar(&hookFile, "hook-file", "", "path to the hook file to use")
	flag.StringVar(&smtpConf.Host, "smtp-host", "localhost", "smtp server")
	flag.IntVar(&smtpConf.Port, "smtp-port", 25, "smtp port")
	flag.StringVar(&smtpConf.User, "smtp-user", "", "smtp username")
	flag.StringVar(&smtpConf.Pass, "smtp-pass", "", "smtp password")
	flag.StringVar(&mailFrom, "mail-from", "webhook@localhost.local", "mail from address")
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
		commander := buildCommander(hook)
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
	if len(addr) == 0 || addr[0] != '/' {
		return net.Listen("tcp", addr)
	}
	l, err := net.Listen("unix", addr)

	if err != nil {
		return nil, err
	}

	if err := os.Chmod(addr, 0777); err != nil {
		l.Close()
		return nil, err
	}
	return l, nil
}

func buildCommander(hook *Hook) *Commander {
	notifier := NewMultiNotifier()

	for _, url := range hook.Notify.Web {
		notifier.Add(NewWebNotifier(url))
	}

	for _, email := range hook.Notify.Email {
		notifier.Add(NewEmailNotifier(mailFrom, email, smtpConf))
	}
	return NewCommander(hook.Command, notifier)
}
