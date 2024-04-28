# Test task BackDevTask

-   [Usage](#usage)
-   [Description](#description)

## Description

There are 3 endpoints:

1. /register, creates user with given guid, user has additional session object with refresh token, it's epmty at first, later to be hashed and stored in MongoDB
2. /authenticate, creates access and refresh tokens with given guid, and returns tokens
3. /refresh, generates new pair of token for given refresh token, and updates session object, and returns tokens

## Usage

Run by `go run main.go`
Rename [.env.example](.env.example) or create new .env file and change accrodingly to [.env.example](.env.example).

> **Swagger is accesible at /swagger**

[![swagger page](https://i.imgur.com/2nAGisL.png "swagger page")](https://i.imgur.com/2nAGisL.png "swagger page")
