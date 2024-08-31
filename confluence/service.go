package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/beto20/gofluence/model"
	"github.com/beto20/gofluence/tmpl"
)

const (
	BODY_PAGE         = "page"
	BODY_STORAGE      = "storage"
	REQ_PARAM_VERSION = "version"
	REQ_PARAM_SPACE   = "space"
	DEFAULT_EMPTY     = " "
	DEFAULT_INT       = 0
)

type confluenceDto struct {
	id          string
	version     int64
	spaceKey    string
	bitbucket   model.Bitbucket
	projectData []model.Document
}

func toConfluenceDto(b model.Bitbucket, documents []model.Document) confluenceDto {
	return confluenceDto{
		id:          DEFAULT_EMPTY,
		version:     DEFAULT_INT,
		spaceKey:    DEFAULT_EMPTY,
		bitbucket:   b,
		projectData: documents,
	}
}

func toTable(cdto confluenceDto) []tmpl.Table {
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

func Execute(b model.Bitbucket, documents []model.Document) {
	fmt.Println("Start Confluence Processing")
	cdto := toConfluenceDto(b, documents)
	findByTitle(&cdto)

	if cdto.id == "" || cdto.version == 0 {
		fmt.Println("Create New Page")
		// createPage(cdto)
	}

	fmt.Println("Update Existing Page")
	updatePage(cdto)

	fmt.Println("End Confluence Processing")
}

func updatePage(cdto confluenceDto) {
	content := tmpl.BuildPage(cdto.bitbucket.RepoFullName, "", toTable(cdto))

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

// const table_content = `
// 	<table>
// 		<tbody>
// 			<tr>
// 				<th>modulo</th>
// 				<th>version</th>
// 				<th>artifactory</th>
// 				<th>estado</th>
// 				<th>rama</th>
// 			</tr>
// 			<tr>
// 				<td>assi-ifx-associated-services</td>
// 				<td>0.0.1-6</td>
// 				<td>SNAPSHOT - RELEASE</td>
// 				<td>ESTABLE</td>
// 				<td>develop - master</td>
// 			</tr>
// 			<tr>
// 				<td>assi-ifx-customer</td>
// 				<td>0.0.1-7</td>
// 				<td>SNAPSHOT - RELEASE</td>
// 				<td>ESTABLE</td>
// 				<td>develop - master</td>
// 			</tr>
// 			<tr>
// 				<td>assi-ifx-savings-account</td>
// 				<td>0.0.1-6</td>
// 				<td>SNAPSHOT - RELEASE</td>
// 				<td>OBSOLETO</td>
// 				<td>develop - master</td>
// 			</tr>
// 		</tbody>
// 	</table>
// `
