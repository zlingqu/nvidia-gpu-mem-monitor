package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
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

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return "error"
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return "error"
	}

	//读取所有输出
	myByte, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return "error"
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return "error"
	}
	// fmt.Printf("stdout:\n\n %s", bytes)
	return string(myByte[0 : len(myByte)-1])
}
