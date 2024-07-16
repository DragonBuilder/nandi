# Nandi

Authentication Server

## Goal

Learn authentication strategies. 

Web apps are incredibly powerful. Such products needn't be installed locally and hence aren't generally constrained by the capabilities of the client machines. 

Since web apps are hosted on the internet they are generally accessible by anyone with an internet connection. The need to identify users and authenticate users is very important so as to avoid unauthorized access and also to tailor custom experiences. Think x.com feeding you tweets from the accounts you follow and not a general feed or Netflix recommending movies based on your preferences. 

So how do we trust the users to tell us who they are? Well that's where authentication & authorisation can help. This repository is a experimentation to explore various strategies for the same.

## Dependencies

Since the goal is to learn the insides of various authentication and authorisation strategies I will be keeping the dependencies to a bare minimum. 

1. Go language
2. sqlite3. Using [go-sqlite3](https://github.com/mattn/go-sqlite3) as the Go-SQLite driver. 
   
    The initial `go get` of this driver would need `gcc` locally installed and the `CGO_ENABLED` environment variable set to `1`, like `export CGO_ENABLED=1`. 
   
   Eg :- `export CGO_ENABLED=1 && go get github.com/mattn/go-sqlite3@v1.14`
3. [reflex](https://github.com/cespare/reflex), if you want to run `./watch.sh` which will reload the app whenever source code is edited.

## Running the code

Simply run `./start.sh` to start the app or run `./watch.sh` to start and also watch for changes in source code and automatically reload the app. 

## HTTP Session based authentication ([v0.1](#v01))

The user provides a unique ID and a password. Here I am using `email` and `password` but it could might as well have been `username` or `ID` or anything. As long as it's unique across the app it should be able to identify a user correctly.

We verify the `email` and `password`. If they are valid, we save the user info in HTTP session cookie so that they don't have to provide the email and password with every subsequent request. 

Any URL that is supposed to be used by an autheticated user can now check for the presence of the the user info in the session cookies and take the appropriate action. 

### Sign Up
The ability for a user to sign-up is important if you want to provide custom experiences for users. The simplest sign-up flow would require an email & a password to uniquely identify and authenticate the user in the future.

A HTTP `POST` endpoint `/api/user/sign-up` is created for anyone to start signing up. 

A handler `SignUpHandler` is defined in [handlers.go](handlers.go) to handle the requests on this endpoint. 

The `SignUpHandler` expects the request body to be a JSON and conforming to the type `UserSignUp` as defined in [handlers.go](handlers.go).

One the body is decoded into an instance of `UserSignUp`, some basic validations such as making sure none of the fields are empty and such. 

Then we generate the hash of the `Password` field and also compare it with value at `ConfirmPassword` field. This is to ensure that the user didn't accidently make a typo while entering the password. 

The details are then stored in the `users` table. We store the `email` and the `password` with `created_at` and `updated_at` because why not? The important thing to note here that we **don't store the password in it plain text** format instead we **store the hashed value of the password**. This is generally considered a good practice. I think there might be couple of reasons for it. 

1. In case of a data breach, the passwords are not automatically compromised.
2. The organisation themselves aren't aware of the actual password.

Now the user should be able to login.

## Login
TODO

## Logout
TODO

## Reveal Secret URL
TODO

## Roadmap

### v0.1
- [x] User sign-up
- [x] User login
- [x] An authenticated URI that can only be accessed if the user is logged in.
- [x] User logout
- [ ] Add unit tests

### v0.2
- [ ] Let users authenticate with their Nandi account from other products and services.
- [ ] Let other services register themselves as Nandi clients.
- [ ] Let other services authenticate with Nandi.

