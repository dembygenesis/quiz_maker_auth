package file

type FileNameAndTypes struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	Confirmed bool   `json:"confirmed"`
}
