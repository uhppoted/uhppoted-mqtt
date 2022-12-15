## HOWTO: _AWS Greengrass IAM_

This HOWTO outlines the steps required to create the IAM policies, groups and users with the necessary permissions to setup
and run a system based on _uhppoted-mqtt_. You can tailor it to fit your requirements but basically what you're aiming to create
are:

1. An _AWS Greengrass_ service role with the necessary permissions to provision and manage the Greengrass devices.
2. An IAM policy with the necessary permissions required to create, configure and run a Greengrass `core` 
   device with a _Moquette_ MQTT broker.
3. An IAM group with the necessary policies and permissions for users needed to create and run the devices.
4. An IAM user to use for creating, configuring and running the Greengrass devices. 

For simplicity this HOWTO creates a permanent user which can and should be deleted when no longer required. If you're familiar
with creating and using temporary credentials, rather use those.

### Policies

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_Policies_](https://console.aws.amazon.com/iamv2/home#/policies) tab
3. Create the following policies:

- _uhppoted-greengrass_ policy for  provisioning (a.k.a. installing and configuring) the AWS Greengrass 'core' and 'thing' 
devices. 
- _UhppotedGreengrassTokenExchangeRoleAccess_ policy for the _core_ device to accesss S3 and Cloudwatch

#### _uhppoted-greengrass_

The _uhppoted-greengrass_ policy below is based on the [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html) from the AWS Greengrass Developer Guide.

1. Click on _Create policy_
2. Open the _JSON_ tab and paste the following policy, **replacing the \<account-id\> with the Amazon
   account ID** (it's in the account drop-down at the top right of the page):
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
                "arn:aws:iam::<account-id>:role/UhppotedGreengrassTokenExchangeRole",
                "arn:aws:iam::<account-id>:policy/UhppotedGreengrassTokenExchangeRoleAccess",
                "arn:aws:iam::<account-id>:policy/UhppotedGreengrassCoreTokenExchangeRoleAlias"
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

_NOTE: the resource ARNs in the above policy are permissive for the purposes of this guide. You probably want to restrict
them to a smaller set of resources. The Policy Editor is your friend._

3. Click on _Next: Tags_
4. Click on _Next: Review_
5. Fill in the name and description fields:
   - `Name`: `uhppoted-greengrass`
   - `Description`: _Greengrass policy for deploying _uhppoted-mqtt_
6. Click on _Create Policy_


#### _UhppotedGreengrassTokenExchangeRoleAccess_

The _UhppotedGreengrassTokenExchangeRoleAccess_ policy below is based on the policy from [Service role permissions for core devices](https://docs.aws.amazon.com/greengrass/v2/developerguide/device-service-role.html)
from the AWS Greengrass Developer Guide. The Greengrass installer creates this policy if it doesn't exist but it's messy about it and by creating it upfront you avoid weird warnings and having to
run the installer twice.

1. Click on _Create policy_
2. Open the _JSON_ tab and paste the following policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogStreams",
                "s3:GetBucketLocation"
            ],
            "Resource": "*"
        }
    ]
}
```
3. Click on _Next: Tags_
4. Click on _Next: Review_
5. Fill in the name and description fields:
   - `Name`: `UhppotedGreengrassTokenExchangeRoleAccess`
   - `Description`: _Greengrass policy for uhppoted-greengrass core device access to S3 and CloudWatch
6. Click on _Create Policy_


### Roles

#### AWS Greengrass Service Role

_Ref._ [Greengrass service role](https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-service-role.html)

_Greengrass_ requires a service role to for provisioning and managing AWS IoT devices (the service role is in the
[_AWS Iot console_](https://console.aws.amazon.com/iot/home) under _Settings_, right at the very bottom). 

If you do not have a service role, create one in [IAM](https://console.aws.amazon.com/iamv2/home):
   1. Open the _Roles_ section
   2. Click on _Create_
   3. Choose:
      - _AWS Service_
      - _Greengrass_ (under _Use cases for other AWS services_)
      - Permissions: _AWSGreengrassResourceAccessRolePolicy_
      - Name: _Greengrass_ServiceRole_

   4. In the [_AWS Iot console_](https://console.aws.amazon.com/iot/home) under _Settings_ attach the newly created
     _Greengrass_ServiceRole_


#### Token Exchange Role

_Ref._ [Device service role](https://docs.aws.amazon.com/greengrass/v2/developerguide/device-service-role.html)

The _Token Exchange Role_ sets the permissions required for the _uhppoted-greengrass_ IoT _core_ device to access
S3 and Cloudwatch. The Greengrass installer creates this role if it doesn't exist but creating it upfront avoids
some warnings and errors and having to run the installer twice.

   1. Open the _Roles_ section
   2. Click on _Create_
   3. Choose:
      - _Custom trust policy_
      - Use the following trust policy:
```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "credentials.iot.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```
      - Attach the UhppotedGreengrassTokenExchangeRoleAccess policy created above
      - Update the name and description fields:
        - Name: _UhppotedGreengrassTokenExchangeRole_
        - Description: _Token exchange role for uhppoted-greengrass IoT core device to access S3 and Cloudwatch_

   4. Click on _Create_


### Groups

In the [_AWS IAM console_](https://console.aws.amazon.com/iamv2), create a _uhppoted-greengrass_ group for the users to
be given the permissions required to provision the AWS Greengrass `core` and `thing` devices:

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_User groups_](https://console.aws.amazon.com/iamv2/home#/groups) page
3. Click on _Create group_
4. Enter the group name _uhppoted-greengrass_
5. Attach the _uhppoted-greengrass_ policy created above to the group
6. Attach the _UhppotedGreengrassTokenExchangeRoleAccess_ policy created above to the group
7. Click on _Create group_


### Users

_(skip this section if you're using temporary credentials)_

In the AWS IAM console, create a _uhppoted-greengrass_ user for provisioning the AWS Greengrass 'core' and 'thing'
devices:

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_Users_](https://console.aws.amazon.com/iamv2/home#/users) page
3. Click _Add users_
4. Enter the user name _uhppoted-greengrass_
5. Select the _Access key - Programmatic access_ AWS credential type
6. Click _Next: Permissions_
7. Check the _uhppoted-greengrass_ group under the _Add user to group_ section
8. Click _Next: Tags_
9. Click _Next: Review_
10. Click _Create user_
11. **Copy the access key and secret key for later use**
12. Click _Close_

