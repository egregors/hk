# hk

Integration temperature and humidity sensor (bosch BME280) with Raspberry Pi and Apple HomeKit

## tl;dr

In general, this project is just 3 parts combined together:

* Raspberry Pi + BME280 sensor integration – to get temperature and humidity data
* Simple web server – to expose the data to local network
* HAP server – to expose the data to HomeKit in order to be able to use it in Home app

// TODO: add 3 screenshots here

## Quick start

// TODO: how to run local

```shell
sudo apt-get update
sudo apt-get upgrade
```
How to enable I2C bus on RPi device: If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS level. Go to "Interfaceing Options" menu, to active I2C bus. Probably you will need to reboot to load i2c kernel module. Finally you should have device like /dev/i2c-1 present in the system.

```shell
egregors@pi:~ $ i2cdetect -y 1
     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
00:                         -- -- -- -- -- -- -- --
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
70: -- -- -- -- -- -- -- 77
```

```shell
sudo apt install snapd
sudo snap install go --classic
go version
go1.23.2 linux/arm64
```


// TOOD: how to build and run prod
https://github.com/d2r2/go-bsbmp