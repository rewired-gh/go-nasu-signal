package app

type sessionLookupRequest struct {
	ID string `json:"id"`
}

type sessionSetRequest struct {
	ID      string  `json:"id"`
	Inviter *string `json:"inviter"`
	Invitee *string `json:"invitee"`
}
