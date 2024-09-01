package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/beto20/gofluence/confluence"
	"github.com/beto20/gofluence/java"
	"github.com/beto20/gofluence/model"
)

type Runner interface {
	Init([]string) error
	Name() string
	Run(sc string) error
}

type Command struct {
	fs               *flag.FlagSet
	branch           string
	commit           string
	connectionString string
	repoFullName     string
	prefix           string
	token            string
	url              string
}

func BuildCommand(cm string) *Command {
	cmd := &Command{
		fs: flag.NewFlagSet(cm, flag.ContinueOnError),
	}

	return genericFlags(cmd)
}

func genericFlags(c *Command) *Command {
	c.fs.StringVar(&c.branch, "b", "branch", "master")
	c.fs.StringVar(&c.commit, "c", "commit", "bf372e2")
	c.fs.StringVar(&c.connectionString, "cs", "connection string", "https://connection.storage.example/")
	c.fs.StringVar(&c.repoFullName, "rn", "repository name", "example-confluence-repo")
	c.fs.StringVar(&c.prefix, "p", "prefix", "exmp")
	c.fs.StringVar(&c.token, "t", "token", "9ycC5jb20ucGU6QVR")
	c.fs.StringVar(&c.url, "u", "url", "https://localhost:8080/")

	return c
}

func (b *Command) Init(args []string) error {
	return b.fs.Parse(args)
}

func (b *Command) Name() string {
	return b.fs.Name()
}

func (b *Command) Run(sc string) error {
	fmt.Printf("branch: %s\n", b.branch)
	fmt.Printf("commit: %s\n", b.commit)
	fmt.Printf("connectionString: %s\n", b.connectionString)
	fmt.Printf("repoFullName: %s\n", b.repoFullName)
	fmt.Printf("prefix: %s\n", b.prefix)
	fmt.Printf("token: %s\n", b.token)
	fmt.Printf("url: %s\n", b.url)

	bb := model.Bitbucket{
		Branch:       b.branch,
		Commit:       b.commit,
		RepoFullName: b.repoFullName,
		Token:        b.token,
		Url:          b.url,
	}

	data := java.ReadJavaProject(b.prefix)
	confluence.Execute(bb, data)
	return nil
}

func Root(args []string) error {
	if len(args) < 1 {
		fmt.Print("No inputs")
	}

	command := os.Args[0]
	fmt.Println(command)
	cmds := []Runner{
		BuildCommand(command),
	}

	for _, cmd := range cmds {
		if cmd.Name() == command {
			cmd.Init(os.Args[1:])
			return cmd.Run(os.Args[0])
		}
	}

	return fmt.Errorf("Error")
}
