package pubsub

import (
	"github.com/micro/go-micro/v2"
	"github.com/sidazhang123/f10-go/plugins/zap"
	index "github.com/sidazhang123/f10-go/srv/index/model"
	"sync"
	"time"
)

var (
	indexService *index.Service
	pub          micro.Event
	log          = zap.GetLogger()
	lastRun      = LockTime{}
)

type LockTime struct {
	t  time.Time
	mu sync.RWMutex
}

func (t *LockTime) Time() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.t
}

func (t *LockTime) SetTime(tm time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t = tm
}

const interval = 10 * time.Minute

func Init(publisher micro.Event) {
	indexService, _ = index.GetService()
	pub = publisher
	lastRun.SetTime(time.Now().Add(-interval))
}
