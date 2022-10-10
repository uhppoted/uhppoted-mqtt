# HOWTO: _Provisioning AWS Greengrass_

This HOWTO outlines the steps required to create and provision the `core` and `thing` devices on _AWS Greengrass_. What you want to
end up with is:

- a `core` device that acts as the MQTT broker for _uhppoted-mqtt_ and forwards messages to/from the AWS Greengrass 
  message dispatch system.
- a `thing` device that acts defines the capabilities and permissions for _uhppoted-mqtt_ to act as an AWS Greengrass IoT 
  element.

### Background
The relevant sections in the _AWS Greengrass_ documentation are:

- [Tutorial:Getting started](https://docs.aws.amazon.com/greengrass/v2/developerguide/getting-started.html)
- [Tutorial:Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html)
- [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html)

## Provision `core` device

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

