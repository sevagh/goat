### Goat for ENI

#### Behavior

`goat eni` should behave correctly with no parameters. It is configured entirely with tags (explained [below](#tags)). It logs to `stderr` by default.

It takes some options:

* `--dry` - dry run, don't execute any commands
* `--log-level=<level>` - logrus log levels (i.e. debug, info, warn, error, fatal, panic)
* `--debug` - an interactive debug mode which prompts to continue after every phase so you can explore the state between phases

#### Fresh run

The event flow is roughly the following:

* Get EC2 metadata on the running instance
* Use metadata to establish an EC2 client and scan ENIs
* Attach the ENIs it needs based on their tags

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
