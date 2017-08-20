package model

type Project struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Type string `json:"type"`
  Path string `json:"path"`
  Name string `json:"name"`
  Count int `json:"count"`
}

type KeywordIndex struct {
  Keyword string `json:"keyword"`
	Project string `json:"proejct"`
}

type Keyword struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Type string `json:"type"`
  Path string `json:"path"`
  KeywordIndex string `json:"keyword_index"`
  Count int `json:"count"`
}

type ProjectInfo struct {
	Name string `json:"name"`
	Color string `json:"color"`
}
