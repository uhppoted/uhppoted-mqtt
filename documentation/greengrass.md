**WORK IN PROGRESS**

# HOWTO: Getting started with _uhppoted-mqtt_ with _AWS Greengrass_

This _HOWTO_ is a simplified guide to getting _uhppoted-mqtt_ up and running with _AWS Greengrass_.

Getting started with _HiveMQ_ or _Mosquito_ is relatively straightforward - _AWS Greengrass_ however is
a whole 'nuther beast in terms of complexity and getting just a base system on which to build can be 
daunting. This guide outlines the steps required to _"just get something working"_:

- It is **NOT** intended as a guide to setting up a production ready system. Among other things it uses fairly
  permissive policies and works around some recommended practices (e.g. _Greengrass Discovery_) that add complication
  when you just want to get a system running.

- It is probably egregiously wrong in places i.e. you should read the official documentation.


## Raison d'être

AWS Greengrass has a couple of expectations that make getting _uhppoted-mqtt_ configured to connect to the _Moquette_ 
MQTT broker component non-trivial until you've read a couple of reams of documentation which in turn requires more
than a little coffee:

1. By default, _Moquette_ is configured to require TLS mutual authentication i.e. clients are required to present a valid
   X.509 certificate signed by a common certificate authority during the intial TLS handshake.
2. Clients are expected to obtain the X.509 certificates required to connect to _Moquette_ using AWS _Greengrass Discovery_.
3. The workaround for clients not using Greengrass Discovery is not officially documented (ref. [Add docs about manual connection of client devices to GG Core without cloud discovery](https://github.com/awsdocs/aws-iot-greengrass-v2-developer-guide/issues/20)).
4. Then there's the message routing...

This guide is essentially a desperation resource distilled from:

- The discussion _[uhppoted-mqtt to AWS Greengrass aws.greengrass.clientdevices.mqtt.Moquette client](https://github.com/uhppoted/uhppoted/discussions/14)_
- [Tutorial: Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html)
- [Install AWS IoT Greengrass Core software with manual resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/manual-installation.html)
- [Implementing Local Client Devices with AWS IoT Greengrass](https://aws.amazon.com/blogs/iot/implementing-local-client-devices-with-aws-iot-greengrass)
- [How to Bridge Mosquitto MQTT Broker to AWS IoT](https://aws.amazon.com/blogs/iot/how-to-bridge-mosquitto-mqtt-broker-to-aws-iot/)


## Outline

For this guide, the target system will comprise a clean Ubuntu 22.04 LTS VPS with:

- an _AWS Greengrass_ _core_ device with the _Auth_, _Moquette_ and _MQTT Bridge_ components
- an _AWS Greengrass_ _thing_ for _uhppoted-mqtt_
- a daemonized _uhppoted-mqtt_

It should be similar'ish for anything else but YMMV. For the rest of this guide:

- the _core_ device will be named and referred to as _uhppoted-greengrass_. Think of it as the MQTT broker.
- the _thing_ device will be named and referred to as _uhppoted-mqtt_. Think of it as _uhppoted-mqtt_.
- both _core_ and _thing_ will be installed on the same machine. This is not a requirement but does avoid having
  to change firewall rules and NATs, etc.

## Host

The instructions below are for Ubuntu 22.04 LTS - modify as required for other systems.

1. Install Java and (optionally) Go:
```
sudo apt install openjdk-8-jdk
sudo apt install golang
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
sudo adduser  --system ggc_user
sudo addgroup --system ggc_group
```

5. Create folders:
```
sudo mkdir -p /opt/aws
sudo mkdir -p /opt/uhppoted

sudo chown -R admin:admin /opt/aws
sudo chown -R uhppoted:uhppoted /opt/uhppoted
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

In the AWS IAM console, create two policies:

1. A _uhppoted-greengrass_ policy for provisioning (a.ka. installing and configuring) the AWS Greengrass 'core' and
   'thing' devices. 
2. A _uhppoted-greengrass-cli_ policy for the AWS Greengrass CLI

For this HOWTO, the policy is based on the [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html) from the AWS Greengrass Developer Guide.

The _uhppoted-greengrass-cli_ policy is a convenience for this HOWTO and is not required if you don't anticipate needing
to use the AWS Greengrass CLI to debug/manage 'core' or 'thing' devices. Chance are you'll probably need it.

#### `uhppoted-greengrass`

1. Open the AWS [_IAM_](https://console.aws.amazon.com/iamv2) console
2. Copy the _Account ID_ from AWS Account for later
3. Open the [_Policies_](https://console.aws.amazon.com/iamv2/home#/policies) tab
4. Click on _Create policy_
5. Open the _JSON_ tab and paste the following policy, replacing the \<account-id\> with the Amazon
   account ID from step 2:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "CreateTokenExchangeRole",
            "Effect": "Allow",
            "Action": [
                "iam:AttachRolePolicy",
                "iam:CreatePolicy",
                "iam:CreateRole",
                "iam:GetPolicy",
                "iam:GetRole",
                "iam:PassRole"
            ],
            "Resource": [
                "arn:aws:iam::<account-id>:role/GreengrassV2TokenExchangeRole",
                "arn:aws:iam::<account-id>:policy/GreengrassV2TokenExchangeRoleAccess"
            ]
        },
        {
            "Sid": "CreateIoTResources",
            "Effect": "Allow",
            "Action": [
                "iot:AddThingToThingGroup",
                "iot:AttachPolicy",
                "iot:AttachThingPrincipal",
                "iot:CreateKeysAndCertificate",
                "iot:CreatePolicy",
                "iot:CreateRoleAlias",
                "iot:CreateThing",
                "iot:CreateThingGroup",
                "iot:DescribeEndpoint",
                "iot:DescribeRoleAlias",
                "iot:DescribeThingGroup",
                "iot:GetPolicy"
            ],
            "Resource": "*"
        },
        {
            "Sid": "DeployDevTools",
            "Effect": "Allow",
            "Action": [
                "greengrass:CreateDeployment",
                "iot:CancelJob",
                "iot:CreateJob",
                "iot:DeleteThingShadow",
                "iot:DescribeJob",
                "iot:DescribeThing",
                "iot:DescribeThingGroup",
                "iot:GetThingShadow",
                "iot:UpdateJob",
                "iot:UpdateThingShadow"
            ],
            "Resource": "*"
        }
    ]
}
```
6. Click on _Next: Tags_
7. Click on _Next: Review_
8. Fill in fields:
   - `Name`: `uhppoted-greengrass`
   - `Description`: Greengrass policy for deploying uhppoted-mqtt
9. Click on _Create Policy_

#### `uhppoted-greengrass-cli`

1. Open the AWS [_IAM_](https://console.aws.amazon.com/iamv2) console
2. Copy the _Account ID_ from AWS Account for later
3. Open the [_Policies_](https://console.aws.amazon.com/iamv2/home#/policies) tab
4. Click on _Create policy_
5. Open the _JSON_ tab and paste the following policy, replacing the \<account-id\> with the Amazon
   account ID from step 2:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "greengrass:List*",
                "greengrass:Get*"
            ],
            "Resource": "*"
        }
    ]
}
```

6. Click on _Next: Tags_
7. Click on _Next: Review_
8. Fill in fields:
   - `Name`: `uhppoted-greengrass-cli`
   - `Description`: Greengrass policy for the CLI
9. Click on _Create Policy_


### Groups

In the AWS IAM console, create two groups:

1. A _uhppoted-greengrass_ group for the users to be given the permissions required to provision the AWS Greengrass 'core' and
   'thing' devices. 
2. A _uhppoted-greengrass-cli_ group for the users to be given the permissions required to use the AWS Greengrass CLI (optional).

Steps:

1. Open the AWS [_IAM_](https://console.aws.amazon.com/iamv2) console
2. Open the [_User groups_](https://console.aws.amazon.com/iamv2/home#/groups) page
3. Click on _Create group_
4. Enter the group name _uhppoted-greengrass_
5. Tick the _uhppoted-greengrass_ policy created above to attach it to the group
6. Click on _Create group_

Optionally, repeat steps 3-6 to create a _uhppoted-greengrass-cli_ group with the _uhppoted-greengrass-cli_ policy attached.


### Users

_(skip this section if you're using temporary credentials)_

In the AWS IAM console, create two users:

1. A _uhppoted-greengrass_ user for provisioning the AWS Greengrass 'core' and 'thing' devices. 
2. A _uhppoted-greengrass-cli_ user for the AWS Greengrass CLI (optional).

Steps:

1. Open the AWS [_IAM_](https://console.aws.amazon.com/iamv2) console
2. Open the [_Users_](https://console.aws.amazon.com/iamv2/home#/users) page
3. Click _Add users_
4. Enter the user name _uhppoted-greengrass_
5. Select the _Access key - Programmatic access_ AWS credential type
6. Click _Next: Permissions_
7. Check the _uhppoted-greengrass_ group under the _Add user to group_ section
8. Click _Next: Tags_
9. Click _Next: Review_
10. Click _Create user_
11. Copy the access key and secret key for later use
12. Click _Close_

Optionally, repeat steps 3-12 to create a _uhppoted-greengrass-cli_ user in the _uhppoted-greengrass-cli_ group.


## AWS Greengrass

The initially relevant sections in the AWS Greengrass documentation are:

- [Tutorial:Getting started](https://docs.aws.amazon.com/greengrass/v2/developerguide/getting-started.html)
- [Tutorial:Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html)
- [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html)

1. Follow the [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html) instructions to install the _core_ device on the Ubuntu host
- use the access key and secret for the _uhppoted-greengrass_ user

2. Install the following additional components to the _core_ device:
- Auth (client device auth)
- MQTT 3.1.1 broker
- MQTT bridge 
- IPDetector

using the instructions from [Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html) instructions steps 1 and 2.

   _Notes_
   - Use the following policy for the Auth component (it is identical to the one in the AWS tutorial except for the group
and policy names):

```
{
  "deviceGroups": {
    "formatVersion": "2021-03-05",
    "definitions": {
      "UhppotedIoTGroup": {
        "selectionRule": "thingName: uhppoted-mqtt*",
        "policyName": "UhppotedIotPolicy"
      }
    },
    "policies": {
      "UhppotedIotPolicy": {
        "AllowConnect": {
          "statementDescription": "Allow client devices to connect.",
          "operations": [
            "mqtt:connect"
          ],
          "resources": [
            "*"
          ]
        },
        "AllowPublish": {
          "statementDescription": "Allow client devices to publish to all topics.",
          "operations": [
            "mqtt:publish"
          ],
          "resources": [
            "*"
          ]
        },
        "AllowSubscribe": {
          "statementDescription": "Allow client devices to subscribe to all topics.",
          "operations": [
            "mqtt:subscribe"
          ],
          "resources": [
            "*"
          ]
        }
      }
    }
  }
}
```

   - Use the following configuration for the MQTT bridge:
{
  "mqttTopicMapping": {
    "UhppoteIotMapping": {
      "topic": "uhppoted/#",
      "source": "LocalMqtt",
      "target": "IotCore"
    }
  }
}

3. Create a 'thing' device for _uhppoted-mqtt_ via the AWS console.

_TODO_

## _uhppoted-mqtt_

1. uhppoted-mqtt
2. Certificates
3. Firewall

## References

1. [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html)


