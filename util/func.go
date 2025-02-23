package util

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

// 获取当前平台
//
// const (
//
//	WinPlEum   uint = 1 // windows
//	LinuxPlEum uint = 2 // linux
//
// )
func ThisPlatform() (uint, error) {
	sysType := runtime.GOOS
	switch sysType {
	case "windows":
		return WinPlEum, nil
	case "linux":
		return LinuxPlEum, nil
	}
	return 0, errors.New("not windows or linux")
}

func InitConfigFiles(linuxPath, windowsPath, defValue string) (string, error) {
	fmt.Printf("linuxPath: %v\n", linuxPath)
	fmt.Printf("windowsPath: %v\n", windowsPath)
	// 前置校验
	pl, err := ThisPlatform()
	if err != nil {
		return "", err
	}
	switch pl {
	case LinuxPlEum:
		isFile, err := PathExists(linuxPath)
		if err != nil {
			return "", err
		}
		if isFile {
			return linuxPath, nil
		} else {
			err := os.WriteFile(linuxPath, []byte(defValue), 0777)
			if err != nil {
				return linuxPath, err
			}
			return linuxPath, nil
		}
	case WinPlEum:
		isFile, err := PathExists(windowsPath)
		if err != nil {
			return "", err
		}
		if isFile {
			return windowsPath, nil
		} else {
			err := os.WriteFile(windowsPath, []byte(defValue), 0777)
			if err != nil {
				return windowsPath, err
			}
			return windowsPath, nil
		}
	}
	return "", errors.New("init error： not run platform")
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()

}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
