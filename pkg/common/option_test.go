package common_test

import (
	"github.com/jenkins-zh/docker-zh/pkg/common"
	"io/ioutil"
	"os"
	"testing"
)

const (
	fakeLTS = "fake-lts"
	fakeWeekly = "fake-weekly"
)

func TestReadCustomConfig(t *testing.T) {
	file, err := ioutil.TempFile(".", "yaml")
	if err != nil {
		t.Fatalf("cannot create temp file %v", err)
	}
	defer os.RemoveAll(file.Name())

	_, err = file.WriteString(getSampleCustomConfigFile())
	if err != nil {
		t.Fatalf("cannot write temp file %v", err)
	}

	mgr := &common.CustomConfigManager{}
	if err = mgr.Read(file.Name()); err != nil {
		t.Fatalf("cannot load the custom config file %v", err)
	}

	var ok bool
	// lts testing
	if ok, err = mgr.HasTLS(fakeLTS); err != nil || ok {
 		t.Fatalf("%s should not exists, %v", fakeLTS, err)
	}

	if err = mgr.AddTLS(fakeLTS); err != nil {
		t.Fatalf("failed when add lts %s, %v", fakeLTS, err)
	}

	if ok, err = mgr.HasTLS(fakeLTS); err != nil || !ok {
		t.Fatalf("lts %s should exists, %v", fakeLTS, err)
	}

	// weekly testing
	if ok, err = mgr.HasWeekly(fakeWeekly); err != nil || ok {
		t.Fatalf("weeekly %s should not exists, %v", fakeWeekly, err)
	}

	if err = mgr.AddWeekly(fakeLTS); err != nil {
		t.Fatalf("failed when add weekly %s, %v", fakeLTS, err)
	}

	if ok, err = mgr.HasWeekly(fakeLTS); err != nil || !ok {
		t.Fatalf("weekly %s should exists, %v", fakeLTS, err)
	}

	// save config
	if err = mgr.Save(); err != nil {
		t.Fatalf("cannot save config file %v", err)
	}
	// make sure we really saved it
	if err = mgr.Read(file.Name()); err != nil {
		t.Fatalf("cannot load the custom config file %v", err)
	}
	if ok, err = mgr.HasTLS(fakeLTS); err != nil || !ok {
		t.Fatalf("lts %s should exists, %v", fakeLTS, err)
	}
}

func getSampleCustomConfigFile() string {
	return `formulas:
- name: pipeline
  md5: c7689e6aa9c759f18eeb18ad8706b57b
lts:
- 2.204.5
weekly:
- 2.223
`
}