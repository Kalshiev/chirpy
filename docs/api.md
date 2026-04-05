# Chirpy REST API
1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Content types](#content-types)
4. [Common response format](#common-response-format)
5. [Endpoints](#endpoints)
   1. [Health](#health)
   2. [Chirps](#chirps)
   3. [Users](#users)
   4. [Token management](#token-management)
   5. [Webhooks](#webhooks)
   6. [Admin endpoints](#admin-endpoints)
6. [Data models](#data-models)
7. [Notes](#notes)

## Overview

Chirpy exposes a small REST-style API for user management, chirp creation, retrieval, and token handling. The server also exposes health and admin endpoints.

Base URL: `http://localhost:8080`

> The project also serves static files under `/app/` and `/app/assets`, but the following documentation focuses on HTTP API endpoints.

---

## Authentication

### Bearer token

Most protected endpoints require a JWT access token in the `Authorization` header:

```http
Authorization: Bearer <access_token>
```

The login endpoint returns an access token and a refresh token.

### API key

The Polka webhook endpoint requires an API key in the `Authorization` header:

```http
Authorization: ApiKey <polka_key>
```

---

## Content types

- Request JSON: `application/json`
- Response JSON: `application/json`
- Admin metrics endpoint returns HTML.

---

## Common response format

Errors are returned as JSON objects:

```json
{
  "error": "error message"
}
```

---

## Endpoints

### Health

#### GET /api/healthz

- Description: Readiness probe
- Response: `200 OK`
- Response body: `OK`

#### GET /api/tea

- Description: Easter Egg
- Response: `418 I'm a teapot`
- Response body: `I'm a teapot`

---

### Chirps

#### POST /api/chirps

- Description: Create a new chirp
- Authentication: Required (`Authorization: Bearer <access_token>`)
- Request body:

```json
{
  "body": "This is a test chirp!",
  "user_id": "220a4777-ee90-49b8-8875-dcfe306d7471"
}
```

- Notes:
  - `body` must be a string of at most 140 characters.
  - The user making the chirp is determined by the authenticated JWT, not by the supplied `user_id`.

- Success response: `201 Created`
- Success body:

```json
{
  "id": "...",
  "created_at": "...",
  "updated_at": "...",
  "body": "This is a test chirp!",
  "user_id": "..."
}
```

#### GET /api/chirps

- Description: List chirps
- Query parameters:
  - `author_id` (optional): filter by author UUID
  - `sort=desc` (optional): order by `created_at` descending

- Success response: `200 OK`
- Success body: array of chirp objects

#### GET /api/chirps/{chirpID}

- Description: Get a single chirp by ID
- Path parameters:
  - `chirpID`: UUID of the chirp
- Success response: `200 OK`
- Success body: chirp object

#### DELETE /api/chirps/{chirpID}

- Description: Delete a chirp
- Authentication: Required (`Authorization: Bearer <access_token>`)
- Path parameters:
  - `chirpID`: UUID of the chirp
- Notes:
  - Only the chirp owner may delete the chirp.
- Success response: `204 No Content`

---

### Users

#### POST /api/users

- Description: Create a new user
- Request body:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

- Success response: `201 Created`
- Success body:

```json
{
  "id": "...",
  "created_at": "...",
  "updated_at": "...",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### POST /api/login

- Description: Authenticate a user and issue tokens
- Request body:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

- Success response: `200 OK`
- Success body:

```json
{
  "id": "...",
  "created_at": "...",
  "updated_at": "...",
  "email": "user@example.com",
  "token": "<jwt_access_token>",
  "refresh_token": "<refresh_token>",
  "is_chirpy_red": false
}
```

#### PUT /api/users

- Description: Update the authenticated user's email and password
- Authentication: Required (`Authorization: Bearer <access_token>`)
- Request body:

```json
{
  "email": "new-email@example.com",
  "password": "new-password"
}
```

- Success response: `200 OK`
- Success body:

```json
{
  "id": "...",
  "created_at": "...",
  "updated_at": "...",
  "email": "new-email@example.com",
  "is_chirpy_red": false
}
```

---

### Token management

#### POST /api/revoke

- Description: Revoke a refresh token
- Authentication: Required (`Authorization: Bearer <refresh_token>`)
- Success response: `204 No Content`

#### POST /api/refresh

- Description: Refresh the access token using a refresh token
- Authentication: Required (`Authorization: Bearer <refresh_token>`)
- Success response: `200 OK`
- Success body:

```json
{
  "token": "<new_access_token>"
}
```

---

### Webhooks

#### POST /api/polka/webhooks

- Description: Handle Polka webhook events
- Authentication: Required (`Authorization: ApiKey <polka_key>`)
- Request body:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "220a4777-ee90-49b8-8875-dcfe306d7471"
  }
}
```

- Notes:
  - Only `user.upgraded` events are processed.
  - The endpoint upgrades the matching user to Chirpy Red.
- Success response: `204 No Content`

---

### Admin endpoints

#### GET /admin/metrics

- Description: Returns an HTML page with the current count of served static asset requests.
- Success response: `200 OK`

#### POST /admin/reset

- Description: Reset the application state and delete all users
- Notes:
  - Only allowed when `PLATFORM=dev`.
- Success response: `200 OK`
- Failure response: `403 Forbidden` when not in dev mode.

---

## Data models

### Chirp

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "body": "string",
  "user_id": "uuid"
}
```

### User

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "string",
  "is_chirpy_red": true | false,
  "token": "<jwt_access_token>",
  "refresh_token": "<refresh_token>"
}
```

---

## Notes

- All JSON request bodies must be valid JSON.
- Errors always respond with a JSON object containing an `error` field.
- Refresh tokens are passed through the `Authorization: Bearer ...` header for the refresh and revoke endpoints.
- The static asset routes are mounted under `/app/` and `/app/assets` but are not part of the documented REST API.
