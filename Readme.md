# HEM-GW16A prometheus exporter

## Usage

```
./HEM-GW16A-exporter --port 0.0.0.0:9101 --url http://192.168.3.6
```
* `--port` Listen Address and port
* `--url` HEM-GW16A root address
* `--username` HEM-GW16A username
* `--password` HEM-GW16A password


### Systemd sample
```
[Unit]
Description=HEM-GW16A Exporter

[Service]
Type=simple
ExecStart=/usr/local/bin/HEM-GW16A-exporter --port 0.0.0.0:9101 --url "http://192.168.3.45"
PrivateTmp=false

[Install]
WantedBy=multi-user.target
```

## License

HEM-GW16A prometheus exporter by Rimpei Kunimoto is licensed under the Apache License, Version2.0