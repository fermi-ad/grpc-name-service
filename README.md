# Acorn Nameserver prototype

Name server serves as a single, reliable source for EPICS and ACNET control system component information such as device definitions, property configurations, and associated metadata.


## Background
 - Requirements
 - Data model design

## Features

This prototype does the following:
- Hosts component information about the control system devices (see data model design document for more detail), including nodes, locations, alarm and access control properties, etc...
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
- A working container runtime (Docker or [Podman](https://quarkus.io/guides/podman))

**_NOTE:_** Application was tested using Podman.  As such, instructions below use Podman.  To use Docker, replace `podman` with `docker`

For the test client, you need the following installed:
- [Go](https://go.dev/doc/install)
- [Protocol Buffer Compiler](https://protobuf.dev/installation/)

Run the following to install the plugin for [generating Go code](https://protobuf.dev/reference/go/go-generated/)):
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

``shell 
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

Quarkus also comes with a Dev UI, that helps with interacting with application and supporting services.  Go to [SERVER_ADDRESS/q/dev-ui](https://localhost:8443/q/dev-ui) to view.

## Dev service initialization
They Keycloak dev service is prepopulated using the configuration from the [Keycloak Dev Services example](https://quarkus.io/version/3.15/guides/security-openid-connect-dev-services), defined in `config/quarkus-realm.json`. This configuration defines 23 users, `alice`, and `bob`, whose passwords are the same as their usernames.  `alice` is given both `admin` and `user` roles, and `bob` is given `user` role only.

The Postgres dev service has the database automatically setup and is prepopulated using `src/main/resources/import.sql`


## Run test client

1. Go to the `test_client` directory
2. Run `./gen-code.sh` to generate the Go code for gRPC
3. Run `go build`. This will create the executable `nsclient`
4. Lookup the keycloak server port by running `podman ps` and looking at the port of the image running keycloak. For example, if the port shows `0.0.0.0:34753->8080/tcp`, the keycloak server port is 34753. 
4. Run the following to see the existing devices in the server
```shell
./nsclient list device --keycloak-url=http://localhost:KEYCLOAK_PORT
2025/05/08 11:24:23 Connecting to localhost:8443
2025/05/08 11:24:23 List devices: [name:"Device A" description:"Device A description" node_hostname:"Node A" name:"Device B" description:"" node_hostname:"Node A"]
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

## Setting up data schema management
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

Run the application in dev mode and go to the Dev UI, click in the Datasources link in the Flyway pane. Hit the Create Initial Migration button

## Packaging and running the application

The application can be packaged using:

```shell script
./mvnw package
```

It produces the `quarkus-run.jar` file in the `target/quarkus-app/` directory.
Be aware that it’s not an _über-jar_ as the dependencies are copied into the `target/quarkus-app/lib/` directory.

The application is now runnable using `java -jar target/quarkus-app/quarkus-run.jar`.

If you want to build an _über-jar_, execute the following command:

```shell script
./mvnw package -Dquarkus.package.jar.type=uber-jar
```

The application, packaged as an _über-jar_, is now runnable using `java -jar target/*-runner.jar`.

## Creating a native executable

You can create a native executable using:

```shell script
./mvnw package -Dnative
```

Or, if you don't have GraalVM installed, you can run the native executable build in a container using:

```shell script
./mvnw package -Dnative -Dquarkus.native.container-build=true
```

You can then execute your native executable with: `./target/acorn-nameserver-1.0.0-SNAPSHOT-runner`

If you want to learn more about building native executables, please consult <https://quarkus.io/guides/maven-tooling>.

## Provided Code

### REST

Easily start your REST Web Services

[Related guide section...](https://quarkus.io/guides/getting-started-reactive#reactive-jax-rs-resources)

# Personal notes

Start in dev mode
```
quarkus dev
```

Run psql in container
```
podman ps
podman exec -it exciting_leakey psql -U quarkus -d quarkus -c "SELECT * FROM location;"
```
Inspect GRPC API
```
grpcurl -plaintext localhost:9000 list
```

Dev UI URL:  http://localhost:8080/q/dev-ui/io.quarkus.quarkus-grpc/services

The application can be packaged using:

```shell script
./mvnw package
```

It produces the `quarkus-run.jar` file in the `target/quarkus-app/` directory.
Be aware that it’s not an _über-jar_ as the dependencies are copied into the `target/quarkus-app/lib/` directory.

The application is now runnable using `java -jar target/quarkus-app/quarkus-run.jar`.

