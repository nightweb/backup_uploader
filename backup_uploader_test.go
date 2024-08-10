package main

import (
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"io"
)

type DriveService interface {
	FilesCreate(file *drive.File) DriveServiceCreateCall
	FilesUpdate(fileId string, file *drive.File) DriveServiceUpdateCall
	FilesList() DriveServiceListCall
}

type DriveServiceCreateCall interface {
	SupportsAllDrives(supportsAllDrives bool) DriveServiceCreateCall
	Fields(fields ...googleapi.Field) DriveServiceCreateCall
	Do(opts ...googleapi.CallOption) (*drive.File, error)
}

type DriveServiceUpdateCall interface {
	SupportsAllDrives(supportsAllDrives bool) DriveServiceUpdateCall
	Media(r io.Reader, options ...googleapi.MediaOption) DriveServiceUpdateCall
	Do(opts ...googleapi.CallOption) (*drive.File, error)
}

type DriveServiceListCall interface {
	Q(query string) DriveServiceListCall
	Fields(fields ...googleapi.Field) DriveServiceListCall
	SupportsAllDrives(supportsAllDrives bool) DriveServiceListCall
	IncludeItemsFromAllDrives(includeItemsFromAllDrives bool) DriveServiceListCall
	Do(opts ...googleapi.CallOption) (*drive.FileList, error)
}

type MockDriveService struct {
	files map[string]*drive.File
}

func (m *MockDriveService) FilesCreate(file *drive.File) DriveServiceCreateCall {
	return &MockDriveServiceCreateCall{file: file, service: m}
}

func (m *MockDriveService) FilesUpdate(fileId string, file *drive.File) DriveServiceUpdateCall {
	return &MockDriveServiceUpdateCall{fileId: fileId, file: file, service: m}
}

func (m *MockDriveService) FilesList() DriveServiceListCall {
	return &MockDriveServiceListCall{service: m}
}

type MockDriveServiceCreateCall struct {
	file    *drive.File
	service *MockDriveService
}

func (c *MockDriveServiceCreateCall) SupportsAllDrives(supportsAllDrives bool) DriveServiceCreateCall {
	return c
}

func (c *MockDriveServiceCreateCall) Fields(fields ...googleapi.Field) DriveServiceCreateCall {
	return c
}

func (c *MockDriveServiceCreateCall) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	c.service.files[c.file.Name] = c.file
	return c.file, nil
}

type MockDriveServiceUpdateCall struct {
	fileId  string
	file    *drive.File
	service *MockDriveService
}

func (c *MockDriveServiceUpdateCall) SupportsAllDrives(supportsAllDrives bool) DriveServiceUpdateCall {
	return c
}

func (c *MockDriveServiceUpdateCall) Media(r io.Reader, options ...googleapi.MediaOption) DriveServiceUpdateCall {
	return c
}

func (c *MockDriveServiceUpdateCall) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	c.service.files[c.fileId] = c.file
	return c.file, nil
}

type MockDriveServiceListCall struct {
	service *MockDriveService
}

func (c *MockDriveServiceListCall) Q(query string) DriveServiceListCall {
	return c
}

func (c *MockDriveServiceListCall) Fields(fields ...googleapi.Field) DriveServiceListCall {
	return c
}

func (c *MockDriveServiceListCall) SupportsAllDrives(supportsAllDrives bool) DriveServiceListCall {
	return c
}

func (c *MockDriveServiceListCall) IncludeItemsFromAllDrives(includeItemsFromAllDrives bool) DriveServiceListCall {
	return c
}

func (c *MockDriveServiceListCall) Do(opts ...googleapi.CallOption) (*drive.FileList, error) {
	files := []*drive.File{}
	for _, file := range c.service.files {
		files = append(files, file)
	}
	return &drive.FileList{Files: files}, nil
}
