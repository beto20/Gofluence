package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/beto20/gofluence/model"
	"github.com/beto20/gofluence/tmpl"
)

const (
	BODY_PAGE              = "page"
	BODY_STORAGE           = "storage"
	REQ_PARAM_VERSION      = "version"
	REQ_PARAM_SPACE        = "space"
	DEFAULT_EMPTY          = " "
	DEFAULT_INT            = 0
	BRANCH                 = "master"
	DEPLOYMENT_RELEASE     = "Release"
	DEPLOYMENT_DEVELOPMENT = "Development - Snapshot"
)

type confluenceDto struct {
	id             string
	version        int64
	spaceKey       string
	bitbucket      model.Bitbucket
	projectData    []model.Document
	imageRemoteUrl string
}

func toConfluenceDto(b model.Bitbucket, documents []model.Document) confluenceDto {
	fmt.Println("b.ImageRemoteUrl", b.ImageRemoteUrl)
	return confluenceDto{
		id:             DEFAULT_EMPTY,
		version:        DEFAULT_INT,
		spaceKey:       DEFAULT_EMPTY,
		bitbucket:      b,
		projectData:    documents,
		imageRemoteUrl: b.ImageRemoteUrl,
	}
}

func toTableProject(cdto confluenceDto) []tmpl.Table {
	var tables []tmpl.Table

	for _, pd := range cdto.projectData {
		t := tmpl.Table{
			Module:      pd.Name,
			Version:     pd.Version,
			Environment: cdto.bitbucket.DeploymentEnvironment,
			Branch:      cdto.bitbucket.Branch,
			Commits:     cdto.bitbucket.Commit,
		}
		tables = append(tables, t)
	}

	return tables
}

func toTableDependency(documents []model.Document) []tmpl.TableDependency {
	var tableDep []tmpl.TableDependency

	for _, doc := range documents {
		for _, dep := range doc.Dependencies.Dependency {
			t := tmpl.TableDependency{
				GroupId:    dep.GroupId,
				ArtifactId: dep.ArtifactId,
				Version:    dep.Version,
				Scope:      dep.Scope,
			}
			tableDep = append(tableDep, t)
		}
	}

	return tableDep
}

func buildConfluencePage(cdto confluenceDto, documents []model.Document) string {
	deps := toTableDependency(documents)
	logs := toTableProject(cdto)
	return tmpl.BuildPage(cdto.bitbucket.RepoFullName, "", logs, deps, cdto.imageRemoteUrl)
}

func Execute(b model.Bitbucket, documents []model.Document) {
	fmt.Println("Start Confluence Processing")
	cdto := toConfluenceDto(b, documents)
	getRepositoryName(&cdto)
	getEnvorinment(&cdto)
	findByTitle(&cdto)

	if cdto.id == "" || cdto.version == 0 {
		fmt.Println("Create New Page")
		// createPage(cdto)
	}

	fmt.Println("Update Existing Page")
	updatePage(cdto, documents)
	fmt.Println("End Confluence Processing")
}

func getRepositoryName(cdto *confluenceDto) {
	if strings.Contains(cdto.bitbucket.RepoFullName, "/") {
		cdto.bitbucket.RepoFullName = strings.Split(cdto.bitbucket.RepoFullName, "/")[1]
	}
}

func getEnvorinment(cdto *confluenceDto) {
	if strings.Contains(cdto.bitbucket.Branch, BRANCH) {
		cdto.bitbucket.DeploymentEnvironment = DEPLOYMENT_RELEASE
	} else {
		cdto.bitbucket.DeploymentEnvironment = DEPLOYMENT_DEVELOPMENT
	}
}

func updatePage(cdto confluenceDto, documents []model.Document) {
	content := buildConfluencePage(cdto, documents)

	request := UpdatePageRequest{
		Id:       cdto.id,
		PageType: BODY_PAGE,
		Title:    cdto.bitbucket.RepoFullName,
		Space: struct {
			Key string `json:"key"`
		}{
			Key: cdto.spaceKey,
		},
		Body: struct {
			Storage struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			} `json:"storage"`
		}{
			Storage: struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			}{
				Value:          content,
				Representation: BODY_STORAGE,
			},
		},
		Version: struct {
			Number int64 `json:"number"`
		}{
			Number: cdto.version + 1,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Print("error marshalling request")
	}

	url := cdto.bitbucket.Url + "/content/" + cdto.id
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", "Basic "+cdto.bitbucket.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("ERROR: ", res.Status)
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
}

// TODO: implement
func createPage(cdto confluenceDto) {
	request := CreatePageRequest{
		PageType: BODY_PAGE,
		Title:    cdto.bitbucket.RepoFullName,
		Ancestors: []struct {
			Id string `json:"id"`
		}{
			{Id: "id"},
		},
		Space: struct {
			Key string `json:"key"`
		}{
			Key: "TC",
		},
		Body: struct {
			Storage struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			} `json:"storage"`
		}{
			Storage: struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			}{
				Value:          "<h1>PRUEBA</h1><ac:structured-macro ac:name=\"table\"><ac:rich-text-body><table><tbody><tr><th>Header 1</th><th>Header 2</th></tr><tr><td>Cell 1</td><td>Cell 2</td></tr><tr><td>Cell 3</td><td>Cell 4</td></tr></tbody></table></ac:rich-text-body></ac:structured-macro>",
				Representation: BODY_STORAGE,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Print("error marshalling request")
	}

	fmt.Println("jsondata: ", request)
	url := cdto.bitbucket.Url + "/content/" + cdto.id

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", "Basic "+cdto.bitbucket.Token)
	req.Header.Set("Content-Type", "application/json")

	// res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	// defer res.Body.Close()
	// if res.StatusCode != 200 {
	// 	fmt.Println("ERROR: ", res.Status)
	// }

	_, err = io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

}

func findByTitle(cdto *confluenceDto) {
	url := cdto.bitbucket.Url + "/content?title=" + cdto.bitbucket.RepoFullName + "&expand=" + REQ_PARAM_VERSION + "," + REQ_PARAM_SPACE
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+cdto.bitbucket.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var response FindByTitleResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Print("error marshalling request")
	}

	if len(response.Results) == 0 {
		fmt.Println("No content found")
	}

	cdto.id = response.Results[0].Id
	cdto.version = response.Results[0].Version.Number
	cdto.spaceKey = response.Results[0].Space.Key
}
