package python

import (
	"context"
	"os"
	"testing"

	"fmt"
	"io"
	"time"

	"io/ioutil"
	"strings"

	"os/exec"

	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	ctx := context.Background()
	script := "/Users/caoyawen/GOPATH/src/github.com/bobcyw/octopus/runner/python/testPython/hello.py"
	reader, writer := io.Pipe()
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			fmt.Println("write", i)
			writer.Write([]byte(fmt.Sprintf("sleep %d\n", i)))
		}
		//writer.Write([]byte("EOF"))
		writer.Close()
	}()
	wait, err := Run(ctx, script, reader, os.Stdout)
	<-wait
	assert.Equal(t, err, nil, "应该为nil")
}

func TestFileRunner(t *testing.T) {
	ctx := context.Background()
	script := "/Users/caoyawen/GOPATH/src/github.com/bobcyw/octopus/runner/python/testPython/hello.py"
	cmd := exec.CommandContext(ctx, "python3", script)
	writer, err := cmd.StdinPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err != nil {
		t.Fatal(err.Error())
	}
	done := make(chan int)
	go func() {
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
		}
		if err = cmd.Wait(); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	go func() {
		for i := 0; i < 4; i++ {
			io.WriteString(writer, fmt.Sprintf("[%d]\n", i))
			time.Sleep(1 * time.Second)
			//writer.Write([]byte(fmt.Sprintf("[%d]", i)))
		}
		//io.WriteString(writer, "EOF\n")
		//writer.Write([]byte("EOF"))
		writer.Close()
	}()

	<-done
}

func TestPipeline(t *testing.T) {
	reader, writer := io.Pipe()
	done := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			io.Copy(writer, strings.NewReader(fmt.Sprintf("Hello %d", i)))
			writer.Close()
			time.Sleep(time.Second)
		}
		close(done)
	}()
	go func() {
		for {
			data, _ := ioutil.ReadAll(reader)
			fmt.Println("receive:", string(data))
		}
	}()
	<-done
	fmt.Println("all complete")
}
