# Simple IdP

A quick try on Identity Provider just because OpenStack's Keystone is too much. This IdP meant to be used as a micro-service for domains (tenants), users, RBAC (NIST Level 1, non-hierarchical) and sessions.

Currently Simple IdP supports token-based authentication over RESTful API. It does not implement SSL as it is intended to be used behind the proxy/balancer.

The package contains 2 commands (executables):

 * `idp-cli` - a command-line utility to view/manage all entities
 * `idp-api` - a server that exposes functionality via RESTful API (for the moment)


## Database support & data model

Currently the Simple IdP supports the following RDBMS via standard Go's `database/sql` interface:

 * MySQL (http://github.com/go-sql-driver/mysql)

The following are WORK IN PROGRESS:

 * SQLite3 (http://github.com/mattn/go-sqlite3)
 * PostgreSQL (http://github.com/lib/pq)

## Building

You can use either included `Makefile` or simple run the following commands:

    go install github.com/oleksandr/idp/cmd/idp-api
    go install github.com/oleksandr/idp/cmd/idp-cli

The corresponding binaries will be created in your `$GOPATH/bin` directory.


## Configuration

Following the 12-Factor-App methodology (http://12factor.net/) the command line tool (`idp-cli`) and a service itself (`idp-api`) are configured via environment variables. 

 * `IDP_REST_ADDR` - an address/port to bind HTTP server to (e.g. `0.0.0.0:8000`)
 * `IDP_RPC_ADDR` - an address/port to bind Thrift RPC server to (e.g. `0.0.0.0:8001`)
 * `IDP_SESSION_TTL` - session TTL in minutes (e.g. `30`)
 * `IDP_SECRET_SALT` - password hashing secret salt (set once before deployment)
 * `IDP_DB_Driver` - name of the database driver to use (e.g. `mysql`, `postgres`, `sqlite3`)
 * `IDP_DB_DSN` - connection DSN, which format depends on a specific driver.
 * `IDP_SQL_TRACE` - dump SQLs into log (`true`/`false`, default `false`)

You can see example of configuration in the included `env.sh` file.


# Running the API

    $ source env.sh
    $ idp-api
    [main] 2015/04/02 11:54:47 RESTful API Server listening 127.0.0.1:8000
    [main] 2015/04/02 11:54:47 RPC API Server listening 127.0.0.1:8001


## Using CLI

    $ source env.sh
    $ idp-cli -h


## RESTful API

For the moment the following resources and methods are available.

### Sessions

 * POST /v1/sessions
 * GET /v1/sessions/current
 * HEAD /v1/sessions/current
 * DELETE /v1/sessions/current

Creating a session requires posting the following structure:

    {
      "session": {
        "domain": {
          "name": "domain1.com"
        },
        "user": {
          "name": "user1",
          "password": "pass1"
        }
      }
    }

### RBAC

 * HEAD /assert/role/`rolename`
 * HEAD /assert/permission/`permissionname`

As alternative you can use `session.domain.id` instead of a domain's name.

## Apache Thrift API

See `spec/services.thrift` for the services you can consume. Use this file to generate clients for the programming language of your choice.

## Authentication

You need to include the following header in your HTTP request:

    Authorization:"Token token=c25b0ff5-a35c-4f63-8ffa-b218771ad365"

where, `c25b0ff5-a35c-4f63-8ffa-b218771ad365` is a token (session's identifier) receiver after successful session creation (see Example below).

Or you can use `X-Auth-Token` header as show below:

    X-Auth-Token: c25b0ff5-a35c-4f63-8ffa-b218771ad365


## Example

The package includes `test_bootstrap.sh` and `test_login.json` files. The first one after some modification in the header can be used to populate database with various test data (domains, users, roles, permissions). 
The second file can be used for creating a new session via RESTful API. Here's an example of interacting with API using HTTPie client (https://github.com/jakubroztocil/httpie):

    $ http :8000/v1/sessions < test_login.json
    HTTP/1.1 201 Created
    Content-Length: 371
    Content-Type: application/json; charset=utf-8
    Date: Fri, 27 Mar 2015 07:56:55 GMT
    {
        "session": {
            "created_on": "2015-03-27T07:56:55Z",
            "domain": {
                "description": "Test domain #1",
                "enabled": true,
                "id": "48981dda-4dac-4cad-bf99-71e268da5fb5",
                "name": "domain1.com"
            },
            "expires_on": "2015-03-27T08:56:55Z",
            "id": "c25b0ff5-a35c-4f63-8ffa-b218771ad365",
            "updated_on": "2015-03-27T07:56:55Z",
            "user": {
                "enabled": true,
                "id": "4d591a87-e051-4d92-8fbb-f7aa0e0a42ca",
                "name": "user1"
            }
        }
    }

Checking existing session:

    $ http head :8000/v1/sessions/current Authorization:"Token token=c25b0ff5-a35c-4f63-8ffa-b218771ad365"
    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Date: Fri, 27 Mar 2015 08:00:52 GMT

If the session is not valid the response will be:

    HTTP/1.1 401 Unauthorized


## Dealing with date and time

The code takes current time in UTC and stores it in database without a timezone. The date and time returned in responses is UTC.


