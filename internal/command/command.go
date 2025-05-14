package command

import (
	"errors"
	"strings"
)

type Command struct {
	Command string
	args    []string
}

func NewCommand(input string) Command {
	inputArgs := strings.Fields(input)
	command := ""
	args := []string{}

	if len(inputArgs) >= 1 {
		command = inputArgs[0]
	}

	if len(inputArgs) > 1 {
		args = inputArgs[1:]
	}

	return Command{command, args}
}

func (c Command) GetArg(idx int) (string, error) {
	if idx >= 0 && idx < len(c.args) {
		return c.args[idx], nil
	}

	return "", errors.New("index out of range")
}
