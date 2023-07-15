module ble-opi-manager

go 1.17

require (
	github.com/google/uuid v1.3.0
	tinygo.org/x/bluetooth v0.7.0
)

require (
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.0.3 // indirect
	github.com/muka/go-bluetooth v0.0.0-20220830075246-0746e3a1ea53 // indirect
	github.com/saltosystems/winrt-go v0.0.0-20230613063811-c792451fa808 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tinygo-org/cbgo v0.0.4 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
)

//replace tinygo.org/x/bluetooth => github.com/DiscreteTom/bluetooth v0.4.1-0.20220312031942-f13a0ecd4975
replace tinygo.org/x/bluetooth => github.com/owlwang/bluetooth v0.7.1-0.20230710074601-6c4f8a447b0b
