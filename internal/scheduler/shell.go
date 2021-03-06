package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cybertec-postgresql/pg_timetable/internal/pgengine"
)

type commander interface {
	CombinedOutput(context.Context, string, ...string) ([]byte, error)
}

type realCommander struct{}

func (c realCommander) CombinedOutput(ctx context.Context, command string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, command, args...).CombinedOutput()
}

// Cmd executes a command
var Cmd commander = realCommander{}

// ExecuteShellCommand executes program command and returns status code, output and error if any
func ExecuteProgramCommand(ctx context.Context, command string, paramValues []string) (code int, stdout string, stderr error) {

	if strings.TrimSpace(command) == "" {
		return -1, "", errors.New("Program command cannot be empty")
	}
	if len(paramValues) == 0 { //mimic empty param
		paramValues = []string{""}
	}
	for _, val := range paramValues {
		params := []string{}
		if val > "" {
			if err := json.Unmarshal([]byte(val), &params); err != nil {
				return -1, "", err
			}
		}
		out, err := Cmd.CombinedOutput(ctx, command, params...) // #nosec
		cmdLine := fmt.Sprintf("%s %v: ", command, params)
		stdout = strings.TrimSpace(string(out))
		if len(out) > 0 {
			pgengine.LogToDB(ctx, "DEBUG", "Output for command ", cmdLine, string(out))
		}
		if err != nil {
			//check if we're dealing with an ExitError - i.e. return code other than 0
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode := exitError.ProcessState.ExitCode()
				pgengine.LogToDB(ctx, "DEBUG", "Return value of the command ", cmdLine, exitCode)
				return exitCode, stdout, exitError
			}
			return -1, stdout, err
		}
	}
	return 0, stdout, nil
}
