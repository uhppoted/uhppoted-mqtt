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

### AWS Greengrass Service Role

_Ref._ [Greengrass service role](https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-service-role.html)

_Greengrass_ requires a service role to for provisioning and managing AWS IoT devices. The service role is in the
[_AWS Iot console_](https://console.aws.amazon.com/iot/home) under _Settings_ (right at the very bottom). 

If you do not have a service role, create one in [IAM](https://console.aws.amazon.com/iamv2/home):
- Open the _Roles_ section
- Click on _Create_
- Choose:
   - _AWS Service_
   - _Greengrass_ (under _Use cases for other AWS services_)
   - Permissions: _AWSGreengrassResourceAccessRolePolicy_
   - Name: _Greengrass_ServiceRole_

- In the [_AWS Iot console_](https://console.aws.amazon.com/iot/home) under _Settings_ attach the newly created
  _Greengrass_ServiceRole_

### Policies

In the [_AWS IAM console_](https://console.aws.amazon.com/iamv2), create two policies:

1. A _uhppoted-greengrass_ policy for provisioning (a.ka. installing and configuring) the AWS Greengrass 'core' and
   'thing' devices. 
2. An _(optional) _uhppoted-greengrass-cli_ policy for the _AWS Greengrass CLI_

The _uhppoted-greengrass_ policy described below is based on the [Minimal IAM policy for installer to provision resources](https://docs.aws.amazon.com/greengrass/v2/developerguide/provision-minimal-iam-policy.html) from the AWS Greengrass Developer Guide.

The _uhppoted-greengrass-cli_ policy is a convenience for this HOWTO and is not required if you don't anticipate needing
to use the AWS Greengrass CLI to debug/manage 'core' or 'thing' devices. As per the recommendation in the 
[Greengrass CLI guide](https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-cli-component.html), the CLI 
provides an unnecessary level of access to the system and should not be enabled for systems in production i.e. once you're
up and running, it is a **really good idea** to delete the _uhppoted-greengrass-cli_ policy.


#### `uhppoted-greengrass`

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_Policies_](https://console.aws.amazon.com/iamv2/home#/policies) tab
3. Click on _Create policy_
4. Open the _JSON_ tab and paste the following policy, **replacing the \<account-id\> with the Amazon
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

_NOTE: the resource ARNs in the above policy are permissive for the puposes of this guide. You probably want to restrict
them to a smaller set of resources. The Policy Editor is your friend._

6. Click on _Next: Tags_
7. Click on _Next: Review_
8. Fill in the name and description fields:
   - `Name`: `uhppoted-greengrass`
   - `Description`: Greengrass policy for deploying _uhppoted-mqtt_
9. Click on _Create Policy_


#### `uhppoted-greengrass-cli`

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_Policies_](https://console.aws.amazon.com/iamv2/home#/policies) tab
3. Click on _Create policy_
4. Open the _JSON_ tab and paste the following policy:
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

5. Click on _Next: Tags_
6. Click on _Next: Review_
7. Fill in name and description fields:
   - `Name`: `uhppoted-greengrass-cli`
   - `Description`: Greengrass policy for the CLI
8. Click on _Create Policy_


### Groups

In the [_AWS IAM console_](https://console.aws.amazon.com/iamv2), create two groups:

1. A _uhppoted-greengrass_ group for the users to be given the permissions required to provision the AWS Greengrass `core` and
   `thing` devices. 
2. A _uhppoted-greengrass-cli_ group for the users to be given the permissions required to use the AWS Greengrass CLI (optional).

Steps:

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_User groups_](https://console.aws.amazon.com/iamv2/home#/groups) page
3. Click on _Create group_
4. Enter the group name _uhppoted-greengrass_
5. Attach the _uhppoted-greengrass_ policy created above to the group
6. Click on _Create group_

Optionally, repeat steps 3-6 to create a _uhppoted-greengrass-cli_ group with the _uhppoted-greengrass-cli_ policy attached.


### Users

_(skip this section if you're using temporary credentials)_

In the AWS IAM console, create two users:

1. A _uhppoted-greengrass_ user for provisioning the AWS Greengrass 'core' and 'thing' devices. 
2. A _uhppoted-greengrass-cli_ user for the AWS Greengrass CLI (optional).

Steps:

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
11. Copy the access key and secret key for later use
12. Click _Close_

Optionally, repeat steps 3-12 to create a _uhppoted-greengrass-cli_ user in the _uhppoted-greengrass-cli_ group.


