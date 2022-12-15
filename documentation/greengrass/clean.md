## Cleaning up

### AWS IoT

1. Disassociate _uhppoted-thing_ from _uhppoted-greengrass_
2. Delete the _uhppoted-thing_ and _uhppoted-greengrass_ things
3. Detach the service role
4. Revoke the certificates
5. Delete the _UhppotedGreengrassThingPolicy_ and _GreengrassTESCertificatePolicyUhppotedGreengrassCoreTokenExchangeRoleAliasInfo_
   policies
6. Delete the certificates
7. Delete the UhppotedGreengrassCoreTokenExchangeRoleAlias role alias
8. Delete the deployments

### AWS IAM

1. Delete the _uhppoted-greengrass_ user
2. Delete the _uhppoted-greengrass_ group
3. Delete the _Greengrass_ServiceRole_ and _UhppotedGreengrassTokenExchangeRole_ roles
4. Delete the _uhppoted-greengrass_ and _UhppotedGreengrassTokenExchangeRoleAccess_ policies

### VPS

1. Destroy the VPS







