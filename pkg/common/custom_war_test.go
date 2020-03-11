package common_test

import (
	"fmt"
	"github.com/jenkins-zh/docker-zh/pkg/common"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	data := map[string]string{
		"version": "hello",
	}
	result, err := common.RenderTemplate("/Users/rick/Workspace/GitHub/jenkins-zh/docker-zh/formulas/zh.yaml", data)
	fmt.Println(result, err)
}
