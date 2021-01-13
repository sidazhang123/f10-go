package model

import (
	"crypto/md5"
	"fmt"
	log2 "github.com/micro/go-micro/v2/util/log"
	"github.com/panjf2000/ants"
	rabbit_mq "github.com/sidazhang123/f10-go/plugins/rabbit-mq"
	tdxF10Protocol_goVer "github.com/sidazhang123/tdxF10Protocol-goVer"
	"github.com/streadway/amqp"
	"regexp"
	"strings"
	"sync"
	"time"
)

type t struct {
	p  *ants.Pool
	mq *rabbit_mq.RabbitMQ
}

func (s *Service) GetByMQ() (err []error) {
	mq := rabbit_mq.GetRMQ()

	p, e := ants.NewPool(opts.Worker)

	if e != nil {
		err = append(err, fmt.Errorf("Failed to create a pool\n"+e.Error()))
		return
	}
	t := &t{p: p, mq: mq}

	notStarted := true
	var sigChan = make(chan struct{})
	// monitor routine
	go func() {
		c := 0
		preRunningWorkers := 0
		// loss tolerance
		ft := opts.Worker / 2
		for {
			time.Sleep(time.Second * 20)
			pRunning := t.p.Running()

			log2.Info(fmt.Sprintf("pool.running %d, idle %t, 1st %t, 2nd%t", pRunning, notStarted, !notStarted && (pRunning == 0 || c > 15), (pRunning == preRunningWorkers) && (pRunning <= ft)))

			if notStarted && (pRunning > 0) {
				notStarted = false
			} else {
				// when work is complete or frozen for 20*15=5min
				if !notStarted && (pRunning == 0 || c > 15) {
					close(sigChan)
					mq.PurgeQueue()
					break
				}
				// if #worker doesn't change and the loss is tolerable(typically < capacity of the pool)
				if (pRunning == preRunningWorkers) && (pRunning <= ft) {
					c += 1
				} else {
					preRunningWorkers = pRunning
					c = 0
				}
			}
		}
	}()
	log2.Info("get_html:51 RegisterSubscriber")
	e = mq.RegisterSubscriber(t, err, sigChan)
	if e != nil {
		err = append(err, fmt.Errorf("Failed to create a pool/close mq conn\n"+e.Error()))
		p.Release()
		return
	}
	p.Release()
	log2.Info("return  from GetByMQ")
	return
}

// from the 3-ele msg, have 5-ele code-flag msgs back into the queue, and put f=1 in db
func (t *t) Consumer(msg amqp.Delivery, errList []error, mux *sync.Mutex) {

	body := strings.Split(string(msg.Body), ";")
	//fmt.Sprintf("%s;%s;%s;%s;%s;%s", stock.GetCode(), stock.GetName(), stock.GetFlagname(),
	//stock.GetFilename(),stock.GetStart(), stock.GetLength())
	if len(body) == 6 {

		_ = t.p.Submit(func() {
			code, name, flagName, filename, start, length := body[0], body[1], body[2], body[3], body[4], body[5]
			addrs, timeout, maxRetry := initApi()
			api := tdxF10Protocol_goVer.Socket{
				MaxRetry: maxRetry,
			}
			api.Init(addrs, timeout)
			err, info := api.GetCompanyInfoContent(code, filename, start, length)
			if err != nil {
				e := fmt.Errorf("failed to get %s %s %s, %s", code, name, flagName, err.Error())
				log.Error(e.Error())
				rabbit_mq.Nack(msg)
				mux.Lock()
				errList = append(errList, e)
				mux.Unlock()
				return
			}

			err, updateTime := getUpdateTime(info)
			if err != nil {
				log.Error(err.Error())
				mux.Lock()
				errList = append(errList, err)
				mux.Unlock()
				rabbit_mq.Nack(msg)
				return
			}
			loc, _ := time.LoadLocation("Asia/Shanghai")
			now := time.Now().In(loc)
			y, m, d := now.Date()
			updateTimeWLocation, _ := time.ParseInLocation("2006-01-02", updateTime, time.UTC)

			toInsert := StockBody{
				Code:       code,
				Name:       name,
				FlagName:   flagName,
				Flag:       "",
				Body:       info,
				FetchTime:  now.Add(8 * time.Hour),
				UpdateTime: updateTimeWLocation,
				Uid:        fmt.Sprintf("%x", md5.Sum([]byte(code+fmt.Sprintf("%d-%d-%d", y, m, d)))),
			}
			err = InsertOne(toInsert, flagName)
			if err != nil && !strings.Contains(err.Error(), "dup key") {

				e := fmt.Errorf("%s, f=%s, %s", code, flagName, err.Error())
				log.Error(e.Error())
				mux.Lock()
				errList = append(errList, e)
				mux.Unlock()
				rabbit_mq.Nack(msg)
				return
			}
			rabbit_mq.Ack(msg)
		})

	} else {
		e := fmt.Errorf("body malformatted %d, %+v", len(body), body)
		log.Error(e.Error())
		mux.Lock()
		errList = append(errList, e)
		mux.Unlock()
		rabbit_mq.Nack(msg)

	}

}

func (s *Service) GetByInput() {}
func initApi() (addrs []string, timeout, maxRetry int) {

	if opts.Timeout < 0 {
		timeout = 0
	}
	if opts.MaxRetry < 0 {
		maxRetry = 0
	}
	addrSliceRaw := strings.Split(strings.ReplaceAll(opts.Addrs, " ", ""), ",")
	var addrSlice []string
	for _, i := range addrSliceRaw {
		if len(i) > 0 {
			addrSlice = append(addrSlice, i)
		}
	}
	if addrSlice != nil && len(addrSlice) > 0 {
		addrs = addrSlice
	}
	return

}

func getUpdateTime(s string) (error, string) {
	update := regexp.MustCompile(opts.GetUpdatetimeRegex).FindStringSubmatch(s)
	if update == nil {
		return fmt.Errorf("[getUpdateTime] failed to extract updateTime"), ""
	}
	return nil, update[1]
}
