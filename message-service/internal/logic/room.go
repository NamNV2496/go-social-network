package logic

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
	LastMsg string             `json:"last_msg"`
	OldMsg  []*Message         `json:"old_msg"`
	Public  int                `json:"public"`
	Members []string           `json:"members"`
}
