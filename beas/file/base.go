package file

import (
	"errors"
	"fmt"
	"mixed/beas/setting"
	"os"
	"strings"
)

var ErrFd *os.File = genDataFile("err")
var OriFd *os.File = genDataFile("ori")
var NewFd *os.File = genDataFile("new")
var LogFd *os.File = genDataFile("log")

// generate data log file
func genDataFile(flag string) *os.File {
	var fd *os.File
	var err error
	switch flag {
	case "err":
		fd, err = CreateFile(setting.FILE_ERR__LOG_DATA)
	case "new":
		fd, err = CreateFile(setting.FILE_NEW_SQL_DATA)
	case "ori":
		fd, err = CreateFile(setting.FILE_ORIGINAL_DATA)
	case "log":
		fd, err = CreateFile(setting.FILE_LOG_DATA)
	default:
		err = errors.New("Wrong file flag")

	}
	if err != nil {
		fmt.Printf("Create log data file err msg:" + err.Error())
		os.Exit(1)
	}
	return fd
}

// close all the opend file when all the worker finished
func CloseRelatedFile(fileDone <-chan bool) {
	select {
	case <-fileDone:
		ErrFd.Close()
		OriFd.Close()
		NewFd.Close()
		LogFd.Close()
	}
}

//write data
func WriteDataFile(fd *os.File, format string, v ...interface{}) (n int, err error) {
	var content string
	if len(format) > 0 {
		content = fmt.Sprintf(format, v...)
	} else {
		content = fmt.Sprint(v...)
	}
	return fd.WriteString(content)
}

// make directory
func mkDir(fileName string) {
	paths := strings.Split(fileName, string(os.PathSeparator))
	paths = paths[0 : len(paths)-1]
	if err := os.MkdirAll(strings.Join(paths, string(os.PathSeparator)), 0777); err != nil {
		panic(err)
	}
}

// create file
func CreateFile(fileName string) (fd *os.File, err error) {
	fd, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return nil, err
		}
		if os.IsNotExist(err) {
			mkDir(fileName)
		}
		fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_SYNC, 0664)
		if err != nil {
			return nil, err
		}
	}
	return fd, err
}
