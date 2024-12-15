package dto

type UserState struct {
	WorkState     bool `json:"work_state"`
	SendingState  bool `json:"sending_state"`
	CheckingState bool `json:"checking_state"`
}

type SendingState struct {
	FirstPhoto  bool     `json:"first_photo"`
	SecondPhoto bool     `json:"second_photo"`
	Photo       []string `json:"photo"`
}
