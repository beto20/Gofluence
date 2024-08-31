package confluence

type CreatePageRequest struct {
	PageType  string `json:"type"`
	Title     string `json:"title"`
	Ancestors []struct {
		Id string `json:"id"`
	} `json:"ancestors"`
	Space struct {
		Key string `json:"key"`
	} `json:"space"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
}

type UpdatePageRequest struct {
	Id       string `json:"id"`
	PageType string `json:"type"`
	Title    string `json:"title"`
	Space    struct {
		Key string `json:"key"`
	} `json:"space"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
	Version struct {
		Number int64 `json:"number"`
	} `json:"version"`
}
