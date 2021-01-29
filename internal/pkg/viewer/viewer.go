package viewer

import (
	p "github.com/cwd-k2/titania.go/internal/pkg/pretty"
)

type Viewer interface {
	Draw()
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

func NewFancyView(name string, codes, cases int, indices []string) *FancyView {
	return &FancyView{name, codes, cases, make([]int, codes), indices}
}

func NewQuietView(name string, total int) *QuietView {
	return &QuietView{name, total, 0}
}

func (view *FancyView) Draw() {

	p.Printf("%s\n", p.Bold(p.Cyan(view.name)))

	for _, index := range view.indices {
		p.Printf("[%s] %s %s\n", p.Yellow("WAIT"), "START", p.Bold(p.Blue(index)))
	}

}

func (view *FancyView) Update(position int) {
	view.counts[position]++

	p.Up(view.codes - position)
	p.Erase()

	if view.counts[position] == view.cases {
		p.Printf("[%s] ", p.Green("DONE"))
	} else {
		p.Printf("[%s] ", p.Yellow("WAIT"))
	}
	p.Printf("%02d/%02d %s", view.counts[position], view.cases, p.Bold(p.Blue(view.indices[position])))

	p.Down(view.codes - position)
	p.Beginning()
}

func (view *QuietView) Draw() {
	p.Printf("[%s] %s\n", p.Green("LAUNCH"), p.Bold(p.Cyan(view.name)))
}

func (view *QuietView) Update(_ int) {
	view.count++

	if view.count == view.total {
		p.Printf("[%s] %s\n", p.Yellow("FINISH"), p.Bold(p.Cyan(view.name)))
	}

}
