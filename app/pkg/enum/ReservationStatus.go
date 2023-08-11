package enum

type ReservationStatus int

const (
	Waiting ReservationStatus = iota
	InProgress
	Done
	_ // skip iota for 3
	_ // skip iota for 4
	Pending
	_ // skip iota for 6
	_ // skip iota for 7
	_ // skip iota for 8
	Canceled
)

var ReservationStatusNames = map[ReservationStatus]string{
	Waiting:    "Waiting",
	InProgress: "InProgress",
	Done:       "Done",
	Pending:    "Pending",
	Canceled:   "Canceled",
}

var ReservationStatusValues = map[string]ReservationStatus{
	"Waiting":    Waiting,
	"InProgress": InProgress,
	"Done":       Done,
	"Pending":    Pending,
	"Canceled":   Canceled,
}

func (r ReservationStatus) String() string {
	return [...]string{"未案内", "案内中", "案内済", "", "", "保留", "", "", "", "キャンセル"}[r]
}
