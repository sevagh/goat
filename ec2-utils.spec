%define _buildid .20

%bcond_without upstart
%bcond_with systemd

Name:      ec2-utils
Summary:   A set of tools for running in EC2
Version:   0.4
Release: 1%{?_buildid}%{?dist}
License:   Apache License 2.0
Group:     System Tools
Source0:   ec2-metadata
Source1:   ec2udev
Source2:   51-ec2-hvm-devices.rules
Source3:   52-ec2-vcpu.rules
Source4:   53-ec2-network-interfaces.rules
Source5:   75-persistent-net-generator.rules
Source6:   ec2net-functions
Source7:   ec2net.hotplug
Source8:   ec2ifup
Source9:   ec2ifdown
Source10:  ec2dhcp.sh
Source11:  ec2ifup.8
Source12:  ec2ifscan
Source13:  elastic-network-interfaces.conf
Source14:  ec2ifscan.8

Source20:  ixgbevf.conf
Source21:  acpiphp.modules

# fedora stuff
Source30:  elastic-network-interfaces.service
Source31:  60-net-hotplug.rules
Source32:  net.hotplug

URL:       http://developer.amazonwebservices.com/connect/entry.jspa?externalID=1825
BuildArch: noarch
Provides:  ec2-metadata
Obsoletes: ec2-metadata
Requires:  curl
BuildRoot: %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
ec2-utils contains a set of utilities for running in ec2.

%package -n ec2-net-utils
Summary:   A set of network tools for managing ENIs
Group:     System Tools
Requires:  initscripts
Requires:  bash >= 4
Requires:  curl
Requires:  iproute
%if %{with upstart}
Requires:  upstart
%endif
%if %{with systemd}
BuildRequires: systemd-units
%endif

%description -n ec2-net-utils
ec2-net-utils contains a set of utilities for managing elastic network
interfaces.

%prep

%build

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/opt/aws/bin
mkdir -p $RPM_BUILD_ROOT/sbin
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/udev/rules.d/
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/sysconfig/network-scripts/
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/dhcp/dhclient.d/
%if %{with upstart}
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/init/
%endif
mkdir -p $RPM_BUILD_ROOT%{_mandir}/man8/

install -m755 %{SOURCE0} $RPM_BUILD_ROOT/opt/aws/bin/
install -m755 %{SOURCE1} $RPM_BUILD_ROOT/sbin/
install -m755 %{SOURCE8} $RPM_BUILD_ROOT/sbin/
install -m755 %{SOURCE9} $RPM_BUILD_ROOT/sbin/
install -m755 %{SOURCE12} $RPM_BUILD_ROOT/sbin/
install -m644 %{SOURCE2} $RPM_BUILD_ROOT%{_sysconfdir}/udev/rules.d/
install -m644 %{SOURCE3} $RPM_BUILD_ROOT%{_sysconfdir}/udev/rules.d/
install -m644 %{SOURCE4} $RPM_BUILD_ROOT%{_sysconfdir}/udev/rules.d/
install -m644 %{SOURCE5} $RPM_BUILD_ROOT%{_sysconfdir}/udev/rules.d/
install -m644 %{SOURCE6} $RPM_BUILD_ROOT%{_sysconfdir}/sysconfig/network-scripts/
install -m755 %{SOURCE7} $RPM_BUILD_ROOT%{_sysconfdir}/sysconfig/network-scripts/
install -m755 %{SOURCE10} $RPM_BUILD_ROOT%{_sysconfdir}/dhcp/dhclient.d/
%if %{with upstart}
install -m644 %{SOURCE13} $RPM_BUILD_ROOT%{_sysconfdir}/init
%endif
install -m644 %{SOURCE11} $RPM_BUILD_ROOT%{_mandir}/man8/ec2ifup.8
ln -s ./ec2ifup.8.gz $RPM_BUILD_ROOT%{_mandir}/man8/ec2ifdown.8.gz
install -m644 %{SOURCE14} $RPM_BUILD_ROOT%{_mandir}/man8/ec2ifscan.8

%if 0%{?fedora}
install -m644 %{SOURCE31} $RPM_BUILD_ROOT%{_sysconfdir}/udev/rules.d/
install -m755 %{SOURCE32} $RPM_BUILD_ROOT%{_sysconfdir}/sysconfig/network-scripts/
%endif

%if %{with systemd}
%{__install} -d -m 0775 ${RPM_BUILD_ROOT}%{_unitdir}
%{__install} -m 0644 %{SOURCE30} ${RPM_BUILD_ROOT}%{_unitdir}
%endif

# add module configs
install -m644 -D %{SOURCE20} $RPM_BUILD_ROOT/etc/modprobe.d/ixgbevf.conf
install -m755 -D %{SOURCE21} $RPM_BUILD_ROOT/etc/sysconfig/modules/acpiphp.modules

%clean
rm -rf $RPM_BUILD_ROOT

%files
/opt/aws/bin/ec2-metadata
/sbin/ec2udev
%{_sysconfdir}/udev/rules.d/51-ec2-hvm-devices.rules
%{_sysconfdir}/udev/rules.d/52-ec2-vcpu.rules

%files -n ec2-net-utils
/sbin/ec2ifup
/sbin/ec2ifdown
/sbin/ec2ifscan
%{_sysconfdir}/udev/rules.d/53-ec2-network-interfaces.rules
%{_sysconfdir}/udev/rules.d/75-persistent-net-generator.rules
%{_sysconfdir}/modprobe.d/ixgbevf.conf
%{_sysconfdir}/sysconfig/modules/acpiphp.modules
%{_sysconfdir}/sysconfig/network-scripts/ec2net-functions
%{_sysconfdir}/sysconfig/network-scripts/ec2net.hotplug
%{_sysconfdir}/dhcp/dhclient.d/ec2dhcp.sh
%if %{with upstart}
%{_sysconfdir}/init/elastic-network-interfaces.conf
%endif
%{_mandir}/man8/ec2ifup.8.gz
%{_mandir}/man8/ec2ifdown.8.gz
%{_mandir}/man8/ec2ifscan.8.gz

%if 0%{?fedora}
%{_sysconfdir}/udev/rules.d/60-net-hotplug.rules
%{_sysconfdir}/sysconfig/network-scripts/net.hotplug
%endif

%if %{with systemd}
%{_unitdir}/elastic-network-interfaces.service
%endif

%post -n ec2-net-utils
%if %{with systemd}
%systemd_post elastic-network-interfaces.service
%endif

%preun -n ec2-net-utils
%if %{with systemd}
%systemd_preun elastic-network-interfaces.service
%endif

%postun -n ec2-net-utils
%if %{with systemd}
%systemd_postun_with_restart elastic-network-interfaces.service
%endif

%changelog
* Tue Sep 24 2013 Andrew Jorgensen <ajorgens@amazon.com>
- Add hotplug script and module config

* Mon Aug 26 2013 Ben Cressey <bcressey@amazon.com>
- Configure interfaces attached at launch time

* Wed Mar 13 2013 Andrew Jorgensen <ajorgens@amazon.com>
- Use -q to avoid using a user's .curlrc

* Sun Sep 16 2012 Ben Cressey <bcressey@amazon.com>
- Add documentation for ec2ifup and ec2ifdown

* Thu Sep 13 2012 Ben Cressey <bcressey@amazon.com>
- Optimize metadata queries for elastic interfaces

* Tue Sep 11 2012 Ben Cressey <bcressey@amazon.com>
- Adjust route table usage for elastic interfaces
- Update headers to reflect Apache 2.0 license

* Wed Sep 5 2012 Ben Cressey <bcressey@amazon.com>
- Configure elastic network interfaces via DHCP

* Wed Aug 29 2012 Andrew Jorgensen <ajorgens@amazon.com>
- Add dependency on curl for ec2-metadata
- Update ec2-metadata to 0.1.1 bugfix release

* Mon Aug 13 2012 Ben Cressey <bcressey@amazon.com>
- Add rules and scripts for MultiIP / MultiVIF support

* Mon Jul 30 2012 Ethan Faust <efaust@amazon.com>
- Added udev rules to automatically bring up vCPUs when they're added.

* Mon Aug 1 2011 Nathan Blackham <blackham@amazon.com>
- adding BuildRoot directive to specfile.

* Wed Sep 22 2010 Nathan Blackham <blackham@amazon.com>
- move to ec2-utils
- add udev code for symlinking xvd* devices to sd*
- fixing typo in spec file
- adding udev symlinks for xvd* devices

* Tue Sep 07 2010 Nathan Blackham <blackham@amazon.com>
- initial packaging of script as an rpm
- moving rpm to noarch
- adding Group line in specfile
- initial packaging of ec2-metadata
- setup complete for package ec2-metadata
