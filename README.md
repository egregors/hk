# 🇭🇰 hk – HomeKit BME280 sensor integration + web iface

Integration temperature and humidity sensor (bosch BME280) with Raspberry Pi and Apple HomeKit

## tl;dr

In general, this project is just 3 parts combined:

* Raspberry Pi + BME280 sensor integration – to get temperature and humidity data
* Simple web server – to expose the data to local network
* HAP server – to expose the data to HomeKit in order to be able to use it in Home app

### Features:

* Temperature and humidity data from BME280 sensor
* Web server to expose the data
* Simple metrics collection with retention and Braille graph
* Backup and restore collected data (gob dump)
* Logging
* HomeKit integration
* Custom PIN for HomeKit

### Screenshots:
<p align="center">
     <img width="80%" alt="Screenshot 2024-11-29 at 18 29 12" src="https://github.com/user-attachments/assets/04885fc5-93d6-4e4e-8144-925465979f4d">
</p>
<p align="center">
     <img width="30%" alt="Screenshot 2024-11-29 at 18 29 12" src="https://github.com/user-attachments/assets/5cb203cb-40ca-4f28-b27c-affc23ae93a1">
     <img width="30%" alt="Screenshot 2024-11-29 at 18 29 12" src="https://github.com/user-attachments/assets/5657c972-2e06-4b8d-874d-46b81c3f87ea">
</p>

## Quick start

Actual sensor code should be compiled on Raspberry Pi device. You can use any other device to compile the code, but you
should copy the binary to Raspberry Pi device to run it.
In purpose to make development process less painful there is a `ClimateSensor` interface. So, you need to install

```shell
egregors@pi:~/Github/hk $ gh repo sync
✓ Synced the "main" branch from "egregors/hk" to local repository
egregors@pi:~/Github/hk $ make build
mv ./"t-hk-srv" ~/go/bin/
egregors@pi:~/Github/hk $ cd ~/go/bin/
egregors@pi:~/go/bin $ sudo ./restart.sh
prev log:
[INFO] 2024/11/29 15:45:56 main.go:30: 🇭🇰 revision: 830083e
[INFO] 2024/11/29 15:45:56 inmem.go:70: try to restore from dump
[INFO] 2024/11/29 15:45:56 inmem.go:75: got from dump:
[INFO] 2024/11/29 15:45:56 inmem.go:77: -- current_humidity: 74485
[INFO] 2024/11/29 15:45:56 inmem.go:77: -- current_temperature: 74488
[INFO] 2024/11/29 15:45:56 bme280_linux_arm64.go:18: make BME280 sensor
[INFO] 2024/11/29 15:45:56 hap.go:28: make HapSrv
[INFO] 2024/11/29 15:45:56 hap.go:46: set custom PIN
[INFO] 2024/11/29 15:45:56 srv.go:81: start syncing sensor data with 5m0s sleep
[INFO] 2024/11/29 15:45:56 srv.go:93: start web server on http://localhost:80
[INFO] 2024/11/29 15:45:56 srv.go:98: start HAP server
[ERRO] 2024/11/29 16:30:57 srv.go:124: can't get sensor data: write /dev/i2c-1: remote I/O error
[ERRO] 2024/11/29 16:35:57 srv.go:124: can't get sensor data: write /dev/i2c-1: remote I/O error
[ERRO] 2024/11/29 16:40:57 srv.go:124: can't get sensor data: write /dev/i2c-1: remote I/O error
[ERRO] 2024/11/29 16:45:57 srv.go:124: can't get sensor data: write /dev/i2c-1: remote I/O error
[INFO] 2024/11/29 16:49:39 main.go:50: server shutdown...
[INFO] 2024/11/29 16:49:39 main.go:54: ctx cancel
[INFO] 2024/11/29 16:49:39 main.go:57: try make a dump to restore it next time...
[INFO] 2024/11/29 16:49:39 main.go:61: done
[INFO] 2024/11/29 16:49:39 main.go:64: buy
---
kill prev srv
Failed to kill t-hk-srv
nohup: appending output to 'nohup.out'
[INFO] 2024/11/29 17:12:11 main.go:30: 🇭🇰 revision: 829bd6d
[INFO] 2024/11/29 17:12:11 inmem.go:70: try to restore from dump
[INFO] 2024/11/29 17:12:11 inmem.go:75: got from dump:
[INFO] 2024/11/29 17:12:11 inmem.go:77: -- current_humidity: 73734
[INFO] 2024/11/29 17:12:11 inmem.go:77: -- current_temperature: 73737
[INFO] 2024/11/29 17:12:11 bme280_linux_arm64.go:18: make BME280 sensor
[INFO] 2024/11/29 17:12:11 hap.go:28: make HapSrv
[INFO] 2024/11/29 17:12:11 hap.go:46: set custom PIN
[INFO] 2024/11/29 17:12:11 srv.go:98: start HAP server
[INFO] 2024/11/29 17:12:11 srv.go:81: start syncing sensor data with 5m0s sleep
[INFO] 2024/11/29 17:12:11 srv.go:93: start web server on http://localhost:80
done
```

How to enable I2C bus on RPi device: If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS
level. Go to "Interfaceing Options" menu, to active I2C bus. Probably you will need to reboot to load i2c kernel module.
Finally you should have device like /dev/i2c-1 present in the system.

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

## References

* https://github.com/brutella/hap – HomeKit Accessory Protocol implementation in Go
* https://github.com/d2r2/go-bsbmp – BME280 sensor driver in Go

## Contributing

Bug reports, bug fixes and new features are always welcome. Please open issues and submit pull requests for any new
code.

## License

This project is licensed under the MIT License - see
the [LICENSE](https://github.com/egregors/hk/blob/main/LICENSE) file for details.
