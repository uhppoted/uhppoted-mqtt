## TODO

- [ ] Set ARN for service role
- [ ] See if there is anything helpful in https://iot.awsworkshops.com/aws-greengrassv2/lab35-greengrassv2-basics/
- [ ] Check service role permissions?? 
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


#### May be appplicable with IP Detect

_Update the _UhppotedGreengrassCoreTokenExchangeRole_ alias_

_ In the [AWS IoT console](https://console.aws.amazon.com/iot/home), edit the TokenExchangeRole created by the installer and either:_
- _set it to alias the Greengrass\_ServiceRole_
- _in IAM, create an UhppotedGreengrassTokenExchangeRole with the necessary permissions and set the alias to use the newly created role._


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

Ref. [Troubleshooting client devices](https://docs.aws.amazon.com/greengrass/v2/developerguide/troubleshooting-client-devices.html)


### TODO

```
...
2022/10/11 20:39:24 INFO  mqttd        Listening on 155.138.131.33:60001
2022/10/11 20:39:24 INFO  mqttd        Publishing events to uhppoted/gateway/events
2022/10/11 20:39:24 INFO  listen       Initialising event listener
2022/10/11 20:39:24 INFO  listen       Listening
2022/10/11 20:39:24 [client]   x509: cannot validate certificate for 127.0.0.1 because it doesn't contain any IP SANs
2022/10/11 20:39:24 [client]   failed to connect to broker, trying next
```

https://stackoverflow.com/questions/71292261/golang-x509-cannot-validate-certificate


