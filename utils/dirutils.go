package dirutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

func GetIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func CurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	dir = dir + "/dist"
	var ss []string
	if runtime.GOOS == "windows" {
		ss = strings.Split(dir, "\\")
	} else {
		ss = strings.Split(dir, "/")
	}

	currentDirName := ss[len(ss)-1]

	fmt.Println("Current Directory Name: ", currentDirName)

	return currentDirName
}
