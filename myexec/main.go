package myexec

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
)

func GetExecOutByCSV(argu string) [][]string {

	out, err := exec.Command("/bin/bash", "-c", argu).Output()

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	return records
}

func GetExecOutByString(argu string) string {

	cmd := exec.Command("/bin/bash", "-c", argu)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s", cmd, err.Error())
		return ""
	}

	return string(output)

}
