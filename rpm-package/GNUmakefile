all: ebs eni

ebs:
	-rpmlint ebs/specfile.spec
	@cd ebs && cp ../../bin/goat-ebs . &&\
		rpmbuild -ba specfile.spec --define "_sourcedir $$PWD" --define "_version $$GOAT_VERSION"

eni:
	-rpmlint eni/specfile.spec
	@cd eni && cp ../../bin/goat-eni . &&\
		rpmbuild -ba specfile.spec --define "_sourcedir $$PWD" --define "_version $$GOAT_VERSION"

.PHONY: ebs eni
