package python

import (
	"context"
	"io"
	"os/exec"
)

//Run 运行python
func Run(ctx context.Context, scriptName string, in io.Reader, out io.Writer) (<-chan int, error) {
	done := make(chan int)
	cmd := exec.CommandContext(ctx, "python3", scriptName)
	cmd.Stdin = in
	cmd.Stdout = out
	go func() {
		cmd.Start()
		cmd.Wait()
		close(done)
	}()
	return done, nil
}
