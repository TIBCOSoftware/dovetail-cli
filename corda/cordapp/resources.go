// Code generated by go-bindata.
// sources:
// resources/abstractflow.template
// resources/app.template
// resources/defaultschedulable.template
// resources/flow.template
// resources/kotlin.pom.xml
// resources/schedulable.template
// DO NOT EDIT!

package cordapp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _resourcesAbstractflowTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\xc1\x8a\xdb\x40\x0c\x86\xef\x7e\x8a\x39\xe4\xb0\xbb\x2c\x7a\x80\x85\x42\x97\x96\x42\x2f\x65\x21\x79\x01\x79\x2c\x1b\x35\x13\xcd\xa0\x91\x93\x06\x33\xef\x5e\xc6\x4e\x70\x68\x9a\xc3\x5e\x7c\x90\x7e\x7d\xff\xef\x7f\x12\xfa\x3d\x0e\xe4\xa6\x09\x7e\x6d\x4b\x69\x1a\x3e\xa4\xa8\xe6\x7c\x84\x84\x8a\x21\x50\x18\x85\x8f\xa4\x99\xa0\xe7\x96\x34\xc3\x76\xcc\x89\xa4\xc3\x36\xd0\x55\x2d\x64\xe0\xa3\x76\x58\xbf\x04\x7d\x88\xa7\x0c\x2f\x0f\xb6\xa3\x71\x60\x63\xca\xf0\xa1\x71\x50\xca\x79\xa7\xe8\xf7\xa4\xab\xf7\x01\x8c\x5b\x1f\xa1\x8b\x47\x32\xe4\x00\x3e\x8a\x21\x0b\xe9\x02\x4a\x09\xde\x53\xfa\x11\xe2\xe9\x53\x37\x4b\xae\xef\xd4\xe3\x18\x6c\xcb\x83\xec\x14\x25\xa3\x37\x8e\x72\x0b\xfb\x27\x70\x05\x29\x7a\xcb\xf0\x7e\x88\xa3\xd8\x03\x19\x77\x24\xc6\x76\x86\x0f\x54\x3b\x3f\x10\xd9\xea\x98\xa1\x46\xa0\xee\x26\xc4\xf5\xe8\x37\x1e\x71\xee\x09\x5e\x9a\x69\xda\x60\x4a\xee\xed\x8b\xab\xff\x5c\x4a\x1d\x7c\x8b\xd2\x2f\x76\x18\x6a\xf0\x3c\xaf\xef\xa6\xb3\x58\x51\x06\x72\x1b\x16\x36\x46\x8b\xfa\xea\x36\xbe\x0a\x87\xf9\xe6\xe7\x75\xbc\x60\x16\x3c\xe7\x5b\x54\xd5\xb1\x74\xf4\xc7\xfd\xc7\x77\xe5\x96\xd2\x7c\xbd\xd0\x58\x86\xb9\x4e\x6c\xf3\x5c\x9c\xf3\x01\x73\x76\x95\xbc\xaa\x9f\x9e\xdf\xdc\xe5\x11\x9f\x4c\x47\x7a\xbd\x33\x2e\xe5\xd9\x4d\x4d\x4d\x44\xd2\x95\xd2\xfc\x0d\x00\x00\xff\xff\xe5\xfa\xbb\x89\xac\x02\x00\x00")

func resourcesAbstractflowTemplateBytes() ([]byte, error) {
	return bindataRead(
		_resourcesAbstractflowTemplate,
		"resources/abstractflow.template",
	)
}

func resourcesAbstractflowTemplate() (*asset, error) {
	bytes, err := resourcesAbstractflowTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/abstractflow.template", size: 684, mode: os.FileMode(420), modTime: time.Unix(1558641594, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesAppTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x8e\xb1\x4e\xc3\x40\x0c\x86\xf7\x7b\x0a\xb7\x53\xb2\xf8\x01\x22\x21\x14\x31\x75\x61\xa0\x6c\x88\xc1\x3d\xcc\xe9\x4a\xee\xce\xf2\x39\x5d\xa2\x7b\x77\x94\xd0\x8a\x8e\x08\x4f\xb6\xf4\x7f\xfe\x3f\x21\xff\x45\x81\x61\x59\xf0\xf9\xd8\x9a\x73\x31\x49\x51\x03\x5f\x12\x5a\x3c\xf9\x82\x1f\xe5\xc2\x46\x71\x42\x5f\x94\x51\xe7\x6c\x31\x31\xfa\x92\x24\x4e\xac\x15\x47\x91\xa7\xeb\xf1\x2f\xf8\xaf\x90\x69\x0c\x81\x15\x0f\xaf\x3f\x8b\x73\xe5\x74\x66\x6f\xab\xfa\x28\xd2\xda\x21\xc9\x04\x8b\x73\x00\x00\x17\x9a\x80\x44\xe0\x01\xee\xf4\x6e\xbd\xa3\x48\x77\x0f\x0d\x83\x9f\xa8\x56\x3c\xd3\x85\x30\xb0\xbd\x70\x2d\xb3\x7a\x1e\xeb\xd1\x94\x29\x75\x7b\x12\xc1\x73\x2d\x79\xdf\xf7\xdb\xfb\xcf\x39\x43\x60\xbb\x9a\x74\x99\x12\x0f\x70\x34\x8d\x39\xf4\x03\xdc\x0c\x1f\x61\xd9\xd2\xeb\x28\xdb\xac\x79\x75\xc2\x5f\xb0\x76\xfd\x6e\xf7\xb6\xd2\xef\x5b\xb0\xb9\xf6\x1d\x00\x00\xff\xff\x0e\x56\xb0\x28\x90\x01\x00\x00")

func resourcesAppTemplateBytes() ([]byte, error) {
	return bindataRead(
		_resourcesAppTemplate,
		"resources/app.template",
	)
}

func resourcesAppTemplate() (*asset, error) {
	bytes, err := resourcesAppTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/app.template", size: 400, mode: os.FileMode(420), modTime: time.Unix(1564772381, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesDefaultschedulableTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8e\x31\x6a\x04\x31\x0c\x45\x7b\x9d\x42\xe5\x4c\x0a\x1f\x60\x09\x61\xab\x40\x20\x6c\x91\x21\x07\xd0\x6a\xe4\x8d\x89\xc6\x36\xb2\x67\xb6\x30\xbe\x7b\x70\x60\x21\x4d\x1a\x81\xf4\x9f\x3e\x0f\x32\xf1\x37\xdd\x04\x5b\x73\x97\xa5\x77\x80\xb0\xe5\x64\x15\x39\xb9\x4c\x46\xaa\xa2\x7b\x0c\x87\x58\x11\xe7\xc3\x55\xac\xb8\x65\x2f\x59\xe2\x4a\x57\x95\x07\x1d\xa5\x3a\x4e\xb6\xd2\x98\xe2\xbc\xa6\x7b\x71\x4f\xff\xa4\x9c\x62\x35\xe2\x5a\xdc\x52\xa9\xca\x87\x78\x80\xf3\xc2\x5f\xb2\xee\x3a\x4a\x5f\x35\xdd\x81\x95\x4a\x19\x56\x63\xbb\xd0\x26\xbd\xbf\x6d\x59\xa7\x83\x14\x4d\xfc\x09\x1f\xbf\xf3\x09\x07\xf2\x9e\x6e\x81\x9f\x3f\x63\xa8\x2f\x38\xcd\x0d\x10\x11\xcf\x7f\x4d\xc7\x21\x1d\x62\x16\x56\x41\xbf\x47\x64\x52\x9d\x66\x6c\xf0\x1b\x75\x80\xfe\x13\x00\x00\xff\xff\x22\xf6\x7d\x5a\x0d\x01\x00\x00")

func resourcesDefaultschedulableTemplateBytes() ([]byte, error) {
	return bindataRead(
		_resourcesDefaultschedulableTemplate,
		"resources/defaultschedulable.template",
	)
}

func resourcesDefaultschedulableTemplate() (*asset, error) {
	bytes, err := resourcesDefaultschedulableTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/defaultschedulable.template", size: 269, mode: os.FileMode(420), modTime: time.Unix(1563827395, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesFlowTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\x41\x6f\xe3\x36\x13\xbd\xfb\x57\x0c\x8c\x1c\xa4\x85\x40\x7c\x67\x7f\x4d\xb1\xde\x14\x8b\x0d\x90\x4d\x0c\xdb\x45\x8f\x8b\x89\x34\xd6\xb2\xa1\x29\x75\x48\x39\x31\x0c\xfe\xf7\x62\x64\xc9\x92\x65\x7b\x37\x6e\x9b\x43\x02\x52\x8f\x6f\x86\x33\x6f\x5e\x58\x62\xfa\x82\x39\xc1\x6e\xa7\x1e\x17\x21\x8c\x46\x7a\x5d\x16\xec\x21\x2d\x54\x89\x8c\xc6\x90\xa9\xac\xde\x10\x3b\x52\x2b\xfd\x4c\xec\xd4\xa2\x72\x25\xd9\x0c\x9f\x0d\xb5\x68\x4b\x5e\xa5\x05\x67\x28\xbf\x49\xad\x4c\xf1\xea\xd4\x87\x0b\x5f\x2b\xaf\x8d\xf6\x9a\x9c\x9a\x71\x91\x33\x39\xb7\x64\x4c\x5f\x88\xbb\xd8\x6b\xe5\xf5\x73\x5a\xa8\xac\xd8\x90\x47\x6d\x54\x5a\x58\x8f\xda\x12\xef\x89\xca\x52\x4d\xcb\xf2\xb3\x29\x5e\xaf\x3a\x33\x88\xb7\xf0\x54\xba\xab\x08\xf6\x17\xfb\x8d\x56\x58\x19\xbf\xd0\xb9\x5d\x32\x5a\x87\xa9\xd7\x85\xed\x67\x33\xb8\xb1\x10\x31\xa6\xde\xa9\xe9\xba\xa8\xac\xbf\x00\xd3\x19\x59\xaf\xfd\x56\xcd\x90\xfd\xf6\x02\xc8\x77\x11\x9d\x92\x14\x28\xeb\x25\xd1\x1e\xfa\x13\x37\x58\x17\x5a\x7d\x18\xed\x76\x37\x58\x96\x30\xb9\x05\x29\x5a\x08\xb2\x71\x57\xd8\xd5\x3e\x1c\x1a\x49\xdc\xd5\x9f\x4f\x76\x6b\x30\xa3\xcd\x09\x6e\xb4\xd5\x5e\xa3\x2f\x38\x81\x9b\x54\x80\x79\x7d\xe6\xbe\xdd\xde\xd3\xec\xe9\xb5\xeb\x53\x09\x4e\xdb\x8c\xde\xe0\x4c\xdc\x8e\x37\x84\xd1\xc7\x85\x47\xf6\xa2\xac\x4f\xdb\xf9\xec\x6e\x94\x1a\x74\x0e\x84\xb0\x03\xdd\xaf\x4b\x13\x75\x59\x25\x70\x83\xde\xb3\xc4\x68\xd2\x52\x53\xef\xd9\x85\xb0\xdb\xe9\x15\x58\xc1\xc0\xff\x64\x35\x4e\xc6\xf2\x87\x6c\x16\xc2\x06\x8d\xd0\xca\x49\xf5\x88\x6b\x0a\x61\x02\xbb\xdd\x9e\x4a\x2d\xb7\x25\x85\x00\x0d\x34\x9e\x0c\x12\x88\x62\xd8\x8d\x46\x00\x00\xc5\x86\x98\x75\x46\x20\x74\xe5\xb1\xb8\xe0\x16\x1a\x89\xaa\x9c\xfc\xa1\x4a\x03\x09\x46\x27\xc5\x0a\x21\x91\x78\xcd\x55\xbe\xa0\x7b\x7a\x76\xc4\x32\x80\x21\xc4\xfb\xb0\x1f\xfb\x13\x78\x94\xc7\xaa\xb2\x90\xa2\x31\x51\x0c\x13\x38\x11\x07\xec\x6a\xb4\xfc\x6c\x90\x01\x39\x77\x70\x0b\x0f\xda\xbe\x50\xf6\x05\xdd\xf7\xaf\x58\xfe\xb2\xf0\xac\x6d\x9e\xc0\xd4\x6e\x7f\x8d\xe2\x23\x7c\x89\xec\x75\xaa\x4b\xb4\x5e\xce\x4d\x99\x71\xfb\xa0\x9d\x6f\xce\xf4\xe1\xef\x69\xcf\x01\x2c\x79\xa8\xb2\xf2\xd1\xf8\xb8\x23\xe3\x64\xd0\xa2\x3e\xbf\x5e\x01\xfd\xd5\xf4\xab\x1e\x17\x69\x1a\x8c\x1f\x0b\x8f\xbc\x1d\x87\x00\x98\x65\xfb\x45\x74\x99\x84\x8c\x23\xb8\xc0\xd4\x96\xbd\xe1\x6a\x97\xff\x90\x6d\xd6\xd5\x4e\x08\xfb\xa5\x54\x98\x65\x27\x57\x8f\x1b\xf1\x41\x9f\x5f\x36\x0e\x6b\x47\xfe\xa9\xe2\xfb\xc6\x32\xa2\x33\xb5\x19\x0c\xa1\xe7\x8a\x42\x70\xaf\x58\x36\x87\x34\xb9\x83\x2c\x23\x69\x42\x72\x94\x57\x32\x54\x74\x3c\x4c\xe1\x10\xe9\x8c\x58\x9b\x70\x75\x9a\xcd\x5e\x9d\x4d\x1e\x75\xe2\x5e\x90\xcd\x96\x6f\x76\x59\xb4\x88\xaf\x68\x2b\x34\x66\x1b\xc2\x49\x2c\xae\x6a\x83\xad\x0b\xd5\x9b\xc5\x46\x23\x62\x6b\xe2\x0b\x32\x6a\x4b\xd6\x79\x4e\x7c\x82\x8c\x93\x5a\x69\x5d\xa1\x98\x7c\xc5\x16\x0e\x18\x99\x97\xa9\xcd\xee\x8a\xf5\x5a\xfb\x68\x78\xf9\xfa\x58\x18\x89\xbd\xed\x33\x3b\x68\x9c\xc9\x95\x85\xcd\x48\x8c\xf1\x40\x56\x7b\xe3\xbc\xfd\xd2\x99\xe9\x55\xd6\x78\x60\x16\x6b\x6c\x7a\x45\xd9\xa7\x5a\xd2\xbd\xab\x4d\x26\xb5\x4f\xc6\x9d\x5d\xf6\x0e\xd6\x76\x29\xfe\x94\xca\x7f\x1e\x62\x69\xf1\x76\x41\xce\xe9\xc2\x4e\x40\x02\x35\x0b\xb1\x8c\xc6\xb1\xa2\x15\x1a\x47\xc9\x19\x6f\x8a\x1b\x0f\x79\xbf\xf3\x1d\x6a\xf0\x0e\xe7\xbb\xc2\xdf\x7e\x6c\x6f\x3f\x9e\x81\x6e\x86\x8e\x86\x61\x4e\x29\xc9\x4b\x27\x3a\x53\xa9\xd3\x61\xb8\x38\x98\xfe\xbb\x76\xaa\x2f\xd7\x5e\x33\x7e\x22\xd7\x3e\x32\x4e\xc0\x56\xc6\xf4\x2d\xd8\xc0\x86\x58\xaf\xb6\xc2\x0c\xb7\x70\xf9\x19\x72\xee\x06\x27\xba\xe7\xe6\xba\xc7\xb2\xef\x22\x24\xe7\x14\x73\xa1\x0e\x67\xe7\xa2\x68\xa6\xfa\x64\x2c\xda\x71\x6f\xa7\xe2\x2a\x69\xb7\xac\xff\x4a\xd9\xa2\x83\xeb\xa5\xdc\xe6\x3d\x54\xf2\x7f\x26\xdb\x41\x74\x95\x56\xcc\x64\xbd\xbc\x53\xe1\x16\xce\xbd\x5e\xd5\xfc\xf7\xc7\x6f\xd3\xd9\xec\xdb\xe7\x87\xa7\x3f\xfe\x7f\x59\x83\x5d\xd5\x7e\x22\xc1\x1e\xf0\x44\x81\x8d\x70\x5a\xc8\x9c\xe4\x59\xba\x7c\xb3\x57\x0c\x4c\x18\xb5\x2f\xb0\xd1\xdf\x01\x00\x00\xff\xff\xf8\x79\x12\x10\x7a\x0c\x00\x00")

func resourcesFlowTemplateBytes() ([]byte, error) {
	return bindataRead(
		_resourcesFlowTemplate,
		"resources/flow.template",
	)
}

func resourcesFlowTemplate() (*asset, error) {
	bytes, err := resourcesFlowTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/flow.template", size: 3194, mode: os.FileMode(420), modTime: time.Unix(1558641594, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesKotlinPomXml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x57\x4d\x6f\x1b\x37\x13\x3e\xcb\xbf\x82\xaf\x61\x1f\xde\x20\x4b\x4a\x4e\x82\xba\x06\xbb\x97\xa6\x07\x17\x49\x53\xc4\x41\x50\xa0\x27\x8a\x3b\x5a\x51\xe6\x92\x04\xc9\x95\x54\x18\xfe\xef\xc5\x92\x5c\xed\x97\x24\xf7\x2b\xba\x2c\x39\xf3\x0c\xf9\x0c\x67\x86\x43\x51\x63\xf5\x06\xb8\x47\xfb\x4a\x2a\xf7\xc3\xe5\xda\x7b\x73\x47\x48\xc5\xb6\xa0\x30\x33\x8c\xaf\x01\x6b\x5b\x92\x5f\x3f\x7d\x24\x6f\xf1\x1c\xcf\x2f\x23\xf2\x6e\xef\xc4\x01\xbd\xdb\xed\xf0\xee\x4d\xc0\xdd\xcc\xe7\x0b\xf2\xdb\xc7\x0f\x0f\x7c\x0d\x15\xcb\x84\x72\x9e\x29\x0e\x97\x68\xef\xc4\x9d\x0b\xc2\x0f\x9a\x33\x2f\xb4\xfa\x0b\x9b\xa1\x53\x88\xbd\x2b\xa2\x30\x0b\x38\xbc\x77\xc5\x65\x7e\x81\x10\xad\x74\x01\xf2\x2b\x58\x27\xb4\xca\x83\x8e\x92\x81\xac\x01\x19\xab\x0d\x58\x2f\xc0\x35\x53\x84\xe8\x92\x39\x78\x2f\x6c\x4e\x3c\x54\x86\x92\x76\x7a\x31\x9b\xd1\xd2\xea\xda\xdc\x17\x79\xfa\x52\xd2\x0a\x2e\x10\x9a\x51\x66\xbd\x58\x31\xee\xef\x8b\xbc\x1b\x52\xd2\x13\x07\xd8\x36\x6d\x9e\xbe\x94\x6c\x3b\x36\x0d\xe9\xe0\x1e\xd7\x95\x11\x12\x2c\xf6\xcc\x96\xe0\xf3\x05\xbe\xa5\xe4\xb8\xea\xa8\x99\xd3\xb5\xe5\x70\xd4\x2c\xa9\xa2\xd9\xa3\xf6\x52\x28\xdc\x52\x58\xe0\x1b\xfc\xdd\x82\x92\x91\xb8\xf1\x3d\x89\x0e\xcb\x48\xa6\xca\x9a\x95\xf0\xb5\x33\x3d\xd8\x9d\x04\x1d\x5b\x88\x19\x71\x76\x8d\x9e\xfe\x98\xf9\x66\x5b\x7d\xe9\x9d\xd1\x69\xf5\x05\x25\xc3\x58\x1f\xc2\x79\xf5\x94\x46\xcf\x83\x88\xf6\x03\x7a\xf5\xd4\x4d\x9e\xc7\x31\x3d\x84\xf4\xea\x29\x8d\x9e\x07\x51\xa5\x05\x18\x50\x05\x28\x1e\x76\x9e\xcd\x3a\xc1\x1f\x31\x0e\x21\x16\xed\xd6\xda\x96\x78\x03\x7e\x69\x99\x50\x0e\x47\x87\x06\xc4\x12\xbe\x47\x22\x82\x32\xe7\x0b\x29\x96\xd9\xa6\x78\xbc\x1d\x73\x4c\x36\x1d\xd3\x61\x88\x9f\xc7\x69\x48\x86\x14\xc7\x94\x7b\xc5\xa0\xc0\x63\xae\x6d\xc1\x86\xe5\x30\xa8\x87\xa0\xcf\xb8\xb6\x30\xa9\x87\xae\x20\x7e\x7f\x8b\xe7\xaf\xff\x3f\x20\x32\xfb\x26\x3c\x56\x42\x35\xf7\xd0\x39\x2a\x37\xcd\x6d\xf1\x2f\xc8\x70\x5d\x61\x2f\x96\x5c\xe3\x42\x6f\xc1\x33\x21\xcf\xb0\x6a\x21\x59\xa2\x7f\x9a\xd6\x1c\x2f\xf0\xdb\xd7\xe8\x05\x5e\xfd\x8c\x3b\x46\xcc\x30\xa9\x2b\xa1\xb4\x64\x4b\x87\x2b\xf0\x56\x70\x77\x86\x5e\x42\x64\x0a\x76\x99\x05\x29\xf8\x39\x86\x0b\x3c\xc7\xef\x5e\x3a\x38\x84\xae\xaf\x61\xef\xc1\x2a\x26\xaf\xaf\x2f\x06\xe9\xd6\x16\xa7\x05\xa3\x9d\xf0\xda\x76\x37\xf3\x41\x74\xa8\x1a\xea\x14\x33\x6e\xad\xbd\xeb\xe5\x38\x28\xb6\x94\x50\xe4\x2b\x26\x1d\x50\xd2\x4e\x5b\x13\x32\xb1\xa1\xa2\xc8\x39\x28\x6f\x99\xa4\x44\x74\x48\xc5\x2a\xc8\x7f\x8c\x0a\xf4\xf9\xb0\x39\x25\x41\xd1\xa2\x6a\x2b\xf3\xa6\x35\xb9\x3b\x42\x1a\x86\x78\xd2\xa0\x82\xe0\x86\x92\x06\x99\xaa\x6b\xe2\xca\x37\x75\xce\x19\x2b\x54\x39\xf5\x2d\xca\x91\x91\x75\x29\xd4\x79\x0f\xdb\xee\x1b\x3c\x8c\x76\x58\x68\x12\x4d\x5d\x93\x17\xc0\x1c\x90\x33\x4e\xf6\xe6\x6d\x8c\x97\xb5\x90\xa9\x2b\x16\xc2\x02\x0f\xd0\xab\xa7\xd4\x72\x9f\x29\xe9\xa4\x71\xc9\xd8\xbb\xde\x1f\xc1\x12\x67\x39\xa9\x98\x50\xa4\xbd\x30\xc7\xd8\xb8\x82\xae\xbd\xa9\xfd\xd1\x15\x62\x3f\x4d\xf6\x84\x4b\xe6\x1c\x38\x4a\xc6\x16\x6d\x2e\xc6\xf5\xbb\x73\x6e\x25\xbd\x60\x1d\x73\xaa\x23\x7a\x58\x62\xe2\x67\x3c\xac\x41\xa3\x26\xa3\x0d\x69\x73\x89\xc9\x5f\x9a\x30\xa5\xcb\xff\xea\x29\x3d\xdf\x70\xaf\x55\xf5\xa4\xdd\x4d\xdf\x99\xc6\xa5\x52\x10\x07\xb3\x13\x7d\x29\x65\x75\x4c\xf1\x64\xf7\x42\x6f\x8a\x4f\xb3\x0d\xb3\x59\xc4\xbf\xd0\x99\xde\xe0\xf0\x0a\x18\x74\xa3\x00\xe0\x5a\xad\x44\x59\xdb\xf0\x5e\xec\xe4\x2f\x44\xf5\x44\xfc\xe2\x99\x9e\x58\x92\xc2\x1e\x78\xdd\xc8\x5c\x7f\x9f\x4e\x3c\xda\x5d\x14\xb9\xa8\xcc\xe0\xee\x68\x8f\x4e\x33\xe9\xf2\xf0\xc9\x37\xcc\x52\x12\x46\xf1\xe3\x46\x58\xb3\x66\x0e\x72\xc3\xf8\x23\x2b\x81\x92\x38\x1d\x42\xce\x1c\xc1\x7f\x9a\x11\x7d\x9f\x15\xd7\x85\x50\x65\x7e\xff\xf0\x29\xbb\xbd\x7d\xf7\x7d\xb6\x68\x6e\x9d\x24\x1c\xa3\x43\xd1\x88\x95\x00\x9b\x33\x23\x28\xe9\xcd\x11\xfd\x5f\x96\xa1\x2f\x6b\xe1\xd0\x4e\x48\x89\x96\x80\x98\x09\x77\x7e\x81\x76\xc2\xaf\x91\x5f\x03\xda\x30\x8b\xb8\x05\xe6\x5b\xe1\x81\x18\xf2\xac\x44\x5b\x26\x6b\x40\x62\x85\x8c\x05\x07\xca\x67\xd9\x94\xef\x9e\xcb\xba\x00\x37\x56\x74\xaa\xfc\xd5\x2b\xf2\xea\xbe\x32\x12\x07\x76\x94\xb4\xf2\xb3\x16\xcc\x18\xbc\x71\xcd\x73\xbd\x15\x8e\xd1\xf1\x37\xe6\x43\x0e\x84\xc6\xca\x93\xe9\x17\xad\x26\xa9\xd6\x13\xb6\xa5\x4a\xda\x5a\x9d\xcd\x26\x65\x8b\x8e\x3f\x11\x63\x35\x9e\xab\x44\xf4\x0f\xde\xa2\xe8\xef\xbc\x2d\x7b\xc7\x7b\xac\xcc\xc6\xca\xa9\x0e\xb5\x3d\x3b\x3e\xf2\xa7\x75\xd7\xb9\x11\xea\x2c\x7e\x3b\x78\x98\xa1\xa3\x55\x78\xf2\xf4\xa7\x11\x98\x85\xdf\xa9\x9a\x6c\x54\x4a\xef\x98\x55\xb9\xb7\x35\x50\x92\x26\x28\xd6\xc1\x7b\xe1\x9a\xc6\x8d\x1a\x99\x50\xa5\x43\xfd\x54\x9e\x85\xf7\x57\x39\xa2\x16\x85\x79\xb6\x61\x5b\x96\x19\x66\x59\x05\x1e\xac\x6b\x82\x58\xa6\xea\xfa\x29\xbc\x06\x90\x6b\xde\x6b\x1e\x35\x7f\x75\xd1\x4a\x5b\xf4\xf3\xc3\xe7\xec\xcd\xfc\x1d\x62\x4a\x69\x1f\x78\x4e\xf7\x23\x71\xc3\x59\xf2\xea\xe4\xd5\x78\x48\xb9\x76\xd4\xb0\xa4\x24\x75\xf2\xf0\x17\xab\xb9\x5b\xf2\x3f\x03\x00\x00\xff\xff\x06\x4a\x73\x9d\x4c\x10\x00\x00")

func resourcesKotlinPomXmlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesKotlinPomXml,
		"resources/kotlin.pom.xml",
	)
}

func resourcesKotlinPomXml() (*asset, error) {
	bytes, err := resourcesKotlinPomXmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/kotlin.pom.xml", size: 4172, mode: os.FileMode(420), modTime: time.Unix(1558641594, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesSchedulableTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x55\x4d\x8f\xdb\x36\x10\xbd\xeb\x57\xcc\x0a\x3d\x48\x89\x41\xf4\xbc\xa8\xdb\xb8\x8b\x04\x59\x34\xcd\x02\x95\x8b\x1e\x8b\x31\x35\x92\x19\x53\x24\x41\x8e\xbc\x35\x1c\xfd\xf7\x82\x92\x65\x5b\x5e\x6d\x17\xd5\xc1\x86\xa9\x79\x6f\xbe\x1e\x9f\x1d\xca\x1d\xd6\x04\xc7\xa3\xf8\x5a\x74\x5d\x92\xa8\xc6\x59\xcf\x20\xad\x70\xe8\x51\x6b\xd2\xad\x51\x7b\xf2\x81\x44\xa5\x36\xe4\x83\x28\xda\xe0\xc8\x94\xb8\xd1\x34\x46\x1b\x62\x21\xad\x2f\x31\x7e\x92\xa8\xb4\x7d\x0e\xe2\xdd\x85\xab\x11\xac\x36\xd2\x8a\xd2\xee\x89\x51\x69\x21\xad\x61\x54\x86\xfc\x00\x73\x4e\xac\x9c\xfb\xa4\xed\xf3\xff\xc5\x3c\x8c\x87\xaf\x94\x12\x41\x1e\x25\x07\x51\x30\x32\xfd\x41\xd5\x2b\x81\xec\xd1\x04\x94\xac\xac\x09\xa2\x50\xb5\xa1\x72\x7d\x39\x1a\x41\x3b\xcb\x5a\x19\x21\xad\xd6\x74\x8a\xfd\xa2\xcc\x8e\xca\xcf\x18\xb6\xbf\xa3\xbb\x89\xf3\x54\xc5\x38\x51\xb5\x5a\x0b\xe7\x55\x83\xfe\xf0\x60\x4d\x60\xdf\x4a\xb6\x3e\x49\x8e\xc7\x1f\xd0\x39\xb8\x5f\x42\x6c\xa6\xeb\x92\x0f\x85\xdc\x52\xd9\xea\x38\xdd\x7e\x1e\x52\x63\x08\x71\x3d\xf1\xd7\x57\x6c\xa8\xeb\x1e\x1b\xa7\xb3\x3d\x6a\xf0\x54\xdd\xc3\xd8\x57\x7e\x0f\xa7\x19\x66\x15\xea\x40\x0b\xe8\xbf\x72\x38\x26\x09\x00\xc0\x87\xeb\xbd\xc5\x03\xbb\x27\xef\x55\x49\x50\xb5\x06\x24\x6a\x9d\xe5\x70\x0f\x2f\x5a\xff\x05\x8e\x7d\x78\x7c\xb4\xad\x6b\xf2\x42\x99\xca\x66\xe9\xa4\x26\x50\x01\xd8\xab\xf8\x9a\x4a\xd8\x1c\x00\x21\x5c\x3a\x01\xda\x93\xe1\x34\x3f\x13\xb1\x3f\x5c\x58\xe3\x13\x88\x9f\x5a\xff\x58\x92\x61\xc5\x87\x2c\x4f\x26\x6f\x63\xb3\x51\x54\x76\xf3\x0d\x96\x20\x3d\x21\x53\xd1\x6e\xe2\xd1\xd3\xe6\x1b\x49\xce\x3c\x55\xf9\x04\xa2\xaa\x6c\x44\xdc\x2d\xc1\xb4\x5a\xe7\xd3\x8c\xf1\xf1\xc4\xad\x37\x10\xda\xcd\x30\xb7\x01\x30\x25\xea\x80\x74\x20\x78\x09\x9e\x0c\xc3\x58\xde\x2a\x53\x03\x5b\x50\x66\x6f\x77\x94\xe6\xaf\x65\x8b\xb5\x4c\x33\x9c\x7f\x75\x12\x59\x6e\xb3\x8f\xff\x48\x72\x71\xf6\xd0\xaf\xef\x3a\x96\xb7\xde\x3e\x83\xa1\x67\x88\x15\x9f\x03\x33\xba\xa4\x1b\xe8\xba\x61\x82\xce\xab\x3d\xf2\x69\xc7\xf3\x73\xbb\x92\x10\x9c\x35\x74\xbd\xf5\x3d\x7a\x40\x5f\x07\x58\xc2\x44\xec\x3f\x15\xec\x95\xa9\x17\xb0\x32\x87\x9f\xb3\x4b\x01\x31\x56\xb8\x96\xb3\xf4\xea\x52\x3d\x1a\xd7\x72\xba\x80\x40\x7e\xaf\x24\x7d\x6e\x37\x42\x5b\x2c\xfb\xd4\xfd\xf6\x44\x89\x8c\x79\x72\xbd\xf2\x93\xa0\x60\x19\xf5\xff\x68\x14\x2b\x64\x65\xea\x5e\x75\x45\xd7\x89\xe1\xf6\x0c\x17\x42\xd4\xc4\xeb\x21\xfe\x46\x9a\xe9\x94\x54\xb2\x89\x8c\xbe\x35\xfd\xce\x6f\x62\x17\x63\xd2\x45\xdf\xc6\x4d\x3d\x18\x76\x71\x0a\x91\x22\xe6\x3b\x9b\xcf\x2a\x1c\x8c\x5c\xc7\xb7\xd9\xb5\x27\x89\xf5\xaa\xf8\xed\xef\xe2\xcf\x5f\x3f\x7d\x79\xfa\xeb\x42\xa5\xaa\xec\xc4\x34\x08\x13\xbe\x7f\x1f\xa8\x85\x0a\x1f\x1b\x17\xe5\x7f\xa3\xd5\x39\xe5\xcc\xea\x32\x56\x19\x86\x05\xc3\xf2\x44\x5a\x13\x67\x3f\xe6\x80\xe1\x3f\xb6\x37\x7b\xdb\x06\xeb\x59\x42\x3a\x3f\xfd\x14\xde\x8f\xb9\xfa\x1c\xe9\x38\xc6\x34\x87\xf7\x90\xc6\xa5\xa4\xb3\xbc\x27\x2d\x4d\xb0\x2b\x5f\xb7\x0d\x19\x0e\xe9\x5b\x95\x26\xaf\xde\xc1\xa9\x21\x9d\x6f\x62\x9f\x13\x62\xb5\xe7\xa6\xf2\x17\x75\xed\xc6\x66\x1f\xe2\xb7\xa8\xac\x8f\x34\xd9\x05\x21\x06\x4b\x7f\x01\x8c\x3e\x0e\xcb\x13\x7e\xce\xe3\x6f\x01\x38\x76\xda\x0b\xc9\xfa\xbb\xbb\xfe\x9f\xb6\x21\x8e\x7f\xae\x0d\xba\x19\xab\x19\xa7\xd6\xcf\x4a\xb1\x30\xd8\xd0\x8d\x4b\x09\xb6\xeb\x83\xa3\x72\xe5\x3d\xbe\xe5\x9f\x6c\xbd\xe8\x0d\xff\xdd\xb9\x98\x7c\x4e\x6f\x23\x04\xc3\x68\x0b\x57\xea\x3b\xb9\x4c\x97\xfc\x1b\x00\x00\xff\xff\x3e\x3e\xfc\x9b\x43\x08\x00\x00")

func resourcesSchedulableTemplateBytes() ([]byte, error) {
	return bindataRead(
		_resourcesSchedulableTemplate,
		"resources/schedulable.template",
	)
}

func resourcesSchedulableTemplate() (*asset, error) {
	bytes, err := resourcesSchedulableTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/schedulable.template", size: 2115, mode: os.FileMode(420), modTime: time.Unix(1564691621, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"resources/abstractflow.template": resourcesAbstractflowTemplate,
	"resources/app.template": resourcesAppTemplate,
	"resources/defaultschedulable.template": resourcesDefaultschedulableTemplate,
	"resources/flow.template": resourcesFlowTemplate,
	"resources/kotlin.pom.xml": resourcesKotlinPomXml,
	"resources/schedulable.template": resourcesSchedulableTemplate,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"resources": &bintree{nil, map[string]*bintree{
		"abstractflow.template": &bintree{resourcesAbstractflowTemplate, map[string]*bintree{}},
		"app.template": &bintree{resourcesAppTemplate, map[string]*bintree{}},
		"defaultschedulable.template": &bintree{resourcesDefaultschedulableTemplate, map[string]*bintree{}},
		"flow.template": &bintree{resourcesFlowTemplate, map[string]*bintree{}},
		"kotlin.pom.xml": &bintree{resourcesKotlinPomXml, map[string]*bintree{}},
		"schedulable.template": &bintree{resourcesSchedulableTemplate, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
