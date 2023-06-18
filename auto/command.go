package auto

type Command struct {
	Action  string `json:"action"`
	Path    string `json:"path"`
	Context string `json:"context"`
}
