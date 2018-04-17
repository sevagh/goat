all: rpm

rpm:
	-rpmlint specfile.spec
	@rpmbuild -ba specfile.spec --define "_sourcedir $$PWD" --define "_version $$GOAT_VERSION"
