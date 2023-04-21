gen: gen-open-api

gen-open-api: gen-openapi-documents

gen-openapi-documents:
	ogen -allow-remote -convenient-errors=off -target=./../pkg/gen/open-api/documents \
	-package=documents ../api/open-api/documents/spec.yml

