# Todist API

A simple unofficial [Todoist](https://todoist.com/) REST API Client

This is a work in progress and not ready for production usage yet.

## Installing

Requires Go modules. You can `go get` the module and add it to your go.mod:

`go get github.com/treelightsoftware/go-todoist`

## Usage

To get started, you should be familiar with how Todoist organizes its entities in their [docs](https://developer.todoist.com/rest/v1/#overview).

You will need to get the authentication token for the user. This can be done either through an oAuth flow OR you can pass it through the environment. This is a lsight nuance to this library, as most calls have a `token` parameter are the first for the functions. If you pass in an auth token AND leave that field blank, the call will automatically use the environment's passed in auth token. So, for example, to get a list of projects for a user, you could do the following:

```go
import todo "github.com/treelightsoftware/go-todoist"

// this will get all of the projects for the user with that token; useful if you have multiple users on your integration
projects, err := todo.GetAllProjects("user_token")
```

Alternatively, if you are using this SDK in a project that involves only a single user (such as a CLI or local application), and you know the auth token ahead of time (such as from the [integrations](https://todoist.com/prefs/integrations) site), you can pass it into the environment:

`TODOIST_AUTH_TOKEN=the_token ./your_app`

This will set the `config.AuthToken` field at startup, and you would instead make your calls in either of the following ways:

```go
import todo "github.com/treelightsoftware/go-todoist"

// explicitly grab it from the config
projects, err := todo.GetAllProjects(config.AuthToken)

// OR

// implicitly grab it from the config, filled in during the makeCall client function invocation
projects, err := todo.GetAllProjects("")

```

The pattern persists throughout.

### Why pointers for the fields of the params?

The default values have meaning in the Todoist API. In otherwords, if you try to update a task and set the content, but not the description field,
the description will default to "", so it will be set to blank without the user's intention. For further inspiration, see the excellent [Stripe-Go](https://github.com/stripe/stripe-go/) library.

## Contributing

Contributors are welcome. You should raise an issue or communicate with us prior to committing any significant effort to ensure that your desired changes are compatible with where we want this library to go. Read more in the CONTRIBUTING.md document. Make sure your tests pass (you will need to provide your own authentication token to get ideal coverage).

## TODO

The following API end points are provided by v1 of the Todoist API. If it has an X, it's been implemented:

- Projects
  - [X] Get All Projects
  - [X] Create a New Project
  - [X] Get a Project
  - [X] Update a Project
  - [X] Delete a Project
  - [ ] Get All Collaborators
- Sections
  - [X] Get All Sections
  - [X] Create a New Section
  - [X] Get a Single Section
  - [X] Update a Section
  - [X] Delete a Section
- Tasks
  - [X] Get Active Tasks
  - [X] Create a New Task
  - [X] Get an Active Task
  - [X] Update a Task
  - [X] Close a Task
  - [X] Reopen a Task
  - [X] Delete a Task
- Comments
  - [ ] Get All Comments
  - [ ] Create a Comment
  - [ ] Get a Comment
  - [ ] Update a Comment
  - [ ] Delete a Comment
- Labels
  - [X] Get All Labels
  - [X] Create a New Label
  - [X] Get a Label
  - [X] Update a Label
  - [X] Delete a Label
- Webhooks
  - [ ] Add function that takes HTTP request, parses it, and returns relevant information to the user

### Other TODOs

- [ ] Implement code coverage with dummy token in CI/CD
- [ ] Improve error checking; for example, if a task's DueDatetime is set, make sure it includes the time component
- [ ] Provide better documentation on usage
