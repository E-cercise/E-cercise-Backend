package enum

type ImageState string

const (
	Temp    ImageState = "temp"
	Archive ImageState = "archive"
)

func (a ImageState) ToString() string {
	return string(a)
}
