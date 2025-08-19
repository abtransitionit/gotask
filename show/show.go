package show

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
)

// Name: Show
// Purpose: create file from file as sudo user or non-sudo user
func (w *Workflow) Show() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Phase", "Description", "Dependencies"})

	// deterministically order phases
	names := make([]string, 0, len(w.Phases))
	for name := range w.Phases {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		phase := w.Phases[name]
		deps := "none"
		if len(phase.Dependencies) > 0 {
			deps = fmt.Sprintf("%v", phase.Dependencies)
		}
		t.AppendRow(table.Row{phase.Name, phase.Description, deps})
	}

	t.Render()
}
