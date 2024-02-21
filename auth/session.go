package auth

type Session struct {
	Id string
}

var sessions map[string]Session = map[string]Session{}

func SetSession(id string, session Session) {
	sessions[id] = session
}

func GetSession(id string) Session {
	return sessions[id]
}

func ClearSession(id string) {
	delete(sessions, id)
}
