**WORK IN PROGRESS**

# HOWTO: Getting started with _uhppoted-mqtt_ with _AWS Greengrass_

This _HOWTO_ is a simplified guide to getting _uhppoted-mqtt_ up and running with _AWS Greengrass_.

Getting started with _HiveMQ_ or _Mosquito_ is relatively straightforward - _AWS Greengrass_ however is
a whole 'nuther beast in terms of complexity and getting just a base system on which to build can be 
daunting. This guide outlines the steps required to _"just get something working"_:

- It is **NOT** intended as a guide to setting up a production ready system. Among other things it uses fairly
  permissive policies and works around some recommended practices (e.g. _Greengrass Discovery_) that add complication
  when you just want to get a system running.

- It is quite probably egregiously wrong in places i.e. you should read the official documentation.


## Background

AWS Greengrass has a couple of expectations that make getting _uhppoted-mqtt_ configured to connect to the _Moquette_ 
MQTT broker component not entirely trivial - at least not until you've read a couple of reams of documentation along
with more than a reasonable amount of coffee:

1. By default, _Moquette_ is configured to require TLS mutual authentication i.e. clients are required to present a valid
   X.509 certificate signed by a common certificate authority during the intial TLS handshake. For a normal Greengrass
   installation this is a GOOD thing but is somewhat overkill when _uhppoted-mqtt_ and the _Moquette_ are on the same device
   and a reasonable set of firewall rules are in place. Howsoever...

2. Clients are expected to dynamically obtain the X.509 certificates required to connect to _Moquette_ using AWS
   _Greengrass Discovery_.

3. The workaround for clients not using Greengrass Discovery is not officially documented (ref. [Add docs about manual connection of client devices to GG Core without cloud discovery](https://github.com/awsdocs/aws-iot-greengrass-v2-developer-guide/issues/20)).

4. Then there's the message routing...

This guide is essentially a desperation resource distilled from:

- The discussion _[uhppoted-mqtt to AWS Greengrass aws.greengrass.clientdevices.mqtt.Moquette client](https://github.com/uhppoted/uhppoted/discussions/14)_
- [Tutorial: Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html)
- [Install AWS IoT Greengrass Core software with manual resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/manual-installation.html)
- [Implementing Local Client Devices with AWS IoT Greengrass](https://aws.amazon.com/blogs/iot/implementing-local-client-devices-with-aws-iot-greengrass)
- [How to Bridge Mosquitto MQTT Broker to AWS IoT](https://aws.amazon.com/blogs/iot/how-to-bridge-mosquitto-mqtt-broker-to-aws-iot/)

## Quickstart

1. Work through the [Preparation](#preparation) section below
2. Work through the [IAM](IAM.md) _HOWTO_ to setup the AWS policies, groups and users
3. Work through the [Provisioning](provisioning.md) _HOWTO_ to setup the AWS IoT core devices
4. Work through the [uhpppoted-mqtt](uhppoted-mqtt.md) _HOWTO_ to setup and configure _uhppoted-mqtt_
5. Open the [AWS IoT Core Test client](https://console.aws.amazon.com/iot/home#/test) and subscribe to `#`
6. Voilá!

The write-up below adds a bit more detail about the process.

## Outline

For this guide, the target system will comprise a clean Ubuntu 22.04 LTS VPS with:

- an _AWS Greengrass_ `core` device with the _Auth_, _Moquette_, _MQTT Bridge_ ~~and _IP Detector_~~ components
- an _AWS Greengrass_ `thing` device for _uhppoted-mqtt_
- a daemonized _uhppoted-mqtt_

It should be similar'ish for anything else but YMMV. For the rest of this guide:

- the `core` device will be named and referred to as _uhppoted-greengrass_. Think of it as the MQTT broker.
- the `thing` device will be named and referred to as _uhppoted-thing_. Think of it as _uhppoted-mqtt_ - the controllers
  themselves don't feature except as a source/destination of MQTT messages.
- both `core` and `thing` will be installed on the same machine. This is not a requirement but in the interests of keeping
  this _HOWTO_ reasonable it does avoid having to configure firewall rules and NATs, etc.

## Preparation

The instructions below are for Ubuntu 22.04 LTS - modify as required for other systems.

1. Install Java and (optionally) Go:
```
sudo apt install -y openjdk-8-jdk
sudo apt install -y golang
```

   _Note_: 

   - _Debian's default golang package is quite out of date - install the latest and greatest from [https://go.dev/doc/install](https://go.dev/doc/install)._

2. Create _admin_ user:
```
sudo adduser admin
sudo usermod -aG sudo admin
```

3. Create _uhppoted_ user:
```
sudo adduser uhppoted
```

4. Create _Greengrass_ user/group:
```
sudo adduser  --system ggc_user
sudo addgroup --system ggc_group
```

5. Create folders:
```
sudo mkdir -p /opt/aws
sudo mkdir -p /opt/uhppoted
sudo mkdir -p /opt/aws/certificates
sudo mkdir -p /etc/uhppoted/mqtt/greengrass
sudo mkdir -p /var/uhppoted

sudo chown -R admin:admin /opt/aws
sudo chown -R uhppoted:uhppoted /opt/uhppoted
sudo chown -R uhppoted:uhppoted /etc/uhppoted
sudo chown -R uhppoted:uhppoted /var/uhppoted
```

5. Update firewall rules to allow local connections to the MQTT broker (TCP, port 8883) and any UHPPOTE controllers (UDP, port 6000):
```
hostname -I
sudo ufw allow from $(hostname -I | cut -d ' ' -f 1) to any port 8883  proto tcp
sudo ufw allow from $(hostname -I | cut -d ' ' -f 1) to any port 60000 proto udp
```

## AWS IAM

The basic requirements are:

1. A Greengrass service role with the permissions required to deploy and manage Greengrass devices.
2. An IAM policy with the necessary permissions required to create, configure and run a Greengrass 'core' 
   device with a _Moquette_ MQTT broker.
3. An IAM group with the necessary policies and permissions for users needed to create and run the devices.
4. An IAM user to use for creating, configuring and running the Greengrass devices. 

More detail can be found in [HOWTO: Greengrass IAM](IAM.md) for those unfamiliar with IAM or needing more detail, but essentially you want to end up with:

1. A _Greengrass_ServiceRole_ for the AWS Greengrass service
2. A _uhppoted-greengrass_ policy for provisioning (a.ka. installing and configuring) the _AWS Greengrass_ `core` and
   `thing` devices. The AWS [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html)
   in the Greengrass developer guide is a good starting point.
3. A _uhppoted-greengrass_ group for the users to be assigned the permissions required to provision the AWS Greengrass `core` and
   `thing` devices. 
4. A _uhppoted-greengrass_ user for provisioning the AWS Greengrass `core` and `thing` devices. 


## AWS Greengrass

The next step is to provision the `core` and `thing` devices on _AWS Greengrass_. There is more detail in [HOWTO:Provisioning AWS Greengrass](provisioning.md), but essentially 
you want to end up with:

- a `core` device that acts as the MQTT broker for _uhppoted-mqtt_ and forwards messages to/from the AWS Greengrass 
  message dispatch system.
- a `thing` device that acts defines the capabilities and permissions for _uhppoted-mqtt_ to act as an AWS Greengrass IoT 
  element.

1. The `core` device should be provisioned with the following additional components:
   - Auth (client device auth)
   - MQTT 3.1.1 broker
   - MQTT bridge 
   - ~~IPDetector~~
     - Configure the endpoint address manually to the VPS IP address

   _References_:
   -  [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html)
   - [Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html)

2. The `core` device MQTT bridge should be configured to use the following topic mapping:
```
{
  "mqttTopicMapping": {
    "UhppoteIotMapping": {
      "topic": "uhppoted/#",
      "source": "LocalMqtt",
      "target": "IotCore"
    }
  }
}
```

3. Create a `thing` device for _uhppoted-mqtt_ via the AWS console.

## _uhppoted-mqtt_

Installing _uhppoted-mqtt_ is straightforward and described in the [README](https://github.com/uhppoted/uhppoted-mqtt#installation), 
but a step-by-step guide is [here](uhppoted-mqtt.md).

If you're installing it as a service/daemon (recommended) the installation will automatically create (or update) the _/etc/uhppotd/uhppoted.conf_ configuration file. If you're installing it as a console application and don't already have a _uhppoted.conf_ file you can generate one by running the following command:
```
./uhppoted-mqtt config > /etc/uhppoted/uhppoted.conf
```

The default installation is configured with full security enabled which is unnecessary for an integration with _Greengrass_ and
also makes debugging difficult. It can always be re-enabled incrementally once the system is up and running.

To run without internal security, edit the _/etc/uhppoted/uhppoted.conf_ file:
```
...
mqtt.security.HMAC.required = false
mqtt.security.authentication = NONE
mqtt.security.nonce.required = false
mqtt.security.outgoing.sign = false
mqtt.security.outgoing.encrypt = false
...
```

## Certificates

By default, _Greengrass_ expects MQTT clients to connect with mutual TLS authentication - this is a GOOD thing, leave it like
that. It does mean however that you need to provide TLS certificates for _uhppoted-mqtt_, which can be done manually or by retrieving them from the Greengrass certificate server. 

Manually provisioning the certificates is only appropriate for debugging an initial setup - for anything else it is highly 
recommneded that you use a script to retrieve then from the certificate server and to also run that script in a _cron_ job
to regularly update the certificates which can be revoked and updated from the AWS console.

### Provisioning certificates manually

You need the following certificate components (in PEM format):
- AWS Root CA certificate
- MQTT broker certificate
- MQTT client certificates
- MQTT client key

1. The AWS Root CA certificates and client certificate and key can be downloaded from the _AWS IoT_ console while creating 
the _uhppoted-mqtt_ `thing` (see [Provision a thing device for uhppoted-mqtt](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/greengrass/provisioning.md#provision-a-thing-device-for-uhppoted-mqtt)). Copy the certificates
and key to:

| File                       | Folder                                           |
|----------------------------|--------------------------------------------------|
| AmazonRootCA1.pem          | /etc/uhppoted/mqtt/greengrass/AmazonRootCA1.cert |
| AmazonRootCA3.pem          | /etc/uhppoted/mqtt/greengrass/AmazonRootCA3.cert |
| 3e7a...-private.pem.key    | /etc/uhppoted/mqtt/greengrass/thing.key          |
| 3e7a...certificate.pem.crt | /etc/uhppoted/mqtt/greengrass/thing.cert         |


(_Optionally_) Install the AWS Root CA certificate in your system trust store. This should not really be necessary 
unless _OpenSSL_ complains about not being able to verify the trust chain. The instructions for _Ubuntu_ can be found 
[here](https://ubuntu.com/server/docs/security-trust-store) but for reference:
```
sudo apt-get install -y ca-certificates
sudo cp /etc/uhppoted/mqtt/greengrass/AmazonRootCA1.cert /usr/local/share/ca-certificates/AmazonRootCA1.cert
sudo cp /etc/uhppoted/mqtt/greengrass/AmazonRootCA3.cert /usr/local/share/ca-certificates/AmazonRootCA3.cert
sudo update-ca-certificates
```

2. Copy the MQTT broker CA certificate from _/usr/local/etc/uhppoted/mqtt/greengrass_:
```
cp /greengrass/v2/work/aws.greengrass.clientdevices.Auth/CA.pem /usr/local/etc/uhppoted/mqtt/greengrass/CA.cert
```

3. Update the _uhppoted.conf_ file:
```
...
mqtt.connection.client.ID = uhppoted-mqtt
mqtt.connection.broker.certificate = /usr/local/etc/uhppoted/mqtt/greengrass/CA.cert
mqtt.connection.client.certificate = /usr/local/etc/uhppoted/mqtt/greengrass/thing.cert
mqtt.connection.client.key = /usr/local/etc/com.github.uhppoted/mqtt/greengrass/thing.key
...
```

##### Provisioning certificates from the _AWS Greengrass_ certificate server

_(this assumes you included the _IPDetector_ module in the AWS Greengrass setup)_

_tl;dr; The documentation folder contains a [sample script](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/uhppoted-setup.sh) (contributed by Tim Irwin) for provisioning the certificates from the AWS certificate server which can be customized to match your system._


## References

1. [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html)
2. [Stackoverflow:How can I make a topic/action to be allowed only to authorized users?](https://iot.stackexchange.com/questions/5640/how-can-i-make-a-topic-action-to-be-allowed-only-to-authorized-users)
3. [AWS Lambda tar.gz](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/aws-lambda-tar.py)
4. [Add docs about manual connection of client devices to GG Core without cloud discovery](https://github.com/awsdocs/aws-iot-greengrass-v2-developer-guide/issues/20)

