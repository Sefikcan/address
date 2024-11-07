package request

type AddressPatchRequest struct {
	Doc []PatchRequest `json:"doc"`
}

type PatchRequest struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}
