//go:build windows

package route

import "os/exec"

func detach_server_process(command *exec.Cmd) {
	_ = command
}
