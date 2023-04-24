# Local run
```shell
DOCUMENTS_SERVICE_CONTEXT_DOCUMENT=. \
docker-compose -f docker-compose.yml \
-f document/deployments/docker-compose.yml \
up --build
```