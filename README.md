# 🇭🇰 hk – HomeKit BME280 sensor integration + web iface + notifications

[![Test](https://github.com/egregors/hk/actions/workflows/test.yml/badge.svg)](https://github.com/egregors/hk/actions/workflows/test.yml)

Smart temperature and humidity monitoring system with BME280 sensor, Raspberry Pi, Apple HomeKit integration, and intelligent notification system.

## tl;dr

In general, this project is a comprehensive IoT solution combining:

* **Hardware**: Raspberry Pi + BME280 sensor integration – to get temperature and humidity data
* **Web Interface**: Simple web server – to expose the data to local network with metrics visualization
* **HomeKit**: HAP server – to expose the data to HomeKit for use in Apple Home app
* **Monitoring**: Smart notification system – to alert you when sensor issues occur
* **Persistence**: Metrics collection with backup/restore and autosave capabilities

### Features:

* Temperature and humidity data from BME280 sensor
* Web server to expose the data
* Simple metrics collection with retention and Braille graph
* Backup and restore collected data (gob dump)
* Autosave metrics with configurable intervals
* Logging
* HomeKit integration
* Custom PIN for HomeKit
* Smart notification system (ntfy.sh support)
* Automatic error notifications for sensor failures
* USB power control for external devices (like LED garlands)

### Screenshots:
<p align="center">
     <img width="80%" alt="Screenshot 2024-11-29 at 18 29 12" src="https://github.com/user-attachments/assets/04885fc5-93d6-4e4e-8144-925465979f4d">
</p>
<div align="center">
     <img width="30%" alt="Screenshot 2024-11-29 at 18 29 12" src="https://github.com/user-attachments/assets/5cb203cb-40ca-4f28-b27c-affc23ae93a1">
     <img width="43%" alt="Screenshot 2024-11-29 at 18 29 12" src="https://github.com/user-attachments/assets/fd22263f-a434-4c1d-9fac-40a4aad7972d">
</div>

## Quick start

Actual sensor code should be compiled on Raspberry Pi device. You can use any other device to compile the code, but you
should copy the binary to Raspberry Pi device to run it.
In purpose to make development process less painful there is a `ClimateSensor` interface. So, you need to install

### Environment Variables

For production mode, you can configure notifications:

* `NOTIFY_URL` - URL for ntfy.sh notifications (optional). If set, the system will send notifications when sensor errors occur.

Example:
```bash
export NOTIFY_URL="https://ntfy.sh/your-topic-name"
```

### Build and Run

The project supports two build modes:

1. **Development mode** (`cmd/dev/main.go`):
   - Uses NoopHap (fake HomeKit server) for development
   - Uses Noop notifier (no actual notifications)
   - Shorter metrics retention (3600 hours)
   - Suitable for testing and development

2. **Production mode** (`cmd/prod/main.go`):
   - Full HomeKit integration with real HAP server
   - ntfy.sh notifications support
   - Longer metrics retention (30 days)
   - Autosave metrics every 60 minutes
   - Requires `NOTIFY_URL` environment variable for notifications

Build and deploy:

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
[INFO] 2024/11/29 16:49:39 main.go:64: bye
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

## Notification System

The project includes a smart notification system that monitors sensor health and sends alerts when issues occur.

### Features:

* **ntfy.sh integration** - Send push notifications to your devices via ntfy.sh service
* **Automatic error detection** - Monitor BME280 sensor connectivity and data integrity
* **Configurable notifications** - Easy setup with environment variables
* **Graceful fallback** - System continues to work even if notifications fail

### Setup:

1. Choose or create a topic on [ntfy.sh](https://ntfy.sh)
2. Set the `NOTIFY_URL` environment variable:
   ```bash
   export NOTIFY_URL="https://ntfy.sh/your-unique-topic-name"
   ```
3. Subscribe to the topic on your mobile device using the ntfy app
4. Run the production build - you'll receive notifications when sensor errors occur

### Notification Types:

* **Sensor Error** - Triggered when BME280 sensor fails to read temperature or humidity data
* **I/O Errors** - Hardware communication issues (e.g., "write /dev/i2c-1: remote I/O error")
* **Connection Problems** - When sensor becomes unresponsive

## USB Power Control

The project includes USB power control functionality for external devices (like LED garlands) using [uhubctl](https://github.com/mvp/uhubctl).

### Features:

* **Dynamic hub detection** - Automatically detects USB hub location on startup
* **Restart resilient** - Works correctly even after Raspberry Pi restarts when hub addresses may change
* **HomeKit integration** - Control USB power through Apple Home app
* **Smart hub selection** - Prefers hubs with per-port power switching (ppps) capability

### Setup:

1. Install uhubctl on your Raspberry Pi:
   ```bash
   sudo apt-get install uhubctl
   ```

2. Ensure your USB hub supports per-port power switching. Check with:
   ```bash
   sudo uhubctl
   ```

3. The system will automatically detect the correct hub location on startup and log it:
   ```
   [INFO] detected USB hub location: 1-1
   ```

### How It Works:

On startup, the system:
1. Runs `uhubctl` to list all available USB hubs
2. Parses the output to find controllable hubs
3. Prioritizes hubs with per-port power switching (ppps)
4. Stores the detected location for all subsequent power control operations

This ensures reliable USB power control even if hub addresses change after a restart.

## References

* https://github.com/brutella/hap – HomeKit Accessory Protocol implementation in Go
* https://github.com/d2r2/go-bsbmp – BME280 sensor driver in Go
* https://ntfy.sh – Simple HTTP-based pub-sub notification service

## Contributing

Bug reports, bug fixes and new features are always welcome. Please open issues and submit pull requests for any new
code.

## License

This project is licensed under the MIT License - see
the [LICENSE](https://github.com/egregors/hk/blob/main/LICENSE) file for details.
