package wifi

import (
	"ble-opi-manager/internal/shell"
	"ble-opi-manager/internal/transport"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"tinygo.org/x/bluetooth"
)

type wifiInfo struct {
	SSID string // wifi name
	PSK  string // wifi password
}

type request struct {
	RefreshOnly bool
	WIFI        wifiInfo
	StaticIP    string
	Router      string
}

type state struct {
	CurrentWIFI string // current wifi ssid
	CurrentIP   string // current ip address
	Router      string
	StaticIP    string // configured static ip

	WIFI wifiInfo
}

func getState() state {
	return state{
		CurrentWIFI: getCurrentWIFI(),
		CurrentIP:   getCurrentIP(),
		Router:      getRouter(),
	}
}

func getCommandOutput(command string) string {
	out, err := shell.RunCommand(command)
	if len(out) == 0 || err != nil {
		return ""
	}
	return out[:len(out)-1] // remove suffix `\n`
}

func getCurrentIP() string {
	return getCommandOutput("ifconfig wlan0 | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\\.){3}[0-9]*).*/\\2/p'")
}

func getCurrentWIFI() string {
	return getCommandOutput("iwgetid -r")
}

func getScanWIFI() string {
	wifis := strings.Split(getCommandOutput("iw dev wlan0 scan |grep SSID: |cut -d':' -f 2"), "\n")
	//trim
	for i := range wifis {
		wifis[i] = strings.TrimSpace(wifis[i])
	}
	return strings.Join(wifis, "\n")
}

func getRouter() string {
	return getCommandOutput("netstat -nr | awk '$1 == \"0.0.0.0\"{print$2}'")
}

func setWifiInfo(wifi wifiInfo) {
	if wifi.SSID == "" || wifi.PSK == "" || len(wifi.PSK) < 8 {
		return
	}
	shell.RunCommand("nmcli dev wifi")
	cmd := fmt.Sprintf("nmcli dev wifi connect '%s' password '%s'", wifi.SSID, wifi.PSK)
	fmt.Println("connect wifi", cmd)
	shell.RunCommand(cmd)
}

func NewCharacteristicConfig(serviceUUID uuid.UUID) bluetooth.CharacteristicConfig {
	wifiCharUUID := uuid.NewSHA1(serviceUUID, []byte("wifi"))
	wifiCharBleUUID, _ := bluetooth.ParseUUID(wifiCharUUID.String())

	reader := &transport.ReadHandler{}
	writer := transport.NewWriteHandler(func(uuid, content []byte) {
		req := &request{}
		json.Unmarshal(content, req)
		current := getState()
		fmt.Println("current := getState(), req", current, "|#|", req)
		if req.RefreshOnly {
			bytes, err := json.Marshal(current)
			if err != nil {
				bytes = []byte(err.Error())
			}
			reader = transport.NewReadHandler(uuid, bytes)
			return
		}
		if req.WIFI.SSID != current.CurrentWIFI {
			setWifiInfo(req.WIFI)
		}
	})

	return bluetooth.CharacteristicConfig{
		UUID:  wifiCharBleUUID,
		Flags: bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicReadPermission,
		ReadEvent: func(client bluetooth.Connection) ([]byte, error) {
			return reader.Read(), nil
		},
		WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
			writer.Write(value)
		},
	}
}
