package text

// same structures as in github.com/schorlet/gcp/ocr/process/image

type DetectText struct {
	Annotations []Annotation `json:"annotations"`
}

type Annotation struct {
	Vertices    []Vertex `json:"vertices"`
	Description string   `json:"description"`
}

type Vertex struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}
