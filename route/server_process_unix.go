//go:build !windows

package route

import (
	"os/exec"
	"syscall"
)

func detach_server_process(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
