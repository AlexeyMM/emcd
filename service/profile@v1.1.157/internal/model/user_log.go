package model

const (
	UserLogChangeTypeChangeAddress = "address"
)

type UserLog struct {
	ID            int
	UserID        int32
	ChangeType    string
	IP            string
	Token         string
	OldValue      string
	Value         string
	Active        bool
	IsSegmentSent bool
}
