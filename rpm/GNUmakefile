VERSION:=$(shell cat version-file | tr - _)

all: rpm

rpm:
	@rpmbuild -ba specfile.spec --define "_sourcedir $(PWD)" --define "_version $(VERSION)"
