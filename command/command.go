package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/beto20/gofluence/chart"
	"github.com/beto20/gofluence/confluence"
	"github.com/beto20/gofluence/java"
	"github.com/beto20/gofluence/model"
	"github.com/beto20/gofluence/storage"
	"github.com/go-echarts/go-echarts/v2/components"
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
	container        string
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
	c.fs.StringVar(&c.connectionString, "cs", "connection string", "server:xxxxx")
	c.fs.StringVar(&c.container, "ct", "container name", "container")
	c.fs.StringVar(&c.repoFullName, "rn", "repository name", "example-confluence-repo")
	c.fs.StringVar(&c.prefix, "p", "prefix", "exmp")
	c.fs.StringVar(&c.token, "t", "token", "xxxxx")
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
	fmt.Printf("container: %s\n", b.container)
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
	imageRemotePath := print(data, bb.RepoFullName, b.connectionString, b.container)
	bb.ImageRemoteUrl = imageRemotePath
	confluence.Execute(bb, data)
	return nil
}

func print(documents []model.Document, projectName string, connectionString string, container string) string {
	fmt.Println("PRINT INIT")
	page := components.NewPage()
	page.AddCharts(
		chart.GenerateTreeChart(documents, projectName),
	)

	imageName := remove(projectName)
	fmt.Println("project-Name ", imageName)
	x := chart.NewSnapshotConfig(page.RenderContent(), imageName+".png", SetConfig(imageName))
	chart.MakeSnapshot(x)

	fmt.Println("PRINT END")

	return storage.UploadImage(container, connectionString, imageName+".png")
}

func remove(repoFullName string) string {
	if contains(repoFullName, "/") {
		repoName := strings.Split(repoFullName, "/")
		return repoName[1]
	}
	return repoFullName
}

func contains(value string, key string) bool {
	for _, v := range value {
		if string(v) == key {
			return true
		}
	}
	return false
}

func SetConfig(imageName string) chart.SnapshotConfigOption {
	return func(config *chart.SnapshotConfig) {
		config.Renderer = "canvas"
		config.Path = "."
		config.FileName = imageName
		config.Suffix = "png"
		config.Quality = 2
		config.HtmlPath = "."
		config.KeepHtml = true
	}
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
