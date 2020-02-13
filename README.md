# Case study on API development

In this exercise you will authenticate a user with credentials and provide them with a session to perform other actions.

## Part 1 - Authentication

Create a login HTTP handler using the code in `api/login_handler.go`.

The handler func has OpenAPI annotations that describe the expected responses and required parameters.

- All responses in the endpoint description must be handled.
- It must validate the JSON payload.
- The successful response must generate a JWT with a session ID in an HTTP cookie and/or response body.
- The HTTP cookie name should be `invision_jwt` and expire in 30 minutes.
- The JWT must have a `session_id` value that grants the user access to teams (using _Test JWT_ info below).

Verify that your handler works by executing:

```bash
$ go test ./...
```

### Test user

> User ID: 1

> Email: "user@example.com"

> Password: "electric cat festival"

### Test JWT

> Issuer: "InVision"

> Session ID: "56906974-f924-4a1c-889e-a1dd2f395ac2"

> Expiration: 1 day


## Part 2 - Authorization

Write a client to use the JWT from Part 1 to access team info.

You must run the API server:

```bash
$ go run cmd/server/server.go
```

Your API client will run from:

```bash
$ go run cmd/client/client.go
```

1. The client will login using the credentials of the _Test user_ and get a JWT.
2. Get the info of team ID *ch72gsb320000000000000000*.
3. Print the JSON response to stdout.


## Completion

Submit a PR to this repo merging into master with an explanation of your changes.
