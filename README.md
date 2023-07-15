

# ble-opi-manager

Manage Orange Pi using Bluetooth Low Energy.
This project is a modified version of ble-raspi-manager https://github.com/DiscreteTom/ble-raspi-manager

## Features

- WIFI management.
- Run shell commands.

## Build

```bash
go build .
```

## Installation

```bash
# run the following script as root
sudo -i

# create a folder
mkdir /root/bom
cd /root/bom


# install the service
cp bom.service /etc/systemd/system/

# reload systemd
systemctl daemon-reload

# optional: modify config
# vim /root/brm/config.json

# start the service
systemctl start bom

# start the service on system startup
systemctl enable bom
```
