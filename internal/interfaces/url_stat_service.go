package interfaces

import (
	"github.com/reactivex/rxgo/v2"
)

type URLStatService interface {
	GetURLStats(shortID string) rxgo.Observable
	RecordAccess(shortID string) rxgo.Observable
}
