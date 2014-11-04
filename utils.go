package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func goGet(p string) error {
	cmd := "go"
	args := []string{"get", p}
	if err := runCmd(cmd, args...); err != nil {
		return err
	}
	return nil
}

func (e *env) goBuild() error {
	os.Chdir(e.BldDir)
	cmd := "go"
	args := []string{"build", tmpname}
	if err := runCmd(cmd, args...); err != nil {
		return err
	}
	return nil
}

func bldDir() string {
	f, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func cleanDirs(targetDir string) error {
	if err := os.RemoveAll(targetDir); err != nil {
		return err
	}
	return nil
}

func runCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	buf := make([]byte, 1024)
	var n int
	for {
		if n, err = stdout.Read(buf); err != nil {
			break
		}
		log.Print(string(buf[0:n]))
	}
	if err == io.EOF {
		err = nil
	}
	return nil
}

func compare(A, B []string) []string {
	m := make(map[string]int)
	for _, b := range B {
		m[b]++
	}
	var ret []string
	for _, a := range A {
		if m[a] > 0 {
			m[a]--
			continue
		}
		ret = append(ret, a)
	}
	return ret
}

func (e *env) logger(facility, msg string, err error) {
	if e.Debug {
		if err == nil {
			log.Printf("[info] %s: %s\n", facility, msg)
		} else {
			log.Fatalf("[error] %s: %s %v\n", facility, msg, err)
		}
	}
}