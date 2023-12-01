package app

import "github.com/rewired-gh/go-signal-server/internal/util"

type session struct {
	id      string
	inviter []string
	invitee []string
}

func newSession() session {
	return session{
		id:      util.GeneratePIN(),
		inviter: make([]string, 0),
		invitee: make([]string, 0),
	}
}
