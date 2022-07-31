package fs

import (
	"context"
	"fmt"
	"os/exec"
	"errors"
	"go-utils/src/logs"
)

// RunSysCommand run shell commands
func RunSysCommand(ctx context.Context, commands []string, envs map[string]string) error {
	logs.Log.Debugf("command: %+v", commands)
	if len(commands) < 1 {
		return errors.New(fmt.Sprintf("para num error, cmd = %+v", commands))
	}
	cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Log.Errorf("command: %+v with error: %v, %+v", commands, string(output), err)
		return err
	}
	return nil
}

// RunSysCommandRet run shell commands and get stdout info
func RunSysCommandRet(ctx context.Context, commands []string, envs map[string]string) ([]byte, error) {
	logs.Log.Debugf("command: %+v", commands)
	if len(commands) < 1 {
		return nil, errors.New(fmt.Sprintf("para num error, cmd = %+v", commands))
	}
	cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
	output, err := cmd.Output()
	if err != nil {
		logs.Log.Debugf("command: %+v with err: %+v", commands, err)
		return nil, err
	}
	return output, nil
}
