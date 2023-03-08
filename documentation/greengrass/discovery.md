# Greengrass Discovery

Enabling the Greengrass _IP Detector_ components allows the _uhppoted-mqtt_ broker certificates and endpoints to be 
retrieved from AWS rather than configured manually. 

Tim Irwin has kindly provided a [script](uhppoted-setup.sh) that does this (and a bit more) - the write up below
describes what is involved.

As per the AWS [Use IP detector to automatically manage endpoints](https://docs.aws.amazon.com/greengrass/v2/developerguide/manage-core-device-endpoints.html#use-ip-detector) documentation, _IP Detector_ is suitable for 'simple' network topologies.


## Enable IP Detector

In the [_AWS IoT_ console](https://console.aws.amazon.com/iot/home), configure the `uhppoted-greengrass` core device
with the _IP Detector_ component:

   1. Open _Manage/Greengrass devices/Core devices_ and click on the _uhppoted-greengrass_ device.
   2. Open the _Client devices tab
   3. Click on _Configure Cloud discovery_
   4. If you've worked through the initial setup everything should already be configure - except for 
      the IP Detector.
   5. Check the _IP Detector_ component and click on the _Edit configuration_
   6. Update the configuration to include the local IPv4 addresses (since presumable _uhppoted-mqtt_ is running on
      the same host as the Greengrass core device)
```
{
    "includeIPv4LoopbackAddrs": true
}
```
   7. Click _Confirm_ and then _Review and deploy_ and wait for the deployment to complete.
   8. Confirm that _IP Detector_ is running on the host:
```
sudo tail -f /greengrass/v2/logs/greengrass.log 
```
You should see log entries for the _IP Detector_:
```
2023-03-08T18:51:48.449Z [INFO] (pool-1-thread-3) com.aws.greengrass.detector.IpDetectorManager: Acquired host IP addresses. {IpAddresses=[/155.138.156.243, /127.0.0.1]}
2023-03-08T18:52:48.449Z [INFO] (pool-1-thread-2) com.aws.greengrass.detector.IpDetectorManager: Acquired host IP addresses. {IpAddresses=[/155.138.156.243, /127.0.0.1]}
2023-03-08T18:53:48.450Z [INFO] (pool-1-thread-4) com.aws.greengrass.detector.IpDetectorManager: Acquired host IP addresses. {IpAddresses=[/155.138.156.243, /127.0.0.1]}
...
```

## Retrieve broker certificate and endpoints from AWS

The certificate and endpoints for uhppoted-greengrass_ core device can be retrieved via the [_AWS Greengrass 
discovery_ REST API](https://docs.aws.amazon.com/greengrass/v1/developerguide/gg-discover-api.html):

```
curl --cert "/etc/uhppoted/mqtt/greengrass/thing.cert" --key "/etc/uhppoted/mqtt/greengrass/thing.key" https://greengrass-ats.iot.us-east-1.amazonaws.com:8443/greengrass/discover/thing/uhppoted-thing
```
(or using `jq`):
```
curl --cert "/etc/uhppoted/mqtt/greengrass/thing.cert" --key "/etc/uhppoted/mqtt/greengrass/thing.key" https://greengrass-ats.iot.us-east-1.amazonaws.com:8443/greengrass/discover/thing/uhppoted-thing | jq .
```
```
{
  "GGGroups": [
    {
      "GGGroupId": "greengrassV2-coreDevice-uhppoted-greengrass",
      "Cores": [
        {
          "thingArn": "arn:aws:iot:us-east-1:XXXXXXXXXXXX:thing/uhppoted-greengrass",
          "Connectivity": [
            {
              "Id": "155.138.156.243",
              "HostAddress": "155.138.156.243",
              "PortNumber": 8883,
              "Metadata": ""
            },
            {
              "Id": "127.0.0.1",
              "HostAddress": "127.0.0.1",
              "PortNumber": 8883,
              "Metadata": ""
            }
          ]
        }
      ],
      "CAs": [
        "-----BEGIN CERTIFICATE-----\nMIID1DCCArygAwIBAgIUEZB3p/+7ljT0z76YkjChnF5Wm8UwDQYJKoZIhvcNAQEL\nBQAwgYkxCzAJBgNVBAYTAlVTMRgwF
        gYDVQQKDA9BbWF6b24uY29tIEluYy4xHDAa\nBgNVBAsME0FtYXpvbiBXZWIgU2VydmljZXMxEzARBgNVBAgMCldhc2hpbmd0b24x\nEDAOBgNVBAcMB1NlYXR0bG
        ....
        ....
        dkYZ5yV\n-----END CERTIFICATE-----\n"
      ]
    }
  ]
}

```

The [script](uhppoted-setup.sh) shows how to extract this information and update the certificates and endpoints
for _uhppoted-mqtt_.



