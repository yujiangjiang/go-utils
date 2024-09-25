package ufile

import (
	"fmt"
	"io/ioutil"
	"os"
)

// WriteLines create and write lines to file
func WriteLines(f string, lines []string) {

	create, err := os.Create(f)
	if err != nil {
		fmt.Println("file create failed")
		return
	}
	defer create.Close()

	str := ""
	for i := 0; i < len(lines); i++ {
		str += lines[i] + "\n"
	}

	err = ioutil.WriteFile(f, []byte(str), 0644)
	if err != nil {
		fmt.Println("write failed", err)
		return
	}
}
