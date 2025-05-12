# Acorn Nameserver prototype

Name server serves as a single, reliable source for EPICS and ACNET control system component information such as device definitions, property configurations, and associated metadata.
<!-- TOC -->

- [Acorn Nameserver prototype](#acorn-nameserver-prototype)
    - [Background](#background)
    - [Features](#features)
    - [Implementation](#implementation)
    - [Prerequisites](#prerequisites)
    - [Running the application in dev mode](#running-the-application-in-dev-mode)
        - [Dev service initialization](#dev-service-initialization)
    - [Run test client](#run-test-client)
    - [Setting up for production](#setting-up-for-production)
        - [Setting up data schema management](#setting-up-data-schema-management)
        - [Configure application](#configure-application)
    - [Package and deploying the application](#package-and-deploying-the-application)
        - [Package as a JAR file](#package-as-a-jar-file)
        - [Package as a container image](#package-as-a-container-image)
        - [Package as a Kubernetes application](#package-as-a-kubernetes-application)
    - [Other services provided](#other-services-provided)
        - [gRPC Reflection](#grpc-reflection)
        - [Health check](#health-check)
    - [Notes about source code](#notes-about-source-code)

<!-- /TOC -->

## Background
 - [Requirements](docs/nameserver-requirements.pdf)


## Features

This prototype does the following:
- Hosts component information about the control system devices (see [data model design document](docs/nameserver-data-schema-design.pptx) for more detail), including nodes, locations, alarm and access control properties, etc...
- Uses secure gRPC as the protocol for service-to-service communication
- Uses PostgreSQL for storage
- Uses OpenID Connect Role-based Access Control for authentication and authorization that can be integrated with Keycloak
- Database schema version control and automated database deployments using Flyway 

## Implementation

This project uses Quarkus, a Java framework geared towards building Kubernetes native Java applications.

Due to this fact, using this framework makes it easy to add on common performance optimizations, implement common protocols, and integrate with other common services.

For more information, go to [quarkus.io](https://quarkus.io/)

## Prerequisites
For prerequisites for Quarkus, refer to the prerequisite lists in the [Creating your first application](https://quarkus.io/guides/getting-started) and [Your second Quarkus application](https://quarkus.io/guides/getting-started-dev-services) tutorials. Specifically, the following need to be installed:
- JDK 17+ installed with JAVA_HOME configured appropriately
- Apache Maven 3.9.9
- A working container runtime (Docker or [Podman](https://quarkus.io/guides/podman)) is needed to run in dev mode

**_NOTE:_** Application was tested using Podman.  As such, instructions below use Podman.  To use Docker, replace `podman` with `docker`

For the test client, you need the following installed:
- [Go](https://go.dev/doc/install)
- [Protocol Buffer Compiler](https://protobuf.dev/installation/)

Run the following to install the plugin for generating Go code ([see docs](https://protobuf.dev/reference/go/go-generated/)):
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Server was tested on Linux operating system.


## Running the application in dev mode

1. Go to `nameserver` folder, this is the server project directory
2. Go to `src/main/resources/ssl` and run `./gen-certs.sh`.  This will generate a self-signed certificate that can be used for testing.
3. Run your application in dev mode by running:


```shell script
./mvnw quarkus:dev
```

This mode will allow live editing.  The server will by default listen from port 8443. The URL and port it's listening on will be printed on the console with a message like this:

```shell 
acorn-nameserver 1.0.0-SNAPSHOT on JVM (powered by Quarkus 3.18.3) started in 32.681s. Listening on: https://0.0.0.0:8443
```

Dev mode also automatically starts containerized Postgres and Keycloak services for testing as well as a container running the "ryuk" service which is responsible for properly cleaning up services when exiting out of dev mode. Run `podman ps` to see running containers:
```shell
podman ps
CONTAINER ID  IMAGE                                 COMMAND               CREATED         STATUS         PORTS                    NAMES
fbf4f903bac2  docker.io/testcontainers/ryuk:0.11.0  /bin/ryuk             24 minutes ago  Up 24 minutes  0.0.0.0:44455->8080/tcp  testcontainers-ryuk-496632b7-bbff-421c-9081-f6c4b204eaa9
752675033d8f  docker.io/library/postgres:17         --max_prepared_tr...  24 minutes ago  Up 24 minutes  0.0.0.0:41137->5432/tcp  wizardly_borg
dee29616ed5a  quay.io/keycloak/keycloak:26.0.7      start --http-enab...  24 minutes ago  Up 24 minutes  0.0.0.0:34753->8080/tcp  friendly_cori
```

**_NOTE:_** Quarkus also comes with a Dev UI, that helps with interacting with application and supporting services.  Go to `[SERVER_ADDRESS]/q/dev-ui` (for example: <https://localhost:8443/q/dev-ui>) to view.

### Dev service initialization
They Keycloak dev service is prepopulated using the configuration from the [Keycloak Dev Services example](https://quarkus.io/version/3.15/guides/security-openid-connect-dev-services), defined in `config/quarkus-realm.json`. This configuration defines users `alice`, and `bob`, whose passwords are the same as their usernames.  `alice` is given both `admin` and `user` roles, and `bob` is given `user` role only.

The Postgres dev service has the database automatically setup and is prepopulated using `src/main/resources/import.sql`. This only applied in dev mode.


## Run test client

1. Go to the `test_client` directory
2. Run `./gen-code.sh` to generate the Go code for gRPC
3. Run `go build`. This will create the executable `nsclient`
4. Lookup the keycloak server port by running `podman ps` and looking at the port of the image running keycloak. For example, if the port shows `0.0.0.0:34753->8080/tcp`, the keycloak server port is 34753. 
4. Run the following to see the existing devices in the server
```shell
./nsclient list device --keycloak-url=http://localhost:KEYCLOAK_PORT
2025/05/08 11:24:23 Connecting to localhost:8443
2025/05/08 11:24:23 List devices: [name:"Device A" description:"Device A description" node_hostname:"Node A"]
```

Let's now try adding a new device.  Input for the device can be passed in via JSON data.  To see the JSON format run
```shell
./nsclient create device --keycloak-url=http://localhost:KEYCLOAK_PORT --json-template
{"name":"Device","description":"(optional)Description","nodeHostname":"NodeHostname"}
```

To create the device:
```shell
./nsclient create device --keycloak-url=http://localhost:KEYCLOAK_PORT --json='{"name":"Device B", "description": "my new device", "nodeHostname":"Node A"}'
```
This example links the new device to the same node as the already existing device.  

Similar commands can be run on other resources such as channels, nodes, locations, etc. To see list of commands, run `./nsclient -h`.  

## Setting up for production
### Setting up data schema management
In dev mode, the quarkus automatically creates the database at startup using Hibernate ORM database generation based on the Hibernate ORM Java classes defined in `src/main/java/nameserver/model`.  It will also re-create the database tables if any of the Hibernate ORM classes change.  This is great when you're in development! Once you're at the point you want to deploy the server, you can use Flyway to create the initial schema and manage any schema changes.

Open `pom.xml` and uncomment the following lines:

```xml
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-flyway</artifactId>
        </dependency>
        <dependency>
        <groupId>org.flywaydb</groupId>
        <artifactId>flyway-database-postgresql</artifactId>
        </dependency>
```

Run the application in dev mode and go to the Dev UI > Extensions > Datasources link in the Flyway panel (See [previous section](#running-the-application-in-dev-mode) about how to get to the Dev UI). Hit the Migration button (see [documentation](https://quarkus.io/guides/hibernate-orm#flyway)).  This will create the initial data schema file `V1.0.0__acorn-nameserver.sql` under `src/main/resources/db/migration`. 

To have the schema file automatically applied startup, add `quarkus.flyway.migrate-at-start=true` to the `application.properties` file in `src/main/resources`.
To make further schema changes, update the Hibernate ORM Java classes and add a `V<VERSION>__<NAME>.sql` where `<VERSION>` is the version number to the `db/migration` directory.

Flyway keeps track which migration versions have been applied already in a separate table, to prevent migration files from being applied more than once.



### Configure application

Open the `application.properties` file in `src/main/resources/`, make sure the following settings are set correctly to work with your external services

Postgres settings
```
%prod.quarkus.datasource.db-kind = postgresql
%prod.quarkus.datasource.username = quarkus
%prod.quarkus.datasource.password = quarkus
%prod.quarkus.datasource.jdbc.url = jdbc:postgresql://localhost:5432/nameserver
```

Keycloak settings
```
%prod.quarkus.oidc.auth-server-url=http://localhost:8180/realms/quarkus
quarkus.oidc.client-id=backend-service
quarkus.oidc.credentials.secret=secret
```

SSL settings
```
quarkus.http.ssl.certificate.files=ssl/server.crt
quarkus.http.ssl.certificate.key-files=ssl/server.key
```

Note that if you need to have different configuration between development (i.e dev mode), test, and in production, prefix settings with `%dev`, `%test`, or `%prod` to only use the settings in either develeopment, test, or production.


## Package and deploying the application
Quarkus may different options for packaging and deploying applicaitons. The ones covered are:
1. As a JAR file
2. As a container
3. As a Kubernetes application (brief)

###  Package as a JAR file
The instructions for packaging as a JAR file is included here as a quick way to run and test the nameserver as a packaged application.

The application can be packaged using:

```shell script
./mvnw package
```

It produces the `quarkus-run.jar` file in the `target/quarkus-app/` directory.

The application is now runnable using `java -jar target/quarkus-app/quarkus-run.jar`. 

### Package as a container image

Add the quarkus extension for the type of container image you are using.  See list of extensions [here](https://quarkus.io/guides/container-image#container-image-extensions).


Once the extension is added, run the following to create a container:
```shell script
./mvnw package -Dnative -Dquarkus.native.container-build=true -Dquarkus.container-image.build=true
```
This uses `Dockerfile.native` in `src/main/docker` to build the image. 

You should find an container image like this:
```shell
#show image path
$ podman images
REPOSITORY                                         TAG             IMAGE ID      CREATED         SIZE
localhost/echandler/acorn-nameserver               1.0.0-SNAPSHOT  f8c924274a41  2 minutes ago   203 MB
```

Here's an example of running the container. It works with unmodified settings from `applications.properties`

1.  Start up the supporting services using helper scripts.
```shell
$ cd test_scripts
$ ./start-services.sh
Starting postgres
7c1fee4a608b7d4fe1fc6ba49a17a04cc8f0e72af84e65a9474bbbe0805a886c
Creating keycloak test image
STEP 1/3: FROM quay.io/keycloak/keycloak:26.0.7
STEP 2/3: COPY quarkus-realm.json /opt/keycloak/data/import/realm-export.json
--> 5fb5918bc404
STEP 3/3: CMD ["start-dev", "--import-realm"]
COMMIT ns-test-keycloak:0
--> 33fd8e5dc77d
Successfully tagged localhost/ns-test-keycloak:0
33fd8e5dc77d7330169bc3cb7b81a4404cc081090e54ba26e1d7ee7e9f5e3514
Starting keycloak
29fc91bf9c8e72dbcb74bca1dc43dd7c0ee4566644cec12f9928fbf2eb4394ce
All services started
$ podman ps
CONTAINER ID  IMAGE                          COMMAND               CREATED        STATUS        PORTS                   NAMES
7c1fee4a608b  docker.io/library/postgres:17  postgres              3 minutes ago  Up 3 minutes  0.0.0.0:5432->5432/tcp  ns-test-postgres
29fc91bf9c8e  localhost/ns-test-keycloak:0   start-dev --impor...  3 minutes ago  Up 3 minutes  0.0.0.0:8180->8080/tcp  ns-test-keycloak
```
2. Run nameserver container.  Use host network to access exposed ports from postgres and keycloak containers
```shell
$ podman run --network=host -p 8443:8443 localhost/echandler/acorn-nameserver:1.0.0-SNAPSHOT
```

### Package as a Kubernetes application

Quarkus has support for automatically generating Kubernetes resources and deploying your application to a Kubernetes cluster.  This was not tested, however, below are useful links to learn more information:

- <https://quarkus.io/guides/grpc-kubernetes>
- <https://quarkus.io/guides/deploying-to-kubernetes>

Note that SmallRye Health extension is already added to the application and `quarkus.grpc.server.use-separate-server` is set to `false`, so you don't need to change `quarkus.kubernetes.ingress.target-port`.


## Other services provided
Below list several other services provided:
### gRPC Reflection 
Server has gRPC reflection enabled to allow clients lookup what methods are available.
```shell
#List gRPC services available
$ grpcurl -key server.key -cert server.crt -insecure localhost:8443 list
NameService
grpc.health.v1.Health
#Describe the NameService
$ grpcurl -key server.key -cert server.crt -insecure localhost:8443 describe NameService
NameService is a service:
service NameService {
  rpc AddChannelAccessControl ( .AddChannelAccessControlRequest ) returns ( .ChannelAccessControlResponse );
  rpc AddChannelAlarm ( .AddChannelAlarmRequest ) returns ( .ChannelAlarmResponse );
  rpc AddChannelTransform ( .AddChannelTransformRequest ) returns ( .ChannelTransformResponse );
  rpc CreateAlarmType ( .CreateAlarmTypeRequest ) returns ( .AlarmTypeResponse );
  rpc CreateChannel ( .CreateChannelRequest ) returns ( .ChannelResponse );
  rpc CreateDevice ( .CreateDeviceRequest ) returns ( .DeviceResponse );
  rpc CreateLocation ( .CreateLocationRequest ) returns ( .LocationResponse );
  rpc CreateLocationType ( .CreateLocationTypeRequest ) returns ( .LocationTypeResponse );
  rpc CreateNode ( .CreateNodeRequest ) returns ( .NodeResponse );
  rpc CreateRole ( .CreateRoleRequest ) returns ( .RoleResponse );
...
}
#Describe a message type
$ grpcurl -key server.key -cert server.crt -insecure localhost:8443 describe .CreateChannelRequest
CreateChannelRequest is a message:
message CreateChannelRequest {
  .Channel channel = 1;
}
```
### Health check
Quarkus gRPC provides a health check
```shell
$ grpcurl -key server.key -cert server.crt -insecure localhost:8443 grpc.health.v1.Health/Check

{
  "status": "SERVING"
}
```
There is also a Health REST API you can access by going to  `[SERVER_URL]/q/health`. This is part of the `smallrye-health` Quarkus extension. See [documentation](https://quarkus.io/guides/smallrye-health#running-the-health-check) for more details. 

## Notes about source code
Source code is at `src/main/java/nameserver`.  The main entry point is in `NameServiceImpl.java`, which defines the main gRPC service.  Code uses Hiberate ORM to define the SQL data schema and handles interactions with the database.  The Hibernate ORM entities are defined in the `nameserver/model`.

Proto files are defined in `src/main/proto`

