package tester

func Execute(directories, languages []string) []*ShowUnit {
	var outcomes []*ShowUnit

	for _, dirname := range directories {
		testUnit := NewTestUnit(dirname, languages)
		// 実行するテストがない
		if testUnit == nil {
			continue
		}

		fruits := testUnit.Exec()

		outcome := new(ShowUnit)
		outcome.Name = dirname
		outcome.Fruits = fruits
		outcomes = append(outcomes, outcome)
	}

	return outcomes
}
