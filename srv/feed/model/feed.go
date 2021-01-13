package model

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sidazhang123/f10-go/basic/config"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	proto "github.com/sidazhang123/f10-go/srv/feed/proto/feed"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type Service struct {
}

var (
	s      *Service
	m      sync.RWMutex
	Params = &Opts{}
	log    = z.GetLogger()
)

type Focus struct {
	Code      string
	RuleId    string
	Name      string
	Fetchtime string
	Keys      map[string]Contain
}
type Contain struct {
	Msg     string
	Contain []string
}

type FocusItem struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	Gentime       string
	Code          string
	Name          string
	Fetchtime     string
	Tabupdatetime string `json:"tabupdatetime, omitempty" bson:"tabupdatetime, omitempty"`
	Keys          string
	Chan          string
	Fav           int32
	Del           int32
}
type CreateRuleDTO struct {
	TarCol  string `json:"tarCol" bson:"tarCol"`
	Channel string
	Key     string
	Cond1   []string
	Cond2   []string
}
type ReadRuleDTO struct {
	TarCol  string `json:"tarCol" bson:"tarCol"`
	Channel string
	Key     string
	Cond1   []string
	Cond2   []string
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
}
type DaysToDel struct {
	Channel string
	Day     int
}
type feedService interface {
	//rule.go
	CreateRule(rules []*proto.Rule) (error, int)
	ReadRule(rules []*proto.Rule) (error, []ReadRuleDTO)
	UpdateRule(rules []*proto.Rule) (error, int)
	DeleteRule(rules []*proto.Rule) (error, int)
	//focus.go
	GenerateFocus(date string) error
	PurgeFocus(date string) (error, int)
	ToggleFocusDel(objectId string, v int32) error
	ToggleFocusFav(objectId string, v int32) error
	ReadFocus(date string, chanId string, del int32, fav int32) (error, string)
	DeleteOutdatedFocus() (error, int)
	DeleteRecoveryBin() (error, int)
	GetChanODay() (error, map[string]int32)
	SetChanODay([]interface{}) (error, int)
	GetFocusStat() (error, map[string]int32)
	//genCSV.go
	GenerateOperationalAnalysisDiffCSV(ts string) error
	//utils.go
	AddJPushReg(id string) error
	Log(msg string) error
	FindLatestFetchTime(collection string) (error, string)
}

func GetService() (service *Service, err error) {
	if s == nil {
		err = errors.Errorf("[GetService] GetService is not initialized.")
		return
	}
	return s, nil
}

func Init() {
	m.Lock()
	defer m.Unlock()
	if s != nil && (Params != &Opts{}) {
		return
	}

	err := config.C().Path("params", Params)
	if err != nil {
		err = fmt.Errorf("[Init] Failed to get params\n" + err.Error())
		log.Error(err.Error())
		return
	}

	s = &Service{}

}
