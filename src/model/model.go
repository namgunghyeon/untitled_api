package model

type Project struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Type string `json:"type"`
  Path string `json:"path"`
  Name string `json:"name"`
  Count int `json:"count"`
}
