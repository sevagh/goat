%define pkgname goat

Name: %{pkgname}
Version: %{_version}
Release: 1%{?dist}
Summary: Attach and mount EBS and ENI

License: BSD 3-clause
URL: https://github.com/sevagh/goat 
Source0: %{pkgname}
Source1: %{pkgname}@.service

Requires: systemd mdadm


%description
Automatically attach AWS resources to a running EC2 instance.


#%prep
#%setup
#%build


%install
%{__mkdir} -p %{buildroot}/%{_bindir}
%{__mkdir} -p %{buildroot}/%{_unitdir}
%{__install} -m0775 %{SOURCE0} %{buildroot}/%{_bindir}/%{pkgname}
%{__install} -m0777 %{SOURCE1} %{buildroot}/%{_unitdir}/%{pkgname}@.service


%files
%{_bindir}/%{pkgname}
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
* Tue Apr 17 2018 Sevag Hanssian <sevag.hanssian@gmail.com>
- Recombine commands
* Tue Mar 06 2018 Sevag Hanssian <sevag.hanssian@gmail.com>
- Split subcommands into two binaries
* Thu Aug 10 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- Goat subcommands
* Tue Jul 11 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- First RPM package for goat 
