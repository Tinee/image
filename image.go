package progimage

type Image struct {
	ID          string `json:"id"`
	ContentType string `json:"contentType"`
	Body        []byte `json:"body"`
}
