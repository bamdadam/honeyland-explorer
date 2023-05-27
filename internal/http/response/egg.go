package response

type EggsResponse struct {
	Type string `json:"type"`
	Eggs []Egg  `json:"eggs"`
}

type Egg struct {
	Id         uint64 `json:"id"`
	Type       string `json:"type"`
	Generation string `json:"generation"`
	Universe   string `json:"universe"`
	Like       string `json:"like"`
	Dislike    string `json:"dislike"`
}
