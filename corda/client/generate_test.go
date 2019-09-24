package client

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {

	wdir, _ := filepath.Abs("/Users/mwenyan@tibco.com/Downloads/tmp")

	opts := NewOptions("/Users/mwenyan@tibco.com/dovetail/src/github.com/TIBCOSoftware/dovetail-cli/corda/client/iou.json", "1.1.0", wdir, "com.charlie.iou.flows")
	generator := NewGenerator(opts)

	//config, err := model.DecodeApp(bytes.NewReader([]byte(jsonstring)))
	err := generator.Generate()
	if err != nil {
		fmt.Errorf("generate err %v", err)
	}
}
