###############################################################################

# rpmbuilder:relative-pack true

###############################################################################

%define  debug_package %{nil}

###############################################################################

%define _posixroot        /
%define _root             /root
%define _bin              /bin
%define _sbin             /sbin
%define _srv              /srv
%define _home             /home
%define _opt              /opt
%define _lib32            %{_posixroot}lib
%define _lib64            %{_posixroot}lib64
%define _libdir32         %{_prefix}%{_lib32}
%define _libdir64         %{_prefix}%{_lib64}
%define _logdir           %{_localstatedir}/log
%define _rundir           %{_localstatedir}/run
%define _lockdir          %{_localstatedir}/lock/subsys
%define _cachedir         %{_localstatedir}/cache
%define _spooldir         %{_localstatedir}/spool
%define _crondir          %{_sysconfdir}/cron.d
%define _loc_prefix       %{_prefix}/local
%define _loc_exec_prefix  %{_loc_prefix}
%define _loc_bindir       %{_loc_exec_prefix}/bin
%define _loc_libdir       %{_loc_exec_prefix}/%{_lib}
%define _loc_libdir32     %{_loc_exec_prefix}/%{_lib32}
%define _loc_libdir64     %{_loc_exec_prefix}/%{_lib64}
%define _loc_libexecdir   %{_loc_exec_prefix}/libexec
%define _loc_sbindir      %{_loc_exec_prefix}/sbin
%define _loc_bindir       %{_loc_exec_prefix}/bin
%define _loc_datarootdir  %{_loc_prefix}/share
%define _loc_includedir   %{_loc_prefix}/include
%define _loc_mandir       %{_loc_datarootdir}/man
%define _rpmstatedir      %{_sharedstatedir}/rpm-state
%define _pkgconfigdir     %{_libdir}/pkgconfig

###############################################################################

Summary:         Utility for managing Vagrant repositories
Name:            vgrepo
Version:         1.1.0
Release:         0%{?dist}
Group:           Applications/System
License:         MIT
URL:             https://github.com/gongled/vgrepo

Source0:         %{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.8

Provides:        %{name} = %{version}-%{release}

###############################################################################

%description
Simple CLI utility for managing Vagrant repositories.

###############################################################################

%prep
%setup -q

%build
mkdir src && mv {github.com,pkg.re} src

export GOPATH=$(pwd)
pushd src/github.com/gongled/%{name}/
%{__make} %{?_smp_mflags} all
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_sysconfdir}
install -dm 755 %{buildroot}%{_sysconfdir}/%{name}
install -dm 755 %{buildroot}%{_sysconfdir}/%{name}/templates

install -pm 755 src/github.com/gongled/%{name}/%{name} \
                %{buildroot}%{_bindir}/

install -pm 755 src/github.com/gongled/%{name}/templates/default.tpl \
                %{buildroot}%{_sysconfdir}/%{name}/templates/

install -pm 644 src/github.com/gongled/%{name}/%{name}.knf \
                %{buildroot}%{_sysconfdir}/%{name}/

%clean
rm -rf %{buildroot}

###############################################################################

%files
%defattr(-,root,root,-)
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.knf
%config(noreplace) %{_sysconfdir}/%{name}/templates/default.tpl
%{_bindir}/%{name}

###############################################################################

%changelog
* Sat May 27 2017 Gleb Goncharov <gongled@gongled.ru> - 1.1.0-0
- Added index generator by given template file

* Sun May 21 2017 Gleb Goncharov <gongled@gongled.ru> - 1.0.0-0
- Initial build

