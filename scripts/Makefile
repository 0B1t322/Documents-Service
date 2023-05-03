gen: gen-open-api

gen-open-api: gen-openapi-documents gen-openapi-sessions

gen-openapi-documents:
	ogen -allow-remote -convenient-errors=off -target=./../pkg/gen/open-api/documents \
	-package=documents ../api/open-api/documents/spec.yml
gen-openapi-sessions:
	ogen -allow-remote -convenient-errors=off -target=./../pkg/gen/open-api/sessions \
    	-package=sessions ../api/open-api/sessions/spec.yml

