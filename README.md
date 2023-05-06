# Online-Document-Redactor
It's repository contains a two modules for online document system.

[click here](https://drive.google.com/file/d/1VUwKALY6H3INnXHXJ2d5IrbhYbU7TF1E/view?usp=sharing) to see final 
qualifying work.

## Directory description

* `/api` - contains open-api specifications
* `/configs` - contains configs (now only nginx config)
* `/document` - code of document service module
* `/sessions` - code of sessions service module
* `/scripts` - contains scripts such as Makefile and postgresql init

## Open-API code generation

This project use API-First approach. If you make some changes
in specification you need to regenerate code.

To generate open-api code run this command in the root of project:
```shell
make -C scripts
```

All generated code will be places in `/pkg/gen/open-api`.

If you need to change some generator options see: https://github.com/ogen-go/ogen 
and change it in `/scripts/Makefile`.

## Services

### Documents service

This service contains methods to store data about documents.

#### Configurations

```.env
# URL to postgresql database
DOCUMENTS_DATABASE_URL=postgres://user:password@example:5432/Documents?sslmode=disable
# Port on which will be listened HTTP server
DOCUMENTS_APP_PORT=8080
# URL to RabbitMQ 
DOCUMENTS_AMQP_URL=amqp://user:password@localhost:5672/
# Exchange name in which messages will be published
DOCUMENT_AMQP_EXCHANGE_NAME=documents-service.events
```

### Sessions service

This service contains methods to update elements of documents and store history of operations.

#### Configurations

```.env
# URL to postgresql database
DOCUMENTS_DATABASE_URL=postgres://user:password@example:5432/Documents?sslmode=disable
# Port on which will be listened HTTP server
SESSIONS_APP_PORT=8082
# URL to RabbitMQ 
SESSIONS_AMQP_URL=amqp://user:password@localhost:5672/
# Exchange name of documents-service events
SESSIONS_DOCUMENT_AMQP_EXCHANGE_NAME=documents-service.events
# URL to documents-service REST
SESSIONS_DOCUMENT_REST_BASE_URL=http://localhost:8080
# URL to influxdb
SESSIONS_INFLUXDB_URL=http://localhost:8086
# Token for influxdb
SESSIONS_INFLUXDB_TOKEN=some_token
```

## Local run with docker

To run locally with docker pass this command from root of the project:
```shell
DOCUMENTS_SERVICE_CONTEXT_DOCUMENT=. \
SESSIONS_SERVICE_CONTEXT_DOCUMENT=. \
docker-compose -f docker-compose.yml \
-f document/deployments/docker-compose.yml \
-f sessions/deployments/docker-compose.yml \
up --build
```

To see Swagger-UI visit http://localhost:5050/api/swagger
