package profile

type Profile struct {
	ID         int64     `json:"u"`
	Categories []int64   `json:"a"`
	Weights    []float64 `json:"w"`
	Device     string    `json:"d"`
}
