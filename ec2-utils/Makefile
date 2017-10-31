lint:
	@rpmlint ec2-utils.spec

rpm:
	@rpmbuild -ba ec2-utils.spec --define "_sourcedir $$PWD"

.PHONY: lint, rpm
