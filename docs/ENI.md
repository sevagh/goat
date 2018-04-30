### Goat for ENI

#### Behavior

`goat eni` should behave correctly with no parameters. It is configured entirely with tags (explained [below](#tags)). It logs to `stderr` by default.

It takes some options:

* `-logLevel=<level>` - logrus log levels (i.e. debug, info, warn, error, fatal, panic)
* `-debug` - an interactive debug mode which prompts to continue after every phase so you can explore the state between phases

#### Fresh run

The event flow is roughly the following:

* Get EC2 metadata on the running instance
* Use metadata to establish an EC2 client and scan ENIs
* Attach the ENIs it needs based on their tags

#### Setting up the ENI - ec2-net-utils

There's a project to port [ec2-net-utils](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html#ec2-net-utils), a tool available on Amazon Linux AMIs, to CentOS/systemd, [here](https://github.com/etuttle/ec2-utils).

I [forked](https://github.com/sevagh/ec2-utils) it and created a release. This tool is highly recommended to perform the actual setup of your ENI:

```
sudo yum install -y https://github.com/sevagh/ec2-utils/releases/download/v0.5.3/ec2-net-utils-0.5-2.fc25.noarch.rpm
sudo systemctl enable elastic-network-interfaces
```

#### DeviceIndex

ENI attachments take a parameter called [DeviceIndex](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-network-interface-attachment.html). Goat isn't smart, and always starts from DeviceIndex `1`.

This means that your EC2 instance should have no attached ENIs.

If it does, they should be the ones that `goat` was going to attach anyway, not external ENIs that have no `goat` tags.

#### Tags

These are the tags you need:

| Tag Name             | Description             | EC2     | ENI    | Tag Value (examples)                                             |
| -------------------- | ----------------------- | ------- | -----  | ---------------------------------------------------------------- |
| GOAT-IN:Prefix       | Logical stack name      | *Yes*   | *Yes*  | `my_app_v1.3.4`                                                  |
| GOAT-IN:NodeId       | EC2 id within stack     | *Yes*   | *Yes*  | `0`, `1`, `2` for 3-node kafka                                   |
