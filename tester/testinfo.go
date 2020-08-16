package tester

type TestInfo struct {
	CaseName string `json:"case_name"`
	Result   string `json:"result"`
	Error    string `json:"error"`
	Time     string `json:"time"`
}
