package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/beto20/gofluence/confluence"
	"github.com/beto20/gofluence/model"
)

type Runner interface {
	Init([]string) error
	Name() string
	Run(sc string) error
}

type Command struct {
	fs                    *flag.FlagSet
	branch                string
	buildNumber           string
	commit                string
	deploymentEnvironment string
	repoFullName          string
	prefix                string
	token                 string
	url                   string
}

func BuildCommand(cm string) *Command {
	cmd := &Command{
		fs: flag.NewFlagSet(cm, flag.ContinueOnError),
	}

	return genericFlags(cmd)
}

func genericFlags(c *Command) *Command {
	c.fs.StringVar(&c.branch, "b", "branch", "master")
	c.fs.StringVar(&c.buildNumber, "n", "build number", "20")
	c.fs.StringVar(&c.commit, "c", "commit", "bf372e2")
	c.fs.StringVar(&c.deploymentEnvironment, "e", "deployment environment", "development")
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
	fmt.Printf("buildNumber: %s\n", b.buildNumber)
	fmt.Printf("commit: %s\n", b.commit)
	fmt.Printf("deploymentEnvironment: %s\n", b.deploymentEnvironment)
	fmt.Printf("repoFullName: %s\n", b.repoFullName)
	fmt.Printf("prefix: %s\n", b.prefix)
	fmt.Printf("token: %s\n", b.token)
	fmt.Printf("url: %s\n", b.url)

	bb := model.Bitbucket{
		Branch:                b.branch,
		BuildNumber:           b.buildNumber,
		Commit:                b.commit,
		DeploymentEnvironment: b.deploymentEnvironment,
		RepoFullName:          b.repoFullName,
		Prefix:                b.prefix,
		Token:                 b.token,
		Url:                   b.url,
	}

	// java.ReadJavaProject(bb.Prefix)
	confluence.Execute(bb)
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
