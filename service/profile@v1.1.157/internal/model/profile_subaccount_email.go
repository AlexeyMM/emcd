package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SubUserEmail тип для хранения email подпользователя.
//
// SubUserEmail должен формироваться по правилу:
//
//	<email родительского пользователя>.<timestamp>
type SubUserEmail string

func NewSubUserEmailByParentEmail(parentEmail string) SubUserEmail {
	// мы должны обеспечить уникальность email'ов subaccount'ов.
	// для этого добавим timestamp в конец email'а.
	// длина добавляемого сегмента не имеет значениея, т.к. для email'а не ограничена
	return SubUserEmail(fmt.Sprintf("%s.%d", parentEmail, time.Now().UnixMilli()))
}

func NewSubUserEmailFromString(s string) SubUserEmail {
	return SubUserEmail(s)
}

func (s SubUserEmail) String() string {
	return strings.ToLower(string(s))
}

// Index возвращает индекс подпользователя:
// - для новых подпользователей (c 10.10.24) это timestamp
// - для существующих подпользователей это порядковый номер
// - не всегда порядоквый номер старых пользователей - валидное число
func (s SubUserEmail) Index() int64 {
	if s == "" {
		return 0
	}

	i := strings.LastIndex(s.String(), ".")
	if i == -1 {
		return 0
	}

	index, err := strconv.ParseInt(s.String()[i+1:], 10, 64)
	if err != nil {
		return 0
	}

	return index
}
