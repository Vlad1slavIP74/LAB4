package commands

import (
	"../engine"
	"fmt"
	"strings"
)

type printCommand struct {
	arg string
}

type concatCommand struct {
	arg1, arg2 string
}

func (concat *concatCommand) Execute(loop engine.Handler) {
	var res []string
	res = append(res, concat.arg1, concat.arg2)
	fmt.Println(strings.Join(res, ""))
}

func (p printCommand) Execute(loop engine.Handler) {
	fmt.Println(p.arg)
}

func Parse(commandLine string) engine.Command {
	args := strings.Fields(commandLine)
	switch {
	case len(args) < 1:
		return &printCommand{"SYNTAX ERROR: no commands"}
	case args[0] == "print" && len(args) == 2:
		return &printCommand{arg: args[1]}
	case args[0] == "cat" && len(args) == 3:
		return &concatCommand{arg1: args[1], arg2: args[2]}
	case ((args[0] == "print" && len(args) < 2) || (args[0] == "cat" && len(args) < 3)):
		return &printCommand{"SYNTAX ERROR: not enough arguments"}
	case ((args[0] == "print" && len(args) > 2) || (args[0] == "cat" && len(args) > 3)):
		return &printCommand{"SYNTAX ERROR: too much arguments"}
	default:
		return &printCommand{"SYNTAX ERROR: invalid command"}
	}
}
