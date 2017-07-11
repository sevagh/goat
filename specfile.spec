%define pkgname goat

Name: %{pkgname}
Version: 0.2.0 
Release: 1%{?dist}
Summary: Attach and mount EBS volumes 

License: BSD 3-clause
URL: https://github.com/sevagh/goat 
Source0: %{pkgname}
Source1: %{pkgname}.service

Requires: systemd


%description
Automatically attach and mount EBS volumes to a running EC2 instance.


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
/bin/systemctl enable goat.service >/dev/null 2>&1 || :
/bin/systemctl start goat.service >/dev/null 2>&1 || :


%preun
if [ $1 -eq 0 ] ; then
        # Package removal, not upgrade
        /bin/systemctl disable goat.service >/dev/null 2>&1 || :
        /bin/systemctl stop goat.service >/dev/null 2>&1 || :
fi


%postun
/bin/systemctl daemon-reload >/dev/null 2>&1 || :


%changelog
* Tue Jul 11 2017 Sevag Hanssian <sevag.hanssian@gmail.com>
- First RPM package for goat 
