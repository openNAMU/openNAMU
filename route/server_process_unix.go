//go:build !windows

package route

import (
	"os/exec"
	"syscall"
)

func Detach_server_process(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
