# IoT 

## Policies

UhppotedGreengrassThingPolicy

{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iot:Connect",
        "iot:Publish",
        "iot:Subscribe",
        "iot:Receive",
        "greengrass:*"
      ],
      "Resource": "*"
    }
  ]
}

GreengrassTESCertificatePolicyUhppotedGreengrassCoreTokenExchangeRoleAliasInfo

{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "iot:AssumeRoleWithCertificate",
    "Resource": "arn:aws:iot:us-east-1:026688924291:rolealias/UhppotedGreengrassCoreTokenExchangeRoleAlias"
  }
}


## Role Aliases

UhppotedGreengrassCoreTokenExchangeRoleAliasInfo

Role: UhppotedGreengrassTokenExchangeRole