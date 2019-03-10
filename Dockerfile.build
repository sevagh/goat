FROM golang:1.12-stretch as build

RUN apt-get update -y &&\
	apt-get install -y zip 

COPY . /goat
WORKDIR /goat

RUN make build zip tarball

FROM fedora as rpmbuild

RUN dnf update -y &&\
	dnf install -y make rpm-build alien

COPY ./rpm /goat-rpm-pkg/

COPY --from=build /goat/bin/goat /goat-rpm-pkg/goat
COPY --from=build /goat/bin/version-file /goat-rpm-pkg/version-file

WORKDIR goat-rpm-pkg

RUN make rpm

WORKDIR /root/rpmbuild/RPMS/x86_64/

RUN alien ./*.rpm

FROM alpine

COPY --from=build /goat/pkg /_goat-pkg
COPY --from=rpmbuild /root/rpmbuild/RPMS/x86_64/* /_goat-pkg/

CMD cp _goat-pkg/* goat-pkg/
