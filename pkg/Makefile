VERSION := 0.3.0

rpm:
	@make -C ../ build && cp ../goat ./
	@rpmlint specfile.spec
	@rpmbuild -ba specfile.spec --define "_sourcedir $$PWD" --define "_version $(VERSION)"
