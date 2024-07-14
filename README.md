# go-auth

## Overview

`go-auth` is a minimalistic Go-based authentication system designed to manage user registration, login, and secure session handling. It includes features for creating and verifying user credentials, session management and email verification.

It gives you the chance to choose from one of three databases: file based db, Postgres, MongoDb, and Redis.

you can set what database to use in the config/config.go file, by setting the DBMS constant to
one of file, postgres or mongo and TempStoreType to either redis or file. 
Using other than file as DBMS or TempStoreType requires you to have those databases running.

For easy setup and testing, use the file based database as it requires no installations.

## Features

- User Registration
- Email Verification before user Registration (toggle at config.RequireEmailVerification)
- User Login
- Session Management

## Installation

Clone the repository:
```sh
git clone https://github.com/Iyusuf40/goBackendUtils.git
```

## Navigate to the project directory:

```sh 
cd go-auth
```

## Install dependencies:

```sh 
go mod tidy
```

# Usage

## Setup
- Write your user model in models/User.go. Must have Email and Password fields.
- If you have setup Postgres as your DBMS, update the userSchema variable in storage/UsersStorage.go to reflect your schema.
- Checkout config/config.go to set your preferences. You can choose your database engine, setup credentials, and customize your welcome message.

## Run the application:

```sh 
go run main.go
```

## Routes

### Users Routes
-	POST    "api/users"         PAYLOAD     {data: {userPayload}}
-	GET     "api/users/:id"
-	PUT     "api/users/:id"     PAYLOAD     {data: {field, value}}
-	DELETE  "api/users/:id"
-   GET     "/complete_signup/:signupId"

### AUTH Routes
-   POST    "auth/login"        PAYLOAD     {data: {email, password}}
-   POST    "auth/logout"       PAYLOAD     {data: {sessionId}}
-   POST    "auth/isloggedin"   PAYLOAD     {data: {sessionId}}




# Contributing

Feel free to submit issues or pull requests for improvements and bug fixes.

# License

This project is licensed under the MIT License.