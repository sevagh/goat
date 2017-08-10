%define pkgname goat

Name: %{pkgname}
Version: %{_version}
Release: 1%{?dist}
Summary: Attach and mount EBS volumes 

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
%{__install} -m0777 %{SOURCE1} %{buildroot}/%{_unitdir}/%{pkgname}.service


%files
%{_bindir}/%{pkgname}
%{_unitdir}/%{pkgname}.service


%post
if [ $1 -eq 1 ]; then
        /bin/systemctl daemon-reload >/dev/null 2>&1 || :
fi


%preun
if [ $1 -eq 0 ] ; then
        # Package removal, not upgrade
        /bin/systemctl disable goat@*.service >/dev/null 2>&1 || :
        /bin/systemctl stop goat@*.service >/dev/null 2>&1 || :
fi


%postun
/bin/systemctl daemon-reload >/dev/null 2>&1 || :


%changelog
* Thu Aug 10 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- Goat subcommands
* Tue Jul 11 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- First RPM package for goat 
