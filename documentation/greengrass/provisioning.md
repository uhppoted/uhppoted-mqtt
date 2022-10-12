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

From the [Install AWS IoT Greengrass Core software with automatic resource provisioning](https://docs.aws.amazon.com/greengrass/v2/developerguide/quick-installation.html) instructions to install the `core` device on the Ubuntu host, 
and using the _access key_ and _secret_ for the _uhppoted-greengrass_ user:
```
sudo su admin
export AWS_ACCESS_KEY_ID=<uhppoted-greengrass user access key>
export AWS_SECRET_ACCESS_KEY=<uhppoted-greengrass user secret key>

cd /opt/aws
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

_Notes_:
1. You may get a message that looks like:
```
   Encountered error - User: arn:aws:iam::...:user/uhppoted-greengrass is not authorized to perform: iam:GetPolicy on
   resource: policy arn:aws:iam::aws:policy/UhppotedGreengrassTokenExchangeRoleAccess because no identity-based policy
   allows...
```
   Normally that error is encountered because the _UhppotedGreengrassTokenExchangeRoleAccess_ does not exist (at least initially). 
   The _GreengrassInstaller_ will (most probably) automatically create the policy and the message can be treated as  a warning.

2. The next error you will encounter on a fresh install is:
```
   Error while trying to setup Greengrass Nucleus software.amazon.awssdk.services.iam.model.NoSuchEntityException: The role
   with name UhppotedGreengrassTokenExchangeRole cannot be found. (Service: Iam, Status Code: 404, Request ID: 9fb...
```
   This is mostly because the _Greengrass_ installer expects to find the _UhppotedGreengrassTokenExchangeRole_ but doesn't even
   though it has (most probably) just created it and whoever wrote the installer didn't do a wait-and-retry.

   The solution is to ... \<sigh\> just run the installer again i.e.:
```
sudo -E java -Droot="/greengrass/v2" -Dlog.store=FILE \
  -jar ./GreengrassInstaller/lib/Greengrass.jar \
  --aws-region us-east-1 \
  --thing-name uhppoted-greengrass \
  --thing-group-name UhppotedGreengrassGroup \
  --thing-policy-name UhppotedGreengrassThingPolicy \
  --tes-role-name UhppotedGreengrassTokenExchangeRole \
  --tes-role-alias-name UhppotedGreengrassCoreTokenExchangeRoleAlias \
  --component-default-user ggc_user:ggc_group \
  --provision true \
  --setup-system-service true
```


On successful completeion of the above you should have:
- the AWS Greengrass `core` device installed in _/greengrass/v2_ on the VPS
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
      - _name_: `uhppoted-mqtt`
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
scp AmazonRootCA1.pem          <host>:/etc/uhppoted/mqtt/greengrass/
scp AmazonRootCA3.pem          <host>:/etc/uhppoted/mqtt/greengrass/
scp 3e7a...-private.pem.key    <host>:/etc/uhppoted/mqtt/greengrass/
scp 3e7a...-public.pem.key     <host>:/etc/uhppoted/mqtt/greengrass/
scp 3e7a...certificate.pem.crt <host>:/etc/uhppoted/mqtt/greengrass/
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
       - Associate _uhppoted-mqtt_ (creating it if it wasn't created above)
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
        "selectionRule": "thingName: uhppoted-mqtt*",
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
        },
        {
            "Effect": "Allow",
            "Action": [
                "greengrass:Discover"
            ],
            "Resource": [
                "arn:aws:iot:region:account-id:thing/*"
            ]
        }
      }
    }
  }
}
```

9. Choose _MQTT 3.1.1 broker (Moquette)_ and leave 'as is'
10. Choose _MQTT Bridge_ and update the configuraiton with the following topic mapping:
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
12. Review and deploy

#### Check basic connectivity and certificate chain

```
openssl s_client -connect localhost:8883 -showcerts

openssl s_client \
        -CAfile /etc/uhppoted/mqtt/greengrass/CA.pem \
        -cert /etc/uhppoted/mqtt/greengrass/client.cert \
        -key /etc/uhppoted/mqtt/greengrass/client.key \
        -connect localhost:8883 

openssl s_client \
        -CAfile /etc/uhppoted/mqtt/greengrass/CA.pem \
        -cert /etc/uhppoted/mqtt/greengrass/client.cert \
        -key /etc/uhppoted/mqtt/greengrass/client.key \
        -connect localhost:8883 -showcerts
```

#### ~~Quick and dirty test~~

From [Test client device communications](https://docs.aws.amazon.com/greengrass/v2/developerguide/test-client-device-communications.html?icmpid=docs_gg_console)

```
cd /opt/aws
git clone https://github.com/aws/aws-iot-device-sdk-python-v2.git
python3 -m pip install --user ./aws-iot-device-sdk-python-v2
cd aws-iot-device-sdk-python-v2/samples

python3 basic_discovery.py --thing_name uhppoted-mqtt \
  --topic 'uhppoted/events' \
  --message 'woot!'  \
  --ca_file CA.pem   \
  --cert thing.cert  \
  --key thing.key    \
  --region us-east-1 \
  --verbosity Info

python3 basic_discovery.py   \
  --thing_name uhppoted-mqtt \
  --topic 'uhppoted/events'  \
  --message 'woot!'          \
  --ca_file /etc/uhppoted/mqtt/greengrass/AmazonRootCA.pem \
  --cert /etc/uhppoted/mqtt/greengrass/thing.cert \
  --key  /etc/uhppoted/mqtt/greengrass/thing.key  \
  --region us-east-1 \
  --verbosity Debug
```

```
python3 basic_connect.py \
  --endpoint 127.0.0.1:8883 \
  --ca_file CA.pem \
  --cert thing.cert \
  --key  thing.key

python3 basic_connect.py \
  --endpoint 127.0.0.1:8883 \
  --ca_file /etc/uhppoted/mqtt/greengrass/AmazonRootCA.pem \
  --cert /etc/uhppoted/mqtt/greengrass/thing.cert \
  --key /etc/uhppoted/mqtt/greengrass/thing.key

```

Ref. https://docs.aws.amazon.com/greengrass/v2/developerguide/troubleshooting-client-devices.html


### TODO

```
2022-10-11T20:00:12.213Z [WARN] (pool-1-thread-2) com.aws.greengrass.detector.uploader.ConnectivityUpdater: Failed to upload the IP addresses.. {}
software.amazon.awssdk.services.greengrassv2data.model.UnauthorizedException: Greengrass is not authorized to assume the Service Role associated with this account. (Service: GreengrassV2Data, Status Code: 401, 
Request ID: 7ef21c42-cd28-78ac-3fac-caa7bc792a2e, Extended Request ID: null)
        at software.amazon.awssdk.core.internal.http.CombinedResponseHandler.handleErrorResponse(CombinedResponseHandler.java:123)
        at software.amazon.awssdk.core.internal.http.CombinedResponseHandler.handleResponse(CombinedResponseHandler.java:79)
```

- [green-grass-is-not-authorized-to-assume-the-service-role](https://repost.aws/questions/QUrO84DbX-QLe8I2fiLKEshg/green-grass-is-not-authorized-to-assume-the-service-role)
- [troubleshoot-assume-service-role]( https://docs.aws.amazon.com/greengrass/v1/developerguide/security_iam_troubleshoot.html#troubleshoot-assume-service-role)
- [security_iam_troubleshoot.html#troubleshoot-assume-service-role](https://docs.aws.amazon.com/greengrass/v1/developerguide/service-role.html#manage-service-role-console)
- [greengrass-service-role](https://github.com/awsdocs/aws-iot-greengrass-v2-developer-guide/blob/main/doc_source/greengrass-service-role.md)
- [greengrass-discovery-demo-application-is-not-working](https://stackoverflow.com/questions/49610000/greengrass-discovery-demo-application-is-not-working)