package utils

import (
	"testing"
	"strings"
	"io/ioutil"
	"time"
//	"fmt"
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
	Info("this should be info: %d", 1000)
	Debug("this should be %s", "debug")
	Warn("this should be warn: %f %f", 10.0, 10.0)
	Error("this should be error")
	s := StartSpinner(Spinners[1], 10)
	for i := 1; i < 10; i++ {
		time.Sleep(300 * time.Millisecond)
		s.Update(10)
	}
	s.Stop()

	g := StartGauge(100, 10)
	for i := 1; i < 10; i++ {
		time.Sleep(300 * time.Millisecond)
		g.Update(10)
	}
	g.Stop()
	
}

