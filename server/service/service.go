package service

type Error struct {
	Msg    string
	Status int
}

func (err *Error) Error() string {
	return err.Msg
}
