module github.com/abtransitionit/gotask

// go toolchain version
go 1.24.2

// prod mode
require (
	github.com/abtransitionit/gocore v0.0.1
	github.com/abtransitionit/golinux v0.0.1
)

require (
	github.com/jedib0t/go-pretty/v6 v6.6.8 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

// used in dev mode - removes by CI at tag step - simplify development when working on several inter dependant projects

// direct dependency
replace github.com/abtransitionit/golinux => ../golinux

// indirect dependency
replace github.com/abtransitionit/gocore => ../gocore
