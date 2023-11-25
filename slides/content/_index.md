+++
title = "API Design: OpenAPI vs. gRPC"
outputs = ["Reveal"]
+++

## API Design
# OpenAPI vs. gRPC

Simon Hammes

Cloud Computing WS 2023/24

---

# OpenAPI/Swagger

---

## OpenAPI/Swagger
- created in 2011 by Tony Tam
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

// TODO

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

<!-- TODO -->

---

## Code

```go{}
package main

import "fmt"

func main() {
    fmt.Println("hello world")
}
```
