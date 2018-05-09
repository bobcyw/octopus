package python

import (
	"context"
	"testing"

	"fmt"
	"io"
	"time"

	"strings"

	"os/exec"

	"bufio"

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
	outReader, outWriter := io.Pipe()

	theOutReader := bufio.NewReader(outReader)
	readDone := make(chan int)
	go func() {
		var err error = nil
		var data string
		for ; ; data, err = theOutReader.ReadString('\n') {
			if err != nil {
				break
			}
			fmt.Printf("get: %s", data)
		}
		fmt.Println("err is ", err)
		close(readDone)
	}()
	wait, err := Run(ctx, script, reader, outWriter)
	<-wait
	fmt.Println("wait complete")
	<-readDone
	assert.Equal(t, err, nil, "应该为nil")
}

type Output struct {
}

func (inst *Output) Write(p []byte) (int, error) {
	fmt.Printf("->%s<-\n", string(p))
	return len(p), nil
}

func TestFileRunner(t *testing.T) {
	ctx := context.Background()
	script := "/Users/caoyawen/GOPATH/src/github.com/bobcyw/octopus/runner/python/testPython/hello.py"
	cmd := exec.CommandContext(ctx, "python3", script)
	writer, err := cmd.StdinPipe()
	reader, err := cmd.StdoutPipe()

	//outputReader := bufio.NewReader(reader)
	fmt.Println("start")
	go func() {
		//var err error = nil
		//var data byte
		buff := make([]byte, 1024)
		for {
			_, err := reader.Read(buff)
			if err != nil {
				fmt.Println("get error ", err)
				break
			}
		}

		//for ; err == nil; data, err = outputReader.ReadByte() {
		//	fmt.Println("get: ", data, err)
		//}
		//fmt.Println("err is ", err)
	}()
	fmt.Println("Step1")

	//out := Output{}
	//cmd.Stdout = &out
	//cmd.Stderr = &out
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
	fmt.Println("Step2")

	go func() {
		for i := 0; i < 4; i++ {
			fmt.Println("write ", i)
			io.WriteString(writer, fmt.Sprintf("[%d]\n", i))
			time.Sleep(1 * time.Second)
			//writer.Write([]byte(fmt.Sprintf("[%d]", i)))
		}
		//io.WriteString(writer, "EOF\n")
		//writer.Write([]byte("EOF"))
		writer.Close()
	}()
	fmt.Println("waiting...")
	<-done
}

func TestPipeline(t *testing.T) {
	reader, writer := io.Pipe()
	done := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			io.Copy(writer, strings.NewReader(fmt.Sprintf("Hello %d\n", i)))
			time.Sleep(time.Second)
		}
		reader.Close()
		close(done)
	}()
	go func() {
		theReader := bufio.NewReader(reader)
		for {
			if text, err := theReader.ReadString('\n'); err != nil {
				fmt.Println("got error ", err)
				break
			} else {
				fmt.Println(text)
			}
			//data, _ := ioutil.ReadAll(reader)
			//fmt.Println("receive:", string(data))
		}

	}()
	<-done
	fmt.Println("all complete")
}
