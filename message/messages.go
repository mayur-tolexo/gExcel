package message

const (
	KindConnected = iota + 1
	KindUserJoined
	KindUserLeft
	KindStroke
	KindClear
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type User struct {
	ID    int    `json:"id"`
	Color string `json:"color"`
}

type Connected struct {
	UserID int              `json:"userId"`
	Kind   int              `json:"kind"`
	Color  string           `json:"color"`
	Users  []User           `json:"users"`
	Msg    map[string]EData `json:"msg"`
}

func NewConnected(userID int, color string, users []User, msg map[string]EData) *Connected {
	return &Connected{
		UserID: userID,
		Kind:   KindConnected,
		Color:  color,
		Users:  users,
		Msg:    msg,
	}
}

type UserJoined struct {
	Kind int  `json:"kind"`
	User User `json:"user"`
}

func NewUserJoined(userID int, color string) *UserJoined {
	return &UserJoined{
		Kind: KindUserJoined,
		User: User{ID: userID, Color: color},
	}
}

type UserLeft struct {
	Kind   int `json:"kind"`
	UserID int `json:"userId"`
}

func NewUserLeft(userID int) *UserLeft {
	return &UserLeft{
		Kind:   KindUserLeft,
		UserID: userID,
	}
}

//Stroke data strokes
type Stroke struct {
	Kind   int     `json:"kind"`
	UserID int     `json:"userId"`
	Points []Point `json:"points"`
	Finish bool    `json:"finish"`
	EDatas []EData `json:"eData"`
}

//EData excel data point
type EData struct {
	ID    string `json:"Id"`
	Value string `json:"value"`
}

type Clear struct {
	Kind   int    `json:"kind"`
	UserID int    `json:"userId"`
	ID     string `json:"Id"`
	Color  string `json:"color"`
}
