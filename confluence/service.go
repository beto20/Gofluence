package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/beto20/gofluence/model"
)

func Execute(b model.Bitbucket) {

	// fmt.Println("request:", string(jsonData))

	// Post(jsonData)

	page := findByTitle("", "", b)

	updatePage(b, page.Results[0].Id, page.Results[0].Version.Number)
}

func updatePage(b model.Bitbucket, id string, version int64) {
	// content := fmt.Sprintf("<h1>%s</h1><p>This is a sample page content. branch: %s, build: %s, env: %s </p>", b.RepoFullName, b.Branch, b.BuildNumber, b.DeploymentEnvironment)

	request := UpdatePageRequest{
		Id:       id,
		PageType: "page",
		Title:    b.RepoFullName,
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
				Value:          "<p>MOCK123</p>",
				Representation: "storage",
			},
		},
		Version: struct {
			Number int64 `json:"number"`
		}{
			Number: version + 1,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Print("error marshalling request")
	}

	fmt.Println("jsondata: ", request)
	url := b.Url + "/content/" + id
	fmt.Println("url: ", url)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", "Basic "+b.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("ERROR: ", res.Status)
		fmt.Println("ERROR: ", err)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func createPage() {
	// request := CreatePageRequest{
	// 	PageType: "page",
	// 	Title:    b.RepoFullName,
	// 	Ancestors: []struct {
	// 		Id string `json:"id"`
	// 	}{
	// 		{Id: id},
	// 	},
	// 	Space: struct {
	// 		Key string `json:"key"`
	// 	}{
	// 		Key: "TC",
	// 	},
	// 	Body: struct {
	// 		Storage struct {
	// 			Value          string `json:"value"`
	// 			Representation string `json:"representation"`
	// 		} `json:"storage"`
	// 	}{
	// 		Storage: struct {
	// 			Value          string `json:"value"`
	// 			Representation string `json:"representation"`
	// 		}{
	// 			Value:          "<h1>PRUEBA</h1><ac:structured-macro ac:name=\"table\"><ac:rich-text-body><table><tbody><tr><th>Header 1</th><th>Header 2</th></tr><tr><td>Cell 1</td><td>Cell 2</td></tr><tr><td>Cell 3</td><td>Cell 4</td></tr></tbody></table></ac:rich-text-body></ac:structured-macro>",
	// 			Representation: "storage",
	// 		},
	// 	},
	// }
}

func findByTitle(spaceKey string, title string, b model.Bitbucket) Response {
	fmt.Println("TEST: ", b.Url)
	url := b.Url + "/content?spaceKey=TC&title=prueba-1&expand=version"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+b.Token)
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

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Print("error marshalling request")
	}

	if len(response.Results) == 0 {
		fmt.Println("No content found")
	}

	fmt.Println("res:", response)

	return response
}

type Response struct {
	Results []struct {
		Id       string `json:"id"`
		PageType string `json:"type"`
		Status   string `json:"status"`
		Title    string `json:"title"`
		Version  struct {
			Number int64 `json:"number"`
		} `json:"version"`
	} `json:"results"`
}
