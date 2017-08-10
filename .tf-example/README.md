# example

Terraform recipe for a minimally viable example of the [goat](https://github.com/sevagh/goat) EBS/EC2 tag-based mounting system.

### Iterating

To iterate with this Terraform recipe, it's helpful to export the 3 required variables:

```
$ export TF_VAR_aws_access_key=xxxx
$ export TF_VAR_aws_secret_key=xxxx
$ export TF_VAR_keypair_name=mykeypair
```
