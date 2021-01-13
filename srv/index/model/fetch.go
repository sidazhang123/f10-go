package model

import (
	"fmt"
	rabbitMq "github.com/sidazhang123/f10-go/plugins/rabbit-mq"
	proto "github.com/sidazhang123/f10-go/srv/index/proto/index"
	"github.com/sidazhang123/tdxF10Protocol-goVer"
	"strings"
)

func (s *Service) Fetch() (error, []*proto.Stock) {
	const expCodeVol = 3800
	const tolerance = 0.8
	// get code-name map
	err, codename := tdxF10Protocol_goVer.GetCodeNameFromSina()
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	if len(codename) < expCodeVol*tolerance {
		return fmt.Errorf(fmt.Sprintf("GetCodeNameFromSina returns #code: %d<%.0f", len(codename), expCodeVol*tolerance)), nil
	}
	// prepare the code list
	var codeSlice []string
	for code := range codename {
		codeSlice = append(codeSlice, code)
	}
	// get code-category map
	err, cateMap := getCompanyCategory(codeSlice)
	if err != nil {
		return err, nil
	}

	flagname := getFlagNameSlice()
	var stockInfo []*proto.Stock
	for code, paramMap := range cateMap {

		for _, flag := range flagname {
			infoSlice := paramMap[flag]

			stockInfo = append(stockInfo, &proto.Stock{
				Code:     code,
				Name:     codename[code],
				Flagname: flag,
				Filename: infoSlice[0],
				Start:    infoSlice[1],
				Length:   infoSlice[2],
			})
		}

	}
	stockInfoLen := float64(len(stockInfo))
	tolCatVol := float64(len(codename)*len(flagname)) * tolerance
	if stockInfoLen < tolCatVol {
		return fmt.Errorf(fmt.Sprintf("getCompanyCategory returns less than expected: %.0f < %.0f", stockInfoLen, tolCatVol)), nil
	}

	//send to mq
	mq := rabbitMq.GetRMQ()

	pub := rabbitMq.GetPubChan()
	_, err = mq.PurgeQueue()
	if err != nil {
		return err, nil
	}
	mq.RegisterPublisher(pub)

	for _, stock := range stockInfo {
		pub <- fmt.Sprintf("%s;%s;%s;%s;%s;%s", stock.GetCode(), stock.GetName(), stock.GetFlagname(),
			stock.GetFilename(), stock.GetStart(), stock.GetLength())
	}

	return nil, stockInfo
}

func getFlagNameSlice() (flagname []string) {
	// get the interested flagname Slice
	flagnameParam := strings.Split(strings.ReplaceAll(opts.FlagName, " ", ""), ",")
	for _, v := range flagnameParam {
		flagname = append(flagname, v)
	}
	return
}

func getCompanyCategory(codeSlice []string) (error, map[string]map[string][]string) {
	addrs, timeout, maxRetry := initApi()
	api := tdxF10Protocol_goVer.Socket{
		MaxRetry: maxRetry,
	}
	//fmt.Printf("[[after initApi]]%+v %d %d %d\n",addrs,len(addrs),timeout,maxRetry)
	api.Init(addrs, timeout)
	return api.GetCompanyInfoCategory(codeSlice)
}

func initApi() (addrs []string, timeout, maxRetry int) {
	if opts.Timeout < 0 {
		timeout = 0
	}
	if maxRetry < 0 {
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
