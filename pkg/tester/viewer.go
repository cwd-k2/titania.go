package tester

import (
	"time"

	. "github.com/cwd-k2/titania.go/internal/pkg/pretty"
)

type Viewer interface {
	Init()
	Update(int)
	Done()
}

type QuietView struct {
	name string
}

type FancyView struct {
	name    string
	codes   int
	cases   int
	counts  []int
	indices []string
	time    time.Time
}

func NewView(name string, codes, cases int, indices []string) Viewer {
	if quiet {
		return NewQuietView(name)
	} else {
		return NewFancyView(name, codes, cases, indices)
	}
}

func NewFancyView(name string, codes, cases int, indices []string) *FancyView {
	return &FancyView{name, codes, cases, make([]int, codes), indices, time.Now()}
}

func NewQuietView(name string) *QuietView {
	return &QuietView{name}
}

func (view *FancyView) Init() {
	Printf("%s\n", Bold(Cyan(view.name)))

	for _, index := range view.indices {
		Printf("[%s] %s %s\n", Yellow("WAIT"), "START", Bold(Blue(index)))
	}
}

func (view *FancyView) Update(position int) {
	Up(view.codes - position)
	Erase()

	if view.counts[position]++; view.counts[position] == view.cases {
		Printf("[%s] ", Green("DONE"))
	} else {
		Printf("[%s] ", Yellow("WAIT"))
	}
	Printf("%02d/%02d %s", view.counts[position], view.cases, Bold(Blue(view.indices[position])))

	Down(view.codes - position)
	Beginning()
}

func (view *FancyView) Done() {
	Printf("[%s] Done in %ds.\n", Blue("INFO"), int(time.Now().Sub(view.time).Seconds()))
}

func (view *QuietView) Init() {
	Printf("[%s] %s %s\n", Green("LAUNCH"), time.Now().Format("15:04:05"), Bold(Cyan(view.name)))
}

func (view *QuietView) Update(_ int) {
	// Nothing to do
}

func (view *QuietView) Done() {
	Printf("[%s] %s %s\n", Yellow("FINISH"), time.Now().Format("15:04:05"), Bold(Cyan(view.name)))
}
