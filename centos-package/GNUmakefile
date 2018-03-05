rpm:
	@make -C ../ release && cp ../dist/bin/goat ./
	@rpmlint specfile.spec
	@rpmbuild -ba specfile.spec --define "_sourcedir $$PWD" --define "_version $$GOAT_VERSION"
