package controller

type Error struct {
	Err string `json:"error"`
}

func ErrorJson(err error) Error {
	return Error{err.Error()}
}
