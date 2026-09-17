package utils

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// SelectByFzf run fzf command to select a specific option from in
// The in must contains \n so that fzf will work properly
// It will Trim \n\t and space for the output
func SelectByFzf(in io.Reader) (selected string, err error) {
	var selectedBuf bytes.Buffer

	fzfCmd := exec.Command("fzf")

	fzfCmd.Stdin = in

	fzfCmd.Stdout = &selectedBuf

	err = fzfCmd.Run()
	if err != nil {
		err = fmt.Errorf("error when exec fzf command: %s", err)
		return
	}

	return strings.Trim(selectedBuf.String(), "\n\t "), nil
}
