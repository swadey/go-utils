package utils

import (
	"testing"
	"strings"
	"io/ioutil"
)

func TestZopen(t *testing.T) {
	fns := [...]string{ "test.txt", "test.gz", "test.bz2", "test.xz" }
	for _, fn := range fns {
		r, in := Zopen("test-data/" + fn)
		b, _ := ioutil.ReadAll(r)
		t.Log("testing " + fn)
		if strings.TrimSpace(string(b)) != "hello world" {
			t.Errorf("Error reading from %s\n", fn)
		}
		in.Close()
	}
	
	Zopen("test.bz2")
	Zopen("test.txt")
	Zopen("test.txt.xz")
}

func TestLogger(t *testing.T) {
	log := Logger()
	log.Info("this should be info: %d", 1000)
	log.Debug("this should be %s", "debug")
	log.Warn("this should be warn: %f %f", 10.0, 10.0)
	log.Error("this should be error")
}

