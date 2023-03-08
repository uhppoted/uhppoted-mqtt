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

## Provision a `core` device

Based on the [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html) instructions. 

To install the `core` device on the Ubuntu host using the _access key_ and _secret_ for the _uhppoted-greengrass_ user:
```
sudo su admin
```
```
export AWS_ACCESS_KEY_ID=<uhppoted-greengrass user access key>
```
```
export AWS_SECRET_ACCESS_KEY=<uhppoted-greengrass user secret key>
```
```
cd /opt/aws
curl -s https://d2s8p88vqu9w66.cloudfront.net/releases/greengrass-nucleus-latest.zip > greengrass-nucleus-latest.zip
unzip greengrass-nucleus-latest.zip -d GreengrassInstaller && rm greengrass-nucleus-latest.zip
java -jar ./GreengrassInstaller/lib/Greengrass.jar --version

sudo -E java -Droot="/greengrass/v2" -Dlog.store=FILE \
  -jar ./GreengrassInstaller/lib/Greengrass.jar \
  --aws-region us-east-1 \
  --thing-name uhppoted-greengrass \
  --thing-policy-name UhppotedGreengrassThingPolicy \
  --tes-role-name UhppotedGreengrassTokenExchangeRole \
  --tes-role-alias-name UhppotedGreengrassCoreTokenExchangeRoleAlias \
  --component-default-user ggc_user:ggc_group \
  --provision true \
  --setup-system-service true \
  --deploy-dev-tools true
```

_Notes_:

1. Unless you've created the _UhppotedGreengrassTokenExchangeRoleAccess_ policy and the _UhppotedGreengrassTokenExchangeRole_
   you may get a message that looks like:
```
   Encountered error - User: arn:aws:iam::...:user/uhppoted-greengrass is not authorized to perform: iam:GetPolicy
   on resource: policy arn:aws:iam::aws:policy/UhppotedGreengrassTokenExchangeRoleAccess because no identity-based
   policy allows...
```
   Normally that error is encountered because the _UhppotedGreengrassTokenExchangeRoleAccess_ does not exist (at least initially). 
   The _Greengrass Installer_ will (most probably) automatically create the policy and the message can be treated as  a warning.

2. The next error you may encounter on a fresh install is:
```
   Error while trying to setup Greengrass Nucleus software.amazon.awssdk.services.iam.model.NoSuchEntityException:
   The role with name UhppotedGreengrassTokenExchangeRole cannot be found. (Service: Iam, Status Code: 404, 
   Request ID: 9fb...
```
   This is mostly because the _Greengrass_ installer expects to find the _UhppotedGreengrassTokenExchangeRole_ but 
   doesn't even though it has (presumably) just created it and whoever wrote the installer didn't do a wait-and-retry.

   The solution is to ... _\<sigh\>_ just run the installer again i.e.:
```
sudo -E java -Droot="/greengrass/v2" -Dlog.store=FILE \
  -jar ./GreengrassInstaller/lib/Greengrass.jar \
  --aws-region us-east-1 \
  --thing-name uhppoted-greengrass \
  --thing-policy-name UhppotedGreengrassThingPolicy \
  --tes-role-name UhppotedGreengrassTokenExchangeRole \
  --tes-role-alias-name UhppotedGreengrassCoreTokenExchangeRoleAlias \
  --component-default-user ggc_user:ggc_group \
  --provision true \
  --setup-system-service true \
  --deploy-dev-tools true
```


On successful completion of the above you should have:
- the _AWS Greengrass_ `core` device installed in _/greengrass/v2_ on the VPS
- a _uhppoted-greengrass_ `core` device listed in the [_AWS IoT_ console](https://console.aws.amazon.com/iot/home) under
  _Manage/Greengrass devices/Core devices_

## Provision a `thing` device for _uhppoted-mqtt_

Based on [Tutorial: Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html). 

It's probably easier to just create the `thing` when you configure the `core` device, because that takes care of associating
the `thing` with the `core`, but in the interests of doing things the difficult way:

In the [_AWS IoT_ console](https://console.aws.amazon.com/iot/home), create a new `thing` to represent _uhppoted-mqtt_:

   1. Open _Manage/All devices/Things_ and click _Create things_
   2. Choose _Create single thing_
   3. Create a `thing` with 
      - _name_: `uhppoted-thing`
      - _device shadow_: `No shadow`
   4. Choose _Auto-generate a new certificate_
   5. Attach the _UhppotedGreengrassThingPolicy_ policy
   6. Create `thing` and download certificate and key files:
      - Device certificate
      - Public key file
      - Private key file
      - Amazon Root CA certificates
   7. Copy the certificates to the _VPS_ (or _Raspberry Pi_, etc) e.g.:
```
sudo mkdir -p /etc/uhppoted/mqtt/greengrass
```
```
scp AmazonRootCA1.pem           <host>:/etc/uhppoted/mqtt/greengrass/AmazonRootCA1.cert
scp AmazonRootCA3.pem           <host>:/etc/uhppoted/mqtt/greengrass/AmazonRootCA3.cert
scp 3e7a...-private.pem.key     <host>:/etc/uhppoted/mqtt/greengrass/thing.key
scp 3e7a...-certificate.pem.crt <host>:/etc/uhppoted/mqtt/greengrass/thing.cert
```
```
sudo chown uhppoted:uhppoted /etc/uhppoted/mqtt/greengrass/*
```

You should now have two `things` in the [_AWS IoT_ console](https://console.aws.amazon.com/iot/home):
- _uhppoted-greengrass_ (`core` device, MQTT broker)
- _uhppoted-mqtt_ (`thing` device, MQTT client)

## Configure the `core` device

The next step is to install the following additional components to the _uhppoted-greengrass_ `core` device:
- Auth (client device auth)
- MQTT 3.1.1 broker
- MQTT bridge 
- ~~IPDetector~~

Based on instructions from [Interact with local IoT devices over MQTT](https://docs.aws.amazon.com/greengrass/v2/developerguide/client-devices-tutorial.html) instructions steps 1 and 2:

   1. Open _Manage/Greengrass devices/Core devices_
   2. Open _uhppoted-greengrass_ _Client devices_ tab
   3. Choose _Configure Cloud discovery configuration_
   4. _Step 1: Select target core devices_:
       - Target type: `Core device`
       - Target name: `uhppoted-greengrass`
   5. _Step 2: Associate client devices_:
       - Associate _uhppoted-thing_ 
   6. _Step 3: Configure and deploy Greengrass components_:
   7. Choose _Greengrass nucleus_ and leave '_as is_'
   8. Choose _Client device auth_
       - Use the policy below (it is identical to the one in the AWS tutorial except for the group and policy names):

```
{
  "deviceGroups": {
    "formatVersion": "2021-03-05",
    "definitions": {
      "UhppotedGreengrassGroup": {
        "selectionRule": "thingName: uhppoted-thing*",
        "policyName": "UhppotedGreengrassThingPolicy"
      }
    },
    "policies": {
      "UhppotedGreengrassThingPolicy": {
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

9. Choose _MQTT 3.1.1 broker (Moquette)_ and leave 'as is'
10. Choose _MQTT Bridge_ and update the configuration with the following topic mapping:
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


11. ~~Tick _IP Detector_ and leave 'as is'~~

**For the moment**

- Uncheck _IP Detector_ and add a _managed endpoint_ with the IP address of the machine the`core` device is running on. 


13. Review and deploy

Review and deploy and wait for the deployment to complete.

14. Copy broker CA certificate

Since we're not using _IP Detect_, copy the _Moquette_ broker certificate:
```
sudo cp /greengrass/v2/work/aws.greengrass.clientdevices.Auth/ca.pem /etc/uhppoted/mqtt/greengrass/CA.cert
sudo chown uhppoted:uhppoted /etc/uhppoted/mqtt/greengrass/CA.cert
```

15. Check basic connectivity and certificate chain

- Without client authentication:
```
openssl s_client -connect localhost:8883 -showcerts
```

- With client authentication, using the Amazon Root CA certificate:
```
openssl s_client -CAfile /etc/uhppoted/mqtt/greengrass/AmazonRootCA1.cert \
                 -cert   /etc/uhppoted/mqtt/greengrass/thing.cert \
                 -key    /etc/uhppoted/mqtt/greengrass/thing.key \
                 -connect localhost:8883 -showcerts
```

- With client authentication using the local broker CA certificate:
```
openssl s_client -CAfile /etc/uhppoted/mqtt/greengrass/CA.cert \
                 -cert   /etc/uhppoted/mqtt/greengrass/thing.cert \
                 -key    /etc/uhppoted/mqtt/greengrass/thing.key \
                 -connect localhost:8883 -showcerts
```

