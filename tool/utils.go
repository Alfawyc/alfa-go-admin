package tool

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os/exec"
	"strconv"
)

// 比对密码
func CompareHashAndPassword(original string, input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(original), []byte(input))
	if err != nil {
		return false, err
	}
	return true, nil
}

type CmdResult struct {
	Output string
	Err    error
}

//执行shell 命令
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("cmd", "/C", command)
	var resultChan chan CmdResult = make(chan CmdResult)
	go func() {
		out, err := cmd.CombinedOutput()
		resultChan <- CmdResult{Output: string(out), Err: err}
	}()
	select {
	case <-ctx.Done():
		exec.Command("taskill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
		cmd.Process.Kill()
		log.Println("kill task")
		return "", errors.New("timeout kill")
	case result := <-resultChan:
		//todo 输出转utf8
		return result.Output, result.Err
	}
}
