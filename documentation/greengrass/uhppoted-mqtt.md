**WORK IN PROGRESS**

## HOWTO: Install and configure _uhppoted-mqtt_

### Download and install _uhppoted-mqtt_
```
sudo su uhppoted-mqtt
cd /opt/uhppoted

curl -OL https://github.com/uhppoted/uhppoted-mqtt/releases/download/v0.8.1/uhppoted-mqtt_v0.8.1.tar.gz
tar xvzf uhppoted-mqtt_v0.8.1.tar.gz

mkdir uhppoted-mqtt
cd uhppoted-mqtt
ln -s /opt/uhppoted/uhppoted-mqtt_v0.8.1/linux/uhppoted-mqtt uhppoted-mqtt
./uhppoted-mqtt config > /etc/uhppoted/uhppoted.conf
```

### Update _uhppoted-mqtt_ configuration for _AWS Greengrass_

Edit _/etc/uhppoted/uhppoted.conf_:
```
...
mqtt.connection.client.ID = uhppoted-mqtt
mqtt.connection.broker = tls://127.0.0.1:8883
mqtt.connection.client.ID = uhppoted-mqttd
mqtt.connection.broker.certificate = /etc/uhppoted/mqtt/greengrass/CA.pem
mqtt.connection.client.certificate = /etc/uhppoted/mqtt/greengrass/thing.cert
mqtt.connection.client.key = /etc/uhppoted/greengrass/thing.key

mqtt.security.HMAC.required = false
mqtt.security.authentication = NONE
mqtt.security.nonce.required = false
mqtt.security.outgoing.sign = false
mqtt.security.outgoing.encrypt = false
...
```

### Update firewall rules

Add MQTT to _ufw_:
```
sudo ufw allow from 127.0.0.1 to any port 1883  proto tcp
sudo ufw allow from 127.0.0.1 to any port 8883  proto tcp
sudo ufw allow from 127.0.0.1 to any port 60000 proto udp
```

### Run in console mode

Run _uhppoted-mqtt_ in console mode:
```
./uhppoted-mqtt --debug console
```


### TODO

```
2022/10/11 20:39:24 START
2022/10/11 20:39:24 ERROR: open /etc/uhppoted/mqtt/greengrass/CA.pem : no such file or directory
2022/10/11 20:39:24 ERROR: open /etc/uhppoted/greengrass/thing.key: no such file or directory
2022/10/11 20:39:24 WARN  open /etc/uhppoted/mqtt/cards: no such file or directory
2022/10/11 20:39:24 INFO  mqttd        Listening on 155.138.131.33:60001
2022/10/11 20:39:24 INFO  mqttd        Publishing events to uhppoted/gateway/events
2022/10/11 20:39:24 INFO  listen       Initialising event listener
2022/10/11 20:39:24 INFO  listen       Listening
2022/10/11 20:39:24 [client]   x509: cannot validate certificate for 127.0.0.1 because it doesn't contain any IP SANs
2022/10/11 20:39:24 [client]   failed to connect to broker, trying next

```

https://stackoverflow.com/questions/71292261/golang-x509-cannot-validate-certificate

