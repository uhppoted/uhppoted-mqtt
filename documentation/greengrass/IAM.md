## HOWTO: _AWS Greengrass IAM_

This section outlines the steps required to create an IAM user with the necessary permissions. You can tailor it to
fit your requirements but basically what you're aiming to create are:

1. An IAM policy with the necessary permissions required to create, configure and run a Greengrass 'core' 
   device with a _Moquette_ MQTT broker.
2. An IAM group with the necessary policies and permissions for users needed to create and run the devices.
3. An IAM user to use for creating, configuring and running the Greengrass devices. 

For simplicity this HOWTO creates a permanent user which can/should be deleted when no longer required. If you're familiar
with creating and using temporary credentials, rather use those.

### Policies

In the AWS IAM console, create two policies:

1. A _uhppoted-greengrass_ policy for provisioning (a.ka. installing and configuring) the AWS Greengrass 'core' and
   'thing' devices. 
2. A _uhppoted-greengrass-cli_ policy for the AWS Greengrass CLI

For this HOWTO, the _uhppoted-greengrass_ policy is based on the [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html) from the AWS Greengrass Developer Guide.

The _uhppoted-greengrass-cli_ policy is a convenience for this HOWTO and is not required if you don't anticipate needing
to use the AWS Greengrass CLI to debug/manage 'core' or 'thing' devices. 

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


