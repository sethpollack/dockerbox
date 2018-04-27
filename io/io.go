package io

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func ReadConfig(filename, fileType string) ([]byte, error) {
	var (
		bytes []byte
		err   error
	)

	switch fileType {
	case "file":
		file, err := getFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %v", filename, err)
		}
		defer file.Close()

		bytes, err = ioutil.ReadAll(file)
		if err != nil {
			err = fmt.Errorf("read failed for %s: %v", filename, err)
			return nil, err
		}
	case "url":
		bytes, err = DownloadFile(filename)
		if err != nil {
			err = fmt.Errorf("download failed for %s: %v", filename, err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("format %s not supported", fileType)
	}

	return bytes, nil
}

func ReadFile(filename string) ([]byte, error) {
	file, err := getFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("read failed for %s: %v", filename, err)
		return nil, err
	}

	return bytes, nil
}

func getFile(filename string) (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err == nil {
		return f, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	if err := os.MkdirAll(path.Dir(filename), os.FileMode(0744)); err != nil {
		return nil, err
	}

	return os.Create(filename)
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func EnsureDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return os.MkdirAll(dirname, os.FileMode(0744))
	}

	return nil
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buff bytes.Buffer
	_, err = io.Copy(&buff, resp.Body)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), err
}
