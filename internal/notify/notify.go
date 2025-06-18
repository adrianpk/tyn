package notify

import (
	"fmt"
	"os/exec"
	"time"
)

type Notifier interface {
	Notify(title string, message string) error
	SetTimeout(timeout time.Duration)
}

type LinuxNotifier struct {
	timeout int
}

func NewLinuxNotifier() *LinuxNotifier {
	return &LinuxNotifier{
		timeout: 5000,
	}
}

func (n *LinuxNotifier) Notify(title string, message string) error {
	cmd := exec.Command("notify-send",
		"-t", fmt.Sprintf("%d", n.timeout),
		"-a", "Tyn",
		title, message)

	return cmd.Run()
}

func (n *LinuxNotifier) SetTimeout(timeout time.Duration) {
	n.timeout = int(timeout.Milliseconds())
}

func NotifyDaily() error {
	notifier := NewLinuxNotifier()
	return notifier.Notify("Tyn Journal",
		"Your daily journal has been generated")
}

func NotifyTaskReminder(taskCount int) error {
	notifier := NewLinuxNotifier()

	message := "You have no pending tasks"
	if taskCount == 1 {
		message = "You have 1 pending task"
	} else if taskCount > 1 {
		message = fmt.Sprintf("You have %d pending tasks", taskCount)
	}

	return notifier.Notify("Tyn Tasks", message)
}

func NotifyDueDate(taskTitle, message string) error {
	notifier := NewLinuxNotifier()
	return notifier.Notify("Tyn: Task Overdue", message)
}

func NotifyDueDateReminder(taskTitle, message string) error {
	notifier := NewLinuxNotifier()
	return notifier.Notify("Tyn: Overdue Task Reminder", message)
}
