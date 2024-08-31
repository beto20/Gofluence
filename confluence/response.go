package confluence

type FindByTitleResponse struct {
	Results []struct {
		Id       string `json:"id"`
		PageType string `json:"type"`
		Status   string `json:"status"`
		Title    string `json:"title"`
		Version  struct {
			Number int64 `json:"number"`
		} `json:"version"`
		Space struct {
			Key string `json:"key"`
		} `json:"space"`
	} `json:"results"`
}
