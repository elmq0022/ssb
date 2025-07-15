# Features

- Publish article via CLI with bearer token authentication
- Sever articles publically via HTTPS
- Render Markdown as HTML for display
- Deployable on fly.io
- Lightweight and cheap to run


# Architecture

- backend: go
- database: sqlite
- deploy on: fly.io
- CLI publisher: send artile to backend via CLI

# Authentication

- Bearer Token Authentication for /post
- Token is static
- Token is set on fly.io as a env value


# API Design

## POST ARTICLE
- URL
  - /api/vi/article

- Auth Required: yes - bearer token
- Request body:
{
  "title": "a title",
  "author": "first and last name"
  "body": "the article"
}

- successful response: 201
{
 "id": "the new id"
}

- unauthorized: 401
-

## PUT ARTICLE
- URL
  - /api/v1/article/{id}
- replace an existing resource
- Authentication yes: Bearer Token
- same fields as POST

Successful Response

status: 200
{
  "id": "abc123",
  "title": "Updated Title",
  ...
}

Error Responses
- 400 Bad Request – invalid JSON or missing required fields
- 401 Unauthorized – missing or invalid bearer token
- 403 Forbidden – if token is valid but not authorized
- 404 Not Found – if the article with the given ID doesn't exist

## DELETE ARTICLE
DELETE /api/v1/article/{id}

Successful Responses
status: 204
No body

Error Responses
- 401 Unauthorized
- 403 Forbidden
- 404 Not Found


## GET ARTICLE
- URL
  - /api/v1/article/{id}

# Data Model

| field        | type      |
|++++++++++++++|+++++++++++|
| id           | string    |
| title        | string    |
| author       | string    |
| title        | string    |
| body         | string    |
| published_at | timestamp |
| updated_at   | timestamp |

# Deployment
- use docker to build go app
- deploy on fly.io using flyctl
- set bearer token as fly.io secret

# CLI Publisher

- go commandline tool
- read markdown file from disk
- send POST to /api/v1/article with token
- Optional config file or env vars for:
  - Server URL
  - Auth token
  - Default author


