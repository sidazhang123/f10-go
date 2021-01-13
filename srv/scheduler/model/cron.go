/*
	Read cron job conf file and set the tasks accordingly :
		name	func to call,
		every	int,
		unit	Second,Minute,Hour...,
		at		run at X every 2 days,
		safe	suppress any panic from the callee func
	with a unique timezone setting:
    	loc		i.e."America/New_York"
	It also logs Error & Info to stdout and the log directory.
	Use *Service.SetJobs() to reload tasks, *Service.RemoveTask(taskName string) to cancel a schedule,
	and *Service.ClearCron() to cancel wipe the schedule.

	We run tasks without explicit arguments because which are meant to be in conf files or ENV vars

	micro --registry=consul call sidazhang123.f10.srv.scheduler Scheduler.UpdateSchedule
*/
package model

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	log2 "github.com/micro/go-micro/v2/util/log"
	"github.com/sidazhang123/f10-go/basic/config"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/scheduler/pubsub"
	"reflect"
	"strconv"
	"time"
)

var (
	log   = z.GetLogger()
	tasks []Task
	tz    = TZ{}
)

type TZ struct {
	Loc string `json:"loc"`
}
type Task struct {
	Name  string `json:"name"`
	Every uint64 `json:"every"`
	Unit  string `json:"unit"`
	At    string `json:"at"`
	Safe  bool   `json:"safe"`
}

func (s *Service) SetCron() (err error, name string, t time.Time) {
	tasks = getConfs()
	err, name, t = s.setJobs()
	return

}

//get all cron conf
func getConfs() (ret []Task) {

	err := config.C().App("time_zone", &tz)
	if err != nil {
		log.Error("[Cron] Can't marshal TZ" + "\n" + err.Error())
		return
	}
	i := 1
	for {
		t := Task{}
		err := config.C().App("task"+strconv.Itoa(i), &t)
		if err != nil {
			log.Error("[Cron] Can't marshal task" + strconv.Itoa(i) + "\n" + err.Error())
			break
		}
		if t.Name == "" {
			break
		}
		ret = append(ret, t)
		log2.Logf("[Task interval] %v", t.Every)
		i += 1
	}
	return
}

//set them accordingly
func (s *Service) setJobs() (err error, name string, t time.Time) {
	//external call to re-init
	gocron.Clear()
	for _, task := range tasks {

		err = s.setJob(task)
		if err != nil {
			return
		}

	}
	name, t = s.NextScheduledTask()
	return
}

func (s *Service) setJob(t Task) (e error) {
	name := t.Name
	log.Info("[setJob] call method with name - " + name)
	v := reflect.ValueOf(s)
	v = v.MethodByName(name)
	if !v.IsValid() {
		e = fmt.Errorf("[setJob] Invalid MethodName - " + name)
		return
	}
	//gocron.Every(1).Day().At("10:30:00").Do(task)
	everyResVal := gocron.Every(t.Every)
	u := t.Unit
	if t.Every > 1 {
		u += "s"
	}
	unitMethod := reflect.ValueOf(everyResVal).MethodByName(u)
	if !unitMethod.IsValid() {
		e = fmt.Errorf("[setJob] Invalid Unit - " + u)
		return
	}
	unitResInterface := unitMethod.Call(nil)[0].Interface().(*gocron.Job)
	var doVal, doMethod reflect.Value

	// At is valid it's based on Day
	atTime := t.At
	if t.Unit == "Day" {
		if atTime == "" {
			atTime = "00:00:00"
		}
		if len(atTime) == 5 {
			atTime += ":00"
		}
		atMethod := reflect.ValueOf(unitResInterface).MethodByName("At")
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(atTime)
		atResInterface := atMethod.Call(params)[0].Interface().(*gocron.Job)
		doVal = reflect.ValueOf(atResInterface)
	} else {
		doVal = reflect.ValueOf(unitResInterface)
	}
	if !doVal.IsValid() {
		e = fmt.Errorf("[setJob] unitResInterface err - ")
		return
	}

	if t.Safe {
		doMethod = doVal.MethodByName("DoSafely")
	} else {
		doMethod = doVal.MethodByName("Do")
	}

	if !doMethod.IsValid() {
		e = fmt.Errorf("[setJob] doVal/dosafely err - ")
		return
	}
	params := make([]reflect.Value, 1)
	params[0] = v
	doMethod.Call(params)
	if tz.Loc != "" {
		loc, err := time.LoadLocation(tz.Loc)
		if err != nil {
			e = fmt.Errorf("[setJob] Failed to load TimeZone - %v", err)
			return
		}
		gocron.ChangeLoc(loc)
	}
	go func() { <-gocron.Start() }()
	return
}

func (s *Service) RemoveTask(funcName string) (e error, name string, t time.Time) {
	v := reflect.ValueOf(s).MethodByName(funcName)
	if !v.IsValid() {
		e = fmt.Errorf("[RemoveTask] Invalid MethodName - " + funcName)
		return
	}
	gocron.Remove(v.Interface().(func()))
	name, t = s.NextScheduledTask()
	return
}

func (s *Service) ClearCron() (string, time.Time) {
	gocron.Clear()
	return s.NextScheduledTask()
}

func (s *Service) Once(funcName string) (e error, name string, t time.Time) {
	v := reflect.ValueOf(s).MethodByName(funcName)
	if !v.IsValid() {
		e = fmt.Errorf("[Once] Invalid MethodName - " + funcName)
		return
	}
	v.Call(nil)
	name, t = s.NextScheduledTask()
	return
}

func (s *Service) NextScheduledTask() (string, time.Time) {
	job, t := gocron.NextRun()
	if job == nil {
		return "", time.Time{}
	}
	val := reflect.ValueOf(*job)
	return val.FieldByName("jobFunc").String(), t
}

func (s *Service) SendAlarm(msg string) error {
	return pubsub.SendAlarm(msg)
}
