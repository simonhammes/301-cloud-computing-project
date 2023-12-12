+++
title = "API Design: OpenAPI vs. gRPC"
outputs = ["Reveal"]
+++

## API Design
# OpenAPI vs. gRPC

Simon Hammes

Cloud Computing WS 2023/24

---

## Outline
1. OpenAPI/Swagger
2. gRPC
3. Comparison
4. Demo

---

# OpenAPI/Swagger

---

## OpenAPI/Swagger
- Created in 2011 by Tony Tam
- API specification in JSON/YAML files
- Goals:
  - Automate API documentation
  - Generate code for API clients
- 2016: Specification was renamed to OpenAPI

---

# OpenAPI Specification

---

## Structure

```yaml
openapi: 3.0.0

info:
  title: Students
  description: Students API
  version: 1.0.0

servers:
  - url: https://api.hs-worms.de/v1

paths: {}

components:
  schemas: {}
```

---

## Paths

// TODO: Parameters?

```yaml
paths:
  /students:
    get:
      summary: Get all students
      responses:
        200:
          description: A list of students.
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Student'
```

---

## Components

<!-- TODO -->

```yaml
components:
  schemas:
    Student:
      type: object
      properties:
        id:
          type: integer
          description: Student ID
        name:
          type: string
          description: Name
        courses:
          type: array
          items:
            $ref: '#/components/schemas/Course'
    Course:
      type: object
      properties:
        id:
          type: integer
          description: Course ID
        name:
          type: string
          description: Course name
        description:
          type: string
          description: Course description
```

---

# Swagger Tools

---

## Swagger UI

// TODO: Screenshot of aforementioned API

```shell
docker run -v ${PWD}:/app -e SWAGGER_JSON=/app/students.yaml -p 80:8080 swaggerapi/swagger-ui
```

---

## Swagger Editor

[editor.swagger.io](https://editor.swagger.io)

// TODO: Screenshot of aforementioned API

---

## Swagger Codegen

> The Swagger Codegen is an open source code-generator to build server stubs and client SDKs directly from a Swagger defined RESTful API.
>
> <cite><a href="https://swagger.io/docs/open-source-tools/swagger-codegen/">swagger.io/docs/open-source-tools/swagger-codegen</a></cite>

=> CLI/Docker

---

# gRPC

---

## gRPC
- RPC (Remote Procedure Call) framework
- Created by Google in 2001 ("Stubby")
- Open-sourced in 2015
- Uses HTTP/2 as a transport mechanism
  - Transport layer is abstracted away
- Uses _Protocol Buffers_ as a serialization mechanism
- _Messages_ and _services_ are defined in `.proto` files

---

## Protocol Buffers

```protobuf
// https://grpc.io/docs/what-is-grpc/introduction/

syntax = "proto3";

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply);
}
```

---

## Workflow
1. Define messages and services in `.proto` file(s)
2. Use `protoc` compiler to generate code
3. Server: implement services  
   Client: execute requests

{{% note %}}
Directly supported languages include: C++, C#, Java, Python, Ruby and Go; 3rd party addons
{{% /note %}}

---

# Comparison

---

|                      | OpenAPI                  | gRPC                        |
|----------------------|--------------------------|-----------------------------|
| Specification Format | JSON or YAML             | Protocol Buffer Language    |
| Describes            | HTTP methods + endpoints | Procedures                  |
| Contract             | Optional                 | Strict                      |
| Serialization Format | JSON*                    | Protocol Buffers*           |
| Transport Protocol   | HTTP/1.1                 | HTTP/2                      |
| Browser Support      | ✅                        | ⚠️                          |
| Streaming            | -[CHECK]                 | Server/Client/Bidirectional |
| Documentation        | Swagger UI               | e.g. protoc-gen-doc         |
| Code Generation      | Swagger Codegen          | protoc (built-in)           |

{{% note %}}
- an API described by an OpenAPI specification can be used without the JSON/YAML file
- for gRPC, the .proto file(s) are strictly required
- Serialization Format:
  - https://grpc.io/blog/grpc-with-json/
- grpc-web: https://github.com/grpc/grpc-web [CHECK]
  - gRPC requires _Trailers_, which are not implemented by browsers
  - https://news.ycombinator.com/item?id=18296014
  - https://carlmastrangelo.com/blog/why-does-grpc-insist-on-trailers
{{% /note %}}

---

# Demo

---

## Prerequisites
- Go
- protoc
- protoc plugins for Go
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

{{% note %}}
```shell
go version
protoc --version
```
{{% /note %}}
