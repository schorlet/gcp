package image

import (
	"context"
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"
)

type DetectText struct {
	Annotations []Annotation `json:"annotations"`
	object      string
}

type Annotation struct {
	Vertices    []Vertex `json:"vertices"`
	Description string   `json:"description"`
}

type Vertex struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

const maxResults = 30

func detect(ctx context.Context, bucket, object string) (DetectText, error) {
	var dt DetectText

	if err := setup(ctx); err != nil {
		return dt, fmt.Errorf("setup: %v", err)
	}

	// an image that refers to an object in Google Cloud Storage
	dt.object = fmt.Sprintf("%s/%s", bucket, object)
	img := vision.NewImageFromURI(fmt.Sprintf("gs://%s", dt.object))

	// detect texts
	texts, err := global.VisionClient.DetectTexts(ctx, img, nil, maxResults)
	if err != nil {
		return dt, fmt.Errorf("detect: %v", err)
	}

	dt.Annotations = make([]Annotation, 0, maxResults)

	for _, text := range texts {
		vertices := make([]Vertex, len(text.BoundingPoly.Vertices))
		for i, vertex := range text.BoundingPoly.Vertices {
			vertices[i] = Vertex{X: vertex.X, Y: vertex.Y}
		}

		a := Annotation{
			Vertices:    vertices,
			Description: text.Description,
		}
		dt.Annotations = append(dt.Annotations, a)
	}

	return dt, nil
}
