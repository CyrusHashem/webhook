package main

type Notifier interface {
	Notify(*Notification)
}

type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

func (mn *MultiNotifier) Add(notifiers ...Notifier) {
	mn.notifiers = append(mn.notifiers, notifiers...)
}

func (mn *MultiNotifier) Notify(n *Notification) {
	for _, notifier := range mn.notifiers {
		notifier.Notify(n)
	}
}

type Notification struct {
	Command string `json:"command"`
	Error   string `json:"error"`
	Output  string `json:"output"`
}

func NewNotification(cmd, err, out string) *Notification {
	return &Notification{Command: cmd, Error: err, Output: out}
}
