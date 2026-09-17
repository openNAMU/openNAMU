//go:build windows

package route

import "os/exec"

func Detach_server_process(command *exec.Cmd) {
	_ = command
}
