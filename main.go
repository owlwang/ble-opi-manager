package main

import (
	"ble-opi-manager/internal/characteristics/command"
	"ble-opi-manager/internal/characteristics/wifi"
	"ble-opi-manager/internal/config"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"tinygo.org/x/bluetooth"
)

func main() {
	conf := config.GetConfig()
	fmt.Println("config:", conf)
	namespaceUUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte("lab.acrcloud.cn"))
	serviceUUID := uuid.NewSHA1(namespaceUUID, []byte(conf.Secret))
	serviceBleUUID, _ := bluetooth.ParseUUID(serviceUUID.String())
	localNameSuffix := uuid.NewMD5(uuid.Nil, []byte(mustGetMacAddr())).String()[:8]
	fmt.Println("namespaceUUID:", namespaceUUID)
	fmt.Println("serviceUUID:", serviceUUID)
	fmt.Println("serviceBleUUID:", serviceBleUUID)
	fmt.Println("localNameSuffix:", localNameSuffix)

	adapter := bluetooth.DefaultAdapter

	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	// Define the peripheral device info.
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "ACRCloud - " + localNameSuffix,
		ServiceUUIDs: []bluetooth.UUID{serviceBleUUID},
	}))

	// Start advertising
	must("start adv", adv.Start())

	// Add Services
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceBleUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			wifi.NewCharacteristicConfig(serviceUUID),
			command.NewCharacteristicConfig(serviceUUID),
		},
	}))

	fmt.Println("ble-opi-manager is running...")
	address, _ := adapter.Address()
	for {
		// Here is for Raspberry PI Bluetooth bug. Repeated rebroadcasts
		println("restart advertising...")
		must("stop adv", adv.Stop())
		must("start adv", adv.Start())
		println("BLE MAC /", address.MAC.String())
		time.Sleep(30 * time.Second)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func mustGetMacAddr() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, interf := range interfaces {
		a := interf.HardwareAddr.String()
		if a != "" {
			return a
		}
	}
	panic("no MAC address found.")
}
