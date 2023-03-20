# tracking-service 

### Where to start
- Install go onto your machine (https://golang.org/doc/install)
- clone the repo
- setup the env files from the env.example ( if you don't have the env files, ask the team for them, the app won't work without them )
- run `go mod download` to download all the dependencies
- run `go run main.go` to start the server


### Application Structure
- `main.go` is the entry point of the application
- `aft` contains africas talking related code
- `dto` contains all the data transfer objects, going to be used throughout the application since
- `handlers` contains all the handlers for the application
- `models` contains all the models for the application
- `repository` the repository, is what brings eveything together, i.e the db clients, both mongo and postgres, as well as the africas talking client. It's also where the application routes are setup.
- `storage` this folder contains the logic for initializing connections to the postgres and mongo clients
- `utils` this folder contains all the utility functions for the application


#### Reviewing the code?
- Start from the `main.go` file, this is the entry point of the application