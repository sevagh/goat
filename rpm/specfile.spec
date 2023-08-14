%define pkgname goat

Name: %{pkgname}
Version: %{_version}
Release: 1%{?dist}
Summary: Attach and mount EBS and ENI
BuildRequires: systemd-rpm-macros

License: BSD 3-clause
URL: https://github.com/sevagh/goat 
Source0: %{pkgname}
Source1: %{pkgname}@.service
Source2: %{pkgname}.cfg

Requires: systemd mdadm


%description
Automatically attach AWS resources to a running EC2 instance.


#%prep
#%setup
#%build


%install
%{__mkdir} -p %{buildroot}/%{_sbindir}
%{__mkdir} -p %{buildroot}/%{_unitdir}
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/sysconfig
%{__install} -m0775 %{SOURCE0} %{buildroot}/%{_sbindir}/%{pkgname}
%{__install} -m0644 %{SOURCE1} %{buildroot}/%{_unitdir}/%{pkgname}@.service
%{__install} -m0644 %{SOURCE2} %{buildroot}/%{_sysconfdir}/sysconfig/%{pkgname}


%files
%config(noreplace) /etc/sysconfig/goat
%{_sbindir}/%{pkgname}
%{_unitdir}/%{pkgname}@.service


%post
if [ $1 -eq 1 ]; then
        /bin/systemctl daemon-reload >/dev/null 2>&1 || :
fi


%preun
if [ $1 -eq 0 ] ; then
        # Package removal, not upgrade
        /bin/systemctl disable %{pkgname}@*.service >/dev/null 2>&1 || :
        /bin/systemctl stop %{pkgname}@*.service >/dev/null 2>&1 || :
fi


%postun
/bin/systemctl daemon-reload >/dev/null 2>&1 || :


%changelog
* Wed Sep 25 2019 Andrew Cutler <andrew@panubo.com>
- Add configuration file
* Tue Apr 17 2018 Sevag Hanssian <sevag.hanssian@gmail.com>
- Recombine commands
* Tue Mar 06 2018 Sevag Hanssian <sevag.hanssian@gmail.com>
- Split subcommands into two binaries
* Thu Aug 10 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- Goat subcommands
* Tue Jul 11 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- First RPM package for goat
