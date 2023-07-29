package enum

type WebSocketKey int

const (
	Todo WebSocketKey = iota
	Reservation
)

func (w WebSocketKey) String() string {
	return [...]string{"todo", "reservation"}[w]
}
