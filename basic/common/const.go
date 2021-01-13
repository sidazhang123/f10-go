package common

var FlagNameToCollSuffix = map[string]string{
	"最新提示": "latest_tips",
	"公司概况": "company_introduction",
	"财务分析": "financial_analysis",
	"股东研究": "shareholder_analysis",
	"股本结构": "share_capital_structure",
	"资本运作": "capital_operation",
	"业内点评": "peer_review",
	"行业分析": "industrial_analysis",
	"公司大事": "company_news",
	"港澳特色": "gao_special",
	"经营分析": "operational_analysis",
	"主力追踪": "core_tracking",
	"分红扩股": "dividends",
	"高层治理": "upper_management",
	"龙虎榜单": "longhu_rankings",
	"关联个股": "related_stocks",
}
var CommandChain = map[string]string{
	IndexStart:      FetchStart,
	FetchStart:      ProcessStart,
	ProcessStart:    GenFeedStart,
	GenFeedStart:    AccumulateStart,
	AccumulateStart: "any process after Accumulate?",
}

const (
	TimestampLayout = "2006-01-02T15:04:05.000Z"
	LoggingTopic    = "sidazhang123.f10-go.logging"
	ControlTopic    = "sidazhang123.f10-go.control"
	LogInfolvl      = 0
	LogErrorlvl     = 1
	LogCallbacklvl  = 2
	IndexStart      = "Index"
	IndexComp       = "Index Complete"
	FetchStart      = "Fetch"
	FetchComp       = "Fetch Complete"
	ProcessStart    = "Process"
	ProcessComp     = "Process Complete"
	GenFeedStart    = "GenFeed"
	GenFeedComp     = "GenFeed Complete"
	AccumulateStart = "Accumulate"
	AccumulateComp  = "Accumulate Complete"

	DeleteOutdatedFocus = "DeleteOutdatedFocus"
)
