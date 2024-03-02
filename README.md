# Chat API Server: Golang

This repo implements the same chat-server api spec as the others in the parent project repo: [Chat](https://github.com/MikeJeffers/chat)

This is best operated as a gitsubmodule to that project.  Please follow the link and instructions for cloning the parent repo.

## API Spec:
The Chat API server is primarily for generating and managing the identity of users.
A user may register an account (username and password), login with those credentials, and receive a JWT with which the chat client can connect to the websocket chat servers.
The Chat API server must implement the following routes:
 - POST /register application/json body `{username:string, password:string}`
   - 201 Response: body `{token:string, user:{id:number, username:string}}`
   - 4xx Response for any input data issues, non-uniqueness, validation
   - 5xx Response for any server errors
 - POST /login application/json body `{username:string, password:string}`
   - 200 Response: body `{token:string, user:{id:number, username:string}}`
   - 4xx Response for any input data issues, non-uniqueness, validation
   - 5xx Response for any server errors


## Local setup (non-docker)
Running this repo locally still requires the dependent services be upped in the parent repo.  The `./run.sh` script assumes this relationship.

### Install
```sh
./scripts/install.sh
```
Prep .env file
```sh
./scripts/init.sh
```
### Run
```sh
./scripts/run.sh
```
`ctrl-c` to stop the process and use `./scripts/deps-down.sh` to down and clear volumes