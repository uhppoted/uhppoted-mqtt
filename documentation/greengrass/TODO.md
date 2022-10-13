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