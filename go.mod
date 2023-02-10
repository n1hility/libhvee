module github.com/baude/hyperv_kvp

go 1.19

require (
	github.com/n1hility/hypervctl v0.0.0-20230210055924-91c3470fe725
	golang.org/x/sys v0.5.0
)

replace github.com/baude/hyperv_kvp => ../../baude/hypervctl

require (
	github.com/drtimf/wmi v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
)
