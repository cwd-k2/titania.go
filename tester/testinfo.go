package tester

type TestInfo struct {
	UnitName string `json:"unit_name"`
	CaseName string `json:"case_name"`
	Language string `json:"language"`
	Result   string `json:"result"`
	Error    string `json:"error"`
	Time     string `json:"time"`
}
