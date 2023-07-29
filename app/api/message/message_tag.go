package message

type Tag string

func (c Tag) String() string {
	return string(c)
}

var SUCCESS = Tag("success")
