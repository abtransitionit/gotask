module github.com/abtransitionit/gotask

// go toolchain version
go 1.24.2

// prod mode
require (
	github.com/abtransitionit/gocore v0.0.1
	github.com/abtransitionit/golinux v0.0.1
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/jedib0t/go-pretty/v6 v6.6.8 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
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
