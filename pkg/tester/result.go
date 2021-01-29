package tester

type Outcome struct {
	Name       string   `json:"name"`
	TestMethod string   `json:"test_method"`
	Fruits     []*Fruit `json:"fruits"`
}

type Fruit struct {
	TestTarget string    `json:"test_target"`
	Language   string    `json:"language"`
	Expect     string    `json:"expect"`
	Details    []*Detail `json:"details"`
}

type Detail struct {
	TestCase   string `json:"test_case"`
	Result     string `json:"result"`
	IsExpected bool   `json:"is_expected"`
	Time       string `json:"time"`
	Output     string `json:"output"`
	Error      string `json:"error"`
}
