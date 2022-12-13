**WORK IN PROGRESS**

## HOWTO: _AWS Greengrass CLI_

Supplements the README with a description of the IAM policies, roles, etc for installing AWS Greengrass CLI
should it become necessary. As per the recommendation in the [Greengrass CLI guide](https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-cli-component.html), the CLI provides an unnecessary level of 
access to the system and should not be enabled for systems in production i.e. once you're up and running, 
it is a **really good idea** to delete the _uhppoted-greengrass-cli_ policy, group and user.


For more information, see [Greengrass CLI](https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-cli-component.html).

### AWS IAM

The basic requirements are:

1. A _uhppoted-greengrass-cli_ policy for the AWS Greengrass CLI. 
2. A _uhppoted-greengrass-cli_ group for the users to be assigned the permissions required to use the AWS Greengrass
    CLI.
3. A _uhppoted-greengrass-cli_ user for the AWS Greengrass CLI.

The CLI setup is a convenience and is not required if you don't anticipate needing to use the _AWS Greengrass CLI_ to debug/manage
`core` or `thing` devices. Chances are you'll probably need it at some point though, particularly if this is your first time
through. 

##### `uhppoted-greengrass-cli` policy

In the [_AWS IAM console_](https://console.aws.amazon.com/iamv2), create a _uhppoted-greengrass-cli_ policy for 
the _AWS Greengrass CLI_:

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


#### `uhppoted-greengrass-cli` group

In the [_AWS IAM console_](https://console.aws.amazon.com/iamv2), create a _uhppoted-greengrass-cli_ group for the
users to be given the permissions required to use the AWS Greengrass CLI (optional):

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_User groups_](https://console.aws.amazon.com/iamv2/home#/groups) page
3. Click on _Create group_
4. Enter the group name _uhppoted-greengrass-cli_
5. Attach the _uhppoted-greengrass-cli_ policy created above to the group
6. Click on _Create group_

#### `uhppoted-greengrass-cli` user

In the AWS IAM console, create a _uhppoted-greengrass-cli_ user for the AWS Greengrass CLI:

1. Open the [_AWS IAM console_](https://console.aws.amazon.com/iamv2)
2. Open the [_Users_](https://console.aws.amazon.com/iamv2/home#/users) page
3. Click _Add users_
4. Enter the user name _uhppoted-greengrass-cli_
5. Select the _Access key - Programmatic access_ AWS credential type
6. Click _Next: Permissions_
7. Check the _uhppoted-greengrass-cli_ group under the _Add user to group_ section
8. Click _Next: Tags_
9. Click _Next: Review_
10. Click _Create user_
11. Copy the access key and secret key for later use
12. Click _Close_

