package client

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {

	wdir, _ := filepath.Abs("/Users/mwenyan@tibco.com/Downloads/tmp")

	opts := NewOptions("/Users/mwenyan@tibco.com/webinar/demo/cli/investor_cp.json", "1.1.0", wdir,
		"com.investor.cp.flows", "/Users/mwenyan@tibco.com/webinar/demo/cli/cp.json",
		"/Users/mwenyan@tibco.com/webinar/demo/cli/cp_dep.pom")
	generator := NewGenerator(opts)

	//config, err := model.DecodeApp(bytes.NewReader([]byte(jsonstring)))
	err := generator.Generate()
	if err != nil {
		fmt.Errorf("generate err %v", err)
	}
}
