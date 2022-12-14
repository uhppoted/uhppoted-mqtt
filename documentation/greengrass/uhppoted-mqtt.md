**WORK IN PROGRESS**

## HOWTO: Install and configure _uhppoted-mqtt_

### Download and install _uhppoted-mqtt_
```
sudo su uhppoted
```
```
cd /opt/uhppoted

curl -OL https://github.com/uhppoted/uhppoted-mqtt/releases/download/v0.8.3/uhppoted-mqtt_v0.8.3.tar.gz
tar xvzf uhppoted-mqtt_v0.8.3.tar.gz

mkdir uhppoted-mqtt
cd uhppoted-mqtt
ln -s /opt/uhppoted/uhppoted-mqtt_v0.8.3/linux/uhppoted-mqtt uhppoted-mqtt
./uhppoted-mqtt config > /etc/uhppoted/uhppoted.conf
```

**IF** you're using the _development_ version:
```
cd /opt/uhppoted

git clone https://github.com/uhppoted/uhppoted-mqtt
cd uhppoted-mqtt
make build
ln -s /opt/uhppoted/uhppoted-mqtt/bin/uhppoted-mqtt uhppoted-mqtt
./uhppoted-mqtt config > /etc/uhppoted/uhppoted.conf
```

### Update _uhppoted-mqtt_ configuration for _AWS Greengrass_

Edit _/etc/uhppoted/uhppoted.conf_:

- Set the MQTT broker connection information **using the host IP address for the MQTT broker***
```
mqtt.connection.client.ID = uhppoted-thing
mqtt.connection.broker = tls://<host-ip-address>:8883
mqtt.connection.broker.certificate = /etc/uhppoted/mqtt/greengrass/CA.cert
mqtt.connection.client.certificate = /etc/uhppoted/mqtt/greengrass/thing.cert
mqtt.connection.client.key = /etc/uhppoted/mqtt/greengrass/thing.key
; mqtt.connection.verify = allow-insecure
```

- Disable security (and add back in as/when required)
```
mqtt.security.HMAC.required = false
mqtt.security.authentication = NONE
mqtt.security.nonce.required = false
mqtt.security.outgoing.sign = false
mqtt.security.outgoing.encrypt = false
```

- Configure any UHPPOTE controllers in the # DEVICES section

### Run in console mode

Run _uhppoted-mqtt_ in console mode:
```
./uhppoted-mqtt run --debug console
```
The console should print a running log that looks something like this:
```
2022/12/12 20:02:39 uhppoted-mqtt service v0.8.2 - Linux (PID 36834)
2022/12/12 20:02:39 WARN  open /etc/uhppoted/mqtt/rsa/signing/mqttd.key: no such file or directory
2022/12/12 20:02:39 WARN  open /etc/uhppoted/mqtt/rsa/encryption/mqttd.key: no such file or directory
2022/12/12 20:02:39 WARN  stat /etc/uhppoted/mqtt/rsa/signing: no such file or directory
2022/12/12 20:02:39 WARN  stat /etc/uhppoted/mqtt/rsa/encryption: no such file or directory
 ... listening
 ... request
 ...          00000000  17 94 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
 ...          00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
 ...          00000020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
 ...          00000030  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
 ...

...
```


