package common_test

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	data := map[string]string{
		"version": "hello",
	}
	result, err := common.RenderTemplate("/Users/rick/Workspace/GitHub/jenkins-zh/jenkins-formulas/formulas/zh.yaml", data)
	fmt.Println(result, err)
}
