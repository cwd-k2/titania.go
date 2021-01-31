package tester

import (
	"time"

	. "github.com/cwd-k2/titania.go/internal/pkg/pretty"
)

type Viewer interface {
	Init()
	Update(int)
}

type QuietView struct {
	name  string
	total int
	count int
}

type FancyView struct {
	name    string
	codes   int
	cases   int
	counts  []int
	indices []string
}

func NewView(name string, codes, cases int, indices []string) Viewer {
	if quiet {
		return NewQuietView(name, codes*cases)
	} else {
		return NewFancyView(name, codes, cases, indices)
	}
}

func NewFancyView(name string, codes, cases int, indices []string) *FancyView {
	return &FancyView{name, codes, cases, make([]int, codes), indices}
}

func NewQuietView(name string, total int) *QuietView {
	return &QuietView{name, total, 0}
}

func (view *FancyView) Init() {

	Printf("%s\n", Bold(Cyan(view.name)))

	for _, index := range view.indices {
		Printf("[%s] %s %s\n", Yellow("WAIT"), "START", Bold(Blue(index)))
	}

}

func (view *FancyView) Update(position int) {
	view.counts[position]++

	Up(view.codes - position)
	Erase()

	if view.counts[position] == view.cases {
		Printf("[%s] ", Green("DONE"))
	} else {
		Printf("[%s] ", Yellow("WAIT"))
	}
	Printf("%02d/%02d %s", view.counts[position], view.cases, Bold(Blue(view.indices[position])))

	Down(view.codes - position)
	Beginning()
}

func (view *QuietView) Init() {
	Printf("[%s] %s %s\n", Green("LAUNCH"), time.Now().Format("15:04:05"), Bold(Cyan(view.name)))
}

func (view *QuietView) Update(_ int) {
	view.count++

	if view.count == view.total {
		Printf("[%s] %s %s\n", Yellow("FINISH"), time.Now().Format("15:04:05"), Bold(Cyan(view.name)))
	}
}
