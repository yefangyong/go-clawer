package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func MapToString(m map[string]string, split1 string, split2 string) string {
	var str string
	for key, val := range m {
		str = fmt.Sprintf("%s%s%s%s%s", str, key, split2, val, split1)
	}
	return str[0 : len(str)-1]
}

func JsonPToJson(jsonp string) string {
	if jsonp[0] != '[' && jsonp[0] != '{' {
		length := len(jsonp[0:strings.Index(jsonp, "(")])
		jsonp = jsonp[length:]
	}
	return strings.Trim(jsonp, "();")
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return path[0 : i+1], nil
}

func MakeLogFile(logPath string, logFolder string) string {
	var logFullPath string
	if runtime.GOOS == "windows" {
		currentPath, err := GetCurrentPath()
		if err != nil {
			panic(err)
		}
		logFullPath = fmt.Sprintf("%s/%s/%s", currentPath, logPath, logFolder)
	} else {
		logFullPath = fmt.Sprintf("%s/%s", logPath, logFolder)
	}
	if err := os.MkdirAll(logFullPath, 0775); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s.log", logFullPath, logFolder)
}

func MD5(str string) string {
	write := md5.New()
	write.Write([]byte(str))
	return hex.EncodeToString(write.Sum(nil))
}
