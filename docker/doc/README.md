# uhppoted-mqtt Docker container

The _docker_ folder contains the files required to create a minimal Docker image that runs a 
basic version of _uhppoted-mqtt_:

| File          | Description                        |
|---------------|------------------------------------|
| Dockferfile   |                                    |
| uhppoted-mqtt | Pre-built binary for Alpine Linux  |
| uhppoted-conf | uhppoted-mqtt configuration file   |
| _broker.pem_  | Sample MQTT broker certificate     |

## Building a Docker image

The basic _Docker_ image is configured to connect using plain TCP.

1. Update the _uhppoted.conf_ file with:
   - The address of the MQTT broker
   - The configured MQTT _client ID_, _username_ and _password_ for _uhppoted-mqtt_
   - The details of the access controllers

2. Build the _Docker_ image:
```
cd docker 
docker build -f Dockerfile -t uhppoted/mqtt .
```

3. Start the _Docker_ container:
```
docker run --name mqttd --rm uhppoted/mqtt
```

4. Once you have established that the container is configured correctly, start the container in detached mode:
```
docker run --detach --name mqttd --rm uhppoted/mqtt
```


## Upgrading to TLS

1. Copy the MQTT broker certificate to the docker folder and add it to the _Dockerfile_:
```
...
# MQTT broker TLS certificate
COPY broker.pem    /etc/uhppoted
...
```

2. Configure the broker connection with the TLS connection in _uhppoted.conf_ e.g.:
```
...
mqtt.connection.broker = tls://192.168.1.100:8883
mqtt.connection.broker.certificate = /etc/uhppoted/broker.pem
...
```

3. Rebuild the docker image:
```
cd docker 
docker build -f Dockerfile -t uhppoted/mqtt .
```


## Upgrading to TLS client authentication

If the MQTT broker requires client TLS authentication:

1. Copy the client key and certificate to the _docker/secure_ folder and add them to the Dockerfile:
```
...
# MQTT broker TLS client key and certificate
COPY secure/client.key  /etc/uhppoted/mqtt
COPY secure/client.cert /etc/uhppoted/mqtt
...
```

2. Configure the broker connection for TLS client authenticationin _uhppoted.conf_ e.g.:
```
...
mqtt.connection.client.certificate = /etc/uhppoted/mqtt/client.cert
mqtt.connection.client.key = /etc/uhppoted/mqtt/client.key
...
```

3. Rebuild the docker image:
```
cd docker 
docker build -f Dockerfile -t uhppoted/mqtt .
```


## Upgrading to a public MQTT broker

If you are using a public MQTT broker it is **highly** recommended that you enable message signing and encryption:

1. Copy the RSA signing and encryption keys to the _docker/secure/rsa_ folder add them to the Dockerfile:
```
...
# uhppoted-mqtt RSA signing and encryption keys
COPY secure/rsa/signing/mqttd.key    /etc/uhppoted/mqtt/rsa/signing
COPY secure/rsa/encryption/mqttd.key /etc/uhppoted/mqtt/rsa/encryption
...
```

2. Configure _uhppoted-mqtt_ for message signing and encryption, e.g.:
```
...
mqtt.security.authentication = RSA
mqtt.security.HMAC.required = true
mqtt.security.HMAC.key = <HMAC key>
mqtt.security.rsa.keys = /etc/uhppoted/mqtt/rsa
mqtt.security.nonce.required = true
mqtt.security.nonce.server = /var/uhppoted/mqtt.nonce
mqtt.security.nonce.clients = /var/uhppoted/mqtt.nonce.counters
mqtt.security.outgoing.sign = true
mqtt.security.outgoing.encrypt = true
...
```

3. Rebuild the docker image:
```
cd docker 
docker build -f Dockerfile -t uhppoted/mqtt .
```
