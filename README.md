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
- Uses OIDC RBAC with HTTPS/gRPC for authentication and authorization that integrates with Keycloak
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

Recommended OS is Linux.

**_NOTE:_** Application was tested using Podman.  As such, instructions below use Podman.  To use Docker, replace `podman` with `docker`

## Running the application in dev mode

You can run your application in dev mode by running:

```shell script
./mvnw quarkus:dev
```

This mode will allow live editing, and automatically starts containerized Postgres and Keycloak services for testing.  The server will by default listen from port 8443. The URL and port it's listening on will be printed on the console with a message like this:

``shell 
acorn-nameserver 1.0.0-SNAPSHOT on JVM (powered by Quarkus 3.18.3) started in 32.681s. Listening on: https://0.0.0.0:8443
```

If TLS is disabled, the default port will change to 8080.


> **_NOTE:_**  Quarkus now ships with a Dev UI, which is available in dev mode only at <http://localhost:8443/q/dev/>.

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

