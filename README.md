# Chirpy
This is a learning project used to learn about Authentication, Authorization, RESTful APIs and HTTP servers in Go.
- This is a guided project that is part of the Backend Developer Path on boot.dev

## Motivation
By designing and building a mockup of X (Formerly twitter) I had the chance to learn about:
- The fundamentals of HTTP servers
- Production-style http servers in Go, without the use of a framework
- JSON, headers, and status codes to communicate with clients via RESTful API
- Type safe SQL to store and retrieve data from a Postgres database
- Secure authorization/authentication systems with well-tested cryptography libraries
- Webhooks and API keys
- Documentation

## Requirements
1. Postgres Database
2. Go Lang

## Installation
1. Git clone this repo
```sh
git clone https://github.com/Kalshiev/chirpy
```
2. Install chirpy
```sh
cd chirpy
go install
```
3. Generate a .env file with the following fields
```md
DB_URL="postgres://..."
PLATFORM=""
SECRET_KEY=""
POLKA_KEY=""
```
4. Run
```sh
./chirpy
```

## API Documentation
Chirpy exposes a small REST-style API for user management, chirp creation, retrieval, and token handling. The server also exposes health and admin endpoints.

[Click here](docs/api.md) to access the complete documentation.