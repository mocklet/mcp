package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// createMultipartFile creates a multipart request body containing a single file field named "file".
// It returns the content type, the body as an io.Reader, and any error.
func createMultipartFile(filePath string) (string, io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", nil, err
	}

	err = writer.Close()
	if err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body, nil
}

// createMultipartFileWithFields creates a multipart request body containing a file and additional text fields.
func createMultipartFileWithFields(filePath string, fields map[string]string) (string, io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", nil, err
	}

	for k, v := range fields {
		err = writer.WriteField(k, v)
		if err != nil {
			return "", nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body, nil
}

