// +build !windows

package tool

import (
	"os/exec"
	"syscall"
)

type CmdResult struct {
	Output string
	Err    error
}

//linux 系统执行定时任务
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	resultChan := make(chan CmdResult)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- CmdResult{Output: string(output), Err: err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return "", errors.New("timeout kill")
	case result := <-resultChan:
		return result.Output, result.Err
	}
}
