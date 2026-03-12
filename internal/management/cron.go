package management

import (
	"bytes"
	"fmt"
	"os/exec"
)

func AddCronJob(schedule, command string) error {
	var output bytes.Buffer

	cmd := exec.Command("crontab", "-l")
	cmd.Stdout = &output

	err := cmd.Run()
	if err != nil {
		output.WriteString("")
	}

	cron := fmt.Sprintln(output.String(), schedule, " ", command)
	cmd = exec.Command("crontab", "-")
	cmd.Stdin = bytes.NewBufferString(cron)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
