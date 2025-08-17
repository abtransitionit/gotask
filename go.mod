module github.com/abtransitionit/gotask

// go toolchain version
go 1.24.2

// prod mode
require github.com/abtransitionit/golinux v1.0.0

require github.com/abtransitionit/gocore v1.0.0

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

// used in dev mode - removes by CI at tag step - simplify development when working on several inter dependant projects

// direct dependency
replace github.com/abtransitionit/golinux => ../golinux

// indirect dependency
replace github.com/abtransitionit/gocore => ../gocore
