## ec2-net-utils for Fedora/systemd

This is a fork of Amazon's ec2-utils with modifications to support Elastic Network Interfaces (ENI) under systemd.

The spec file produces two RPMs: ec2-utils and ec2-net-utils.  The ec2-net-utils RPM contains ENI support.  It allows you to attach an ENI to a running instance and have it work as you would expect.  Hurray!

The ec2-utils RPM just contains a shell script to lookup instance metadata.

## Install

Builds are available on [Copr](https://copr.fedoraproject.org/coprs/etuttle/ec2-utils/).  Drop the .repo in your repos.d, then `yum install ec2-net-utils`

* Imporant! Don't forget to enable the `elastic-network-interfaces` systemd unit, or ENI's won't work at boot!

## OS Support

* ✓ Fedora 20
* ? Fedora 21 (should work, as it is using network-scripts according to [the cloud kickstart](https://git.fedorahosted.org/cgit/spin-kickstarts.git/tree/fedora-cloud-base.ks?id=7f202a0e531ea178243c563a721c0c248af87219#n51))
* ? Fedora 19 not tested with recent changes
* ? CentOS7 I don't think there's an official AMI yet
* ✗ RHEL7 - the AMI uses Network Manager

## How does it work

A udev rule runs `ec2net.hotplug` when a device is added or removed, which is a script that writes interface config, including source route setup.  It relies on the primary interface having come up so it can query AWS metadata.

Another udev rule starts the `ec2-ifup@` service when an interface is added, and a third one runs `/sbin/ifdown` on device removal.  The original version from AWS relied on net.hotplug to do this, which has been removed from Fedora.

Finally, `elastic-network-interfaces.service` is run late in the boot process.  It calls `ec2ifscan` which fires another udev add event for any interface which is not configured.  This handles the case of booting with an ENI that `ec2net.hotplug` hasn't had a chance to configure yet.

## Complications

* udev add events are fired during boot, during 'attach', and a second time during boot for the unconfigured case.  Meanwhile, network-scripts expects to manage any interface with a cfg that exists at boot.  So the udev events have to be ignored in the appropriate cases.
* Fedora 20 uses a kernel feature (address lifetime) which removes expired addresses, even if dhclient isn't running.  So dhclient must be kept running or the address will be dropped.
* Systemd kills any long-running processes that are spawned by scripts that are run by udev.  To be kept alive, dhclient must be started by a service started by udev (hence `ec2-ifup@`).

## This is a mess!

Yeah, but it's not easy to untangle it from network-scripts without porting to NetworkManager, and it's not clear if NetworkManager is even the way forward (with systemd-networkd on the horizon).  If Amazon Linux ever switches to systemd, they'll probably come up with a cleaner solution.