package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Command struct {
	name string
	args []string
	ctx  *context.Context
	exec *exec.Cmd
}

func NewMakeCMD(name string, ctx *context.Context, args ...string) *Command {
	command := exec.Command("make", name)
	command.Args = append(command.Args, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	return &Command{name, args, ctx, command}
}

func (cmd *Command) exc(ctx context.Context) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	cancel()
	go func(ctx context.Context, n int) {
		cmd.exec.Run()
		for i := 0; i <= n; i++ {
			select {
			case t := <-time.Tick(time.Second):
				log.Println(t)
				time.Sleep(time.Second)
			case <-ctx.Done():
				cmd.Stop()
				return
			}
		}
	}(ctx, 10)

}
func (cmd *Command) Execute(c *context.Context) error {
	ctx, cancel := context.WithDeadline(*c, time.Now().Add(10*time.Second))
	defer cancel()
	go cmd.exc(ctx)
	return ctx.Err()
}

func (cmd *Command) Stop() {
	cmd.exec.Process.Release()
	// cmd.exc.Process.Kill()
	// cmd.exc.Process.Signal(os.Interrupt)
	// cmd.exc.Process.Kill().Error()
	// cmd.exc.Process.Signal(os.Kill)
	cmd.exec.Process.Kill()

}

type CLI struct {
	ctx      context.Context
	Cmd      chan string
	Err      chan error
	Out      chan []byte
	env      map[string]string
	commands map[string]*Command
	scan     bufio.Scanner
}

func NewCLI(commands map[string]*Command) *CLI {
	cli := &CLI{
		ctx:      context.Background(),
		Cmd:      make(chan string),
		Err:      make(chan error),
		Out:      make(chan []byte),
		env:      map[string]string{},
		commands: commands,
		scan:     *bufio.NewScanner(os.Stdin),
	}

	return cli
}

func (cli *CLI) GetScanner() bufio.Scanner {
	return cli.scan
}

func (cli *CLI) GetCommand(name string) (*Command, error) {
	cmd, ok := cli.commands[name]
	if !ok {
		return nil, fmt.Errorf("requested command not found")
	}
	return cmd, nil
}

func (cli *CLI) EnvLoad(keys ...string) {
	for _, k := range keys {
		cli.env[k] = os.Getenv(k)
	}
}

func (cli *CLI) EnvSet(key, val string) {
	cli.env[key] = val
	os.Setenv(key, val)
}

func (cli *CLI) EnvUpload(vars map[string]string) {
	for k, v := range vars {
		cli.EnvSet(k, v)
	}
}

func (cli *CLI) EnvGet(key string) (string, error) {
	val, ok := cli.env[key]
	if !ok {
		return "", fmt.Errorf("%s not setted", key)
	}
	return val, nil
}

func (cli *CLI) RedirectErr(err error) {
	cli.Err <- err
	log.Printf("%s err detected", err)
}

func (cli *CLI) ExecWithContext(cmd *Command) error {
	err := cmd.Execute(cmd.ctx)
	if err != nil {
		return err
	}
	return nil

}

func (cli *CLI) RunCommand(name string, ctx *context.Context) error {
	cmd, err := cli.GetCommand(name)
	if err != nil {
		return err
	}
	cmd.ctx = ctx
	err = cli.ExecWithContext(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (cli *CLI) cmd(name string, ctx *context.Context) {
	err := cli.RunCommand(name, ctx)
	if err != nil {
		log.Println(err)
	}
}

func (cli *CLI) Run() {
	ctx, cancel := context.WithCancel(cli.ctx)
	for {
		select {
		case in := <-cli.Cmd:
			cli.cmd(in, &ctx)
		// case out := <-cli.Out:
		// log.Println(string(buff))
		// log.Printf("%b bytes writed", n)
		case err := <-cli.Err:
			// cli.err.Read([]byte(err.Error()))
			log.Println(err)
			cancel()
		}
	}
}
