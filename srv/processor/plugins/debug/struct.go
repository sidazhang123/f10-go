package debug

type Debug struct{}
type debug interface {
	Financial_analysis(string) (map[string]interface{}, []error)
	Latest_tips(string) (map[string]interface{}, []error)
	Shareholder_analysis(string) (map[string]interface{}, []error)
	Operational_analysis(string) (map[string]interface{}, []error)
}
