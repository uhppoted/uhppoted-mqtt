**WORK IN PROGRESS**

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

1. From the [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html) instructions to install the `core` device on the Ubuntu host, 
and using the access key and secret for the _uhppoted-greengrass_ user:
```
sudo su admin
export AWS_ACCESS_KEY_ID=<uhppoted-greengrass user access key>
export AWS_SECRET_ACCESS_KEY=<uhppoted-greengrass user secret key>

cd /opt/admin
curl -s https://d2s8p88vqu9w66.cloudfront.net/releases/greengrass-nucleus-latest.zip > greengrass-nucleus-latest.zip
unzip greengrass-nucleus-latest.zip -d GreengrassInstaller && rm greengrass-nucleus-latest.zip
java -jar ./GreengrassInstaller/lib/Greengrass.jar --version

sudo -E java -Droot="/greengrass/v2" -Dlog.store=FILE \
  -jar ./GreengrassInstaller/lib/Greengrass.jar \
  --aws-region us-east-1 \
  --thing-name uhppoted-greengrass \
  --thing-group-name uhppoted-greengrass \
  --thing-policy-name UhppotedGreengrassThingPolicy \
  --tes-role-name UhppotedGreengrassTokenExchangeRole \
  --tes-role-alias-name UhppotedGreengrassCoreTokenExchangeRoleAlias \
  --component-default-user ggc_user:ggc_group \
  --provision true \
  --setup-system-service true

```

On successful completeion of the above you should have:
- the AWS Greengrass `core` device installed in _/greengrass/v2_ on the VPS
- a `core` device listed in the _AWS IoT_ console

4. Create a 'thing' device for _uhppoted-mqtt_ via the AWS console.

Based on [Tutorial: Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html).

In the _AWS IoT_ console, create a new `thing` to represent _uhppoted-mqtt_:

   4.1 Open _Manage/Things_ and click _Create Things_
   4.2 Choose _Create single thing_
   4.3 Create a `thing` with 
       - Name: _uhppoted-mqtt_
       - Device shadow: _No shadow_
  4.4 Choose _Auto-generate a new certificate_
  4.5 Attach the _UhppotedGreengrassThingPolicy_ policy
  4.6 Create `thing` and download certificate and key files:
   - Device certificate
   - Public key file
   - Private key file
   - Amazon Root CA certificates
   4.7. Copy the certificates to the VPS:
        - /etc/uhppoted/mqtt/aws/...

And ... you should now have two `things`:
- _uhppoted-greengrass_ (core device/MQTT broker)
- _uhppoted-mqtt_ (thing device/MQTT client)


3. Configure the `core` devices

The next step is to install the following additional components to the _uhppoted-greengrass_ device:
- Auth (client device auth)
- MQTT 3.1.1 broker
- MQTT bridge 
- IPDetector

Based on instructions from [Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html) instructions steps 1 and 2:

   3.1 Open _Manage/Greengrass devices/Core devices_
   3.2 Open _uhppoted-greengrass_ _Client devices_ tab
   3.3 Choose _Configure Cloud discovery configuration_:
       - Target type: Core device
       - Target name: uhppoted-greengrass
   3.4 Associate client devices:
       - _uhppoted-mqtt_

   3.5 Greengrass nucleus
       - Leave 'as is'

   3.6 Client device auth
       - Use the policy below (it is identical to the one in the AWS tutorial except for the group and policy names):

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

3.7 MQTT 3.1.1 broker (Moquette) 
    - Leave 'as is'

3.8 MQTT Bridge
    - Use the following configuration for the MQTT bridge:
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

3.9 IP Detector
    - Leave 'as is'

3.10 Review and deploy

On the VPS do a quick and dirty test:
```
python3 basic_discovery.py \
  --thing_name uhppoted-mqtt \
  --topic 'uhppoted/events' \
  --message 'Woot!' \
  --ca_file /opt/aws/AmazonRootCA1.pem \
  --cert /opt/aws/device.pem.crt \
  --key /opt/aws/private.pem.key \
  --region us-east-1 \
  --verbosity Warn
```
