# HOWTO: _uhppoted-mqtt_ with _AWS Greengrass_: Getting Started

This _HOWTO_ is a simplified guide to getting _uhppoted-mqtt_ up and running with _AWS Greengrass_.

Getting started with _HiveMQ_ or _Mosquito_ is relatively straightforward - AWS Greengrass is a whole 'nuther
beast in terms of complexity and getting just a base system on which to build can be daunting. 

This guide outlines the steps required to _"just get something working"_:

- It is **NOT** intended as a guide to setting up a production ready system. Among other things it uses fairly
  permissive policies and works around some recommended practices (e.g. Greengrass Discovery) that add complication
  when you just want to get a system running.

- It is probably egregiously wrong in places i.e. you should read the official documentation.


## Raison d'être

AWS Greengrass has a couple of expectations that make getting _uhppoted-mqtt_ configured to connect to the _Moquette_ 
MQTT broker component non-trivial until you've read a couple of reams of documentation, accompanied by quite a lot
of coffee:

1. By default, _Moquette_ is configured to require TLS mutual authentication i.e. clients are required to present a valid
   X.509 certificate signed by a common certificate authority during the intial TLS handshake.
2. Clients are expected to obtain the X.509 certificates required to connect to _Moquette_ using AWS Greengrass Discovery.
3. The workaround for clients not using Greengrass Discovery is not officially documented (ref. [Add docs about manual connection of client devices to GG Core without cloud discovery](https://github.com/awsdocs/aws-iot-greengrass-v2-developer-guide/issues/20)).
4. Then there's the message routing...

This guide is essentially a desperation resource distilled from:

- The discussion '[uhppoted-mqtt to AWS Greengrass aws.greengrass.clientdevices.mqtt.Moquette client](https://github.com/uhppoted/uhppoted/discussions/14)'
- [Tutorial: Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html)
- [Install AWS IoT Greengrass Core software with manual resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/manual-installation.html)
- [Implementing Local Client Devices with AWS IoT Greengrass](https://aws.amazon.com/blogs/iot/implementing-local-client-devices-with-aws-iot-greengrass)
- [How to Bridge Mosquitto MQTT Broker to AWS IoT](https://aws.amazon.com/blogs/iot/how-to-bridge-mosquitto-mqtt-broker-to-aws-iot/)


## Outline

For this guide, the target system will comprise a clean Ubuntu 22.04 LTS VPS with:
- an _AWS Greengrass_ _core_ device with the _Auth_, _Moquette_ and _MQTT Bridge_ components
- an _AWS Greengrass_ _thing_ for _uhppoted-mqtt_

It should be similar'ish for anything else but YMMV.

For the rest of this guide:

- the _core_ device will be named and referred to as _uhppoted-greengrass_. Think of it as the MQTT broker.
- the _thing_ device will be named and referred to as _uhppoted-mqtt_. Think of it as _uhppoted-mqtt_.
- both _core_ and _thing_ will be installed on the same machine. This is not a requirement but does avoid having
  to change firewall rules and NATs, etc.

## Host

The instructions below are for Ubuntu 22.04 LTS - modify as required for other systems.

1. Install Java:
```
sudo apt install openjdk-8-jdk
```

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
sudo adduser --system ggc_user
sudo addgroup --system ggc_group
```

5. Create folders:
```
sudo mkdir -p /opt/aws
sudo mkdir -p /opt/uhppoted

sudo chown -R admin:admin /opt/aws
sudo chown -R uhppoted:uhppoted /opt/uhppoted
```

6. Optionally, install Go:
```
sudo apt install golang
```

## AWS IAM

_(skip this section if you're comfortable with AWS IAM and/or have already configured your users and policies)_

This section outlines the steps required to create an IAM user with the necessary permissions. You can tailor it to
fit your needs but basically it:

1. Creates an IAM policy with the necessary permissions required to create, configure and run a Greengrass 'core' 
   device with a _Moquette_ MQTT broker.
2. Create an IAM group with the necessary policies and permissions for users needed to create and run the devices.
3. Create an IAM user to use for creating, configuring and running the Greengrass devices. 

For simplicity it creates a permanent user which can/should be deleted when no longer required. If you're familiar
with creating and using temporary credentials, rather use those.

### Policies

### Groups

### Users

## AWS Greengrass

## _uhppoted-mqtt_

1. uhppoted-mqtt
2. Certificates
3. Firewall

## References

1. [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html)


