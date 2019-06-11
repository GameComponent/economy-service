[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGameComponent%2Feconomy-service.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FGameComponent%2Feconomy-service?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/GameComponent/economy-service)](https://goreportcard.com/report/github.com/GameComponent/economy-service)

**State of the project:** Unstable and in active development

# :dollar: Economy Service

The economy service is a service to manage your game's economy.

- [What is the economy service?](#what-is-the-economy-service)
- [Requirements](#requirements)
- [Setup](#setup)
- [Contributing](#contributing)
- [License](#license)

## What is the economy service?

The economy service allows you to give your players access to items, currencies and much more. Anything that touches the economy of your game is parts of this service. This services allows you to define item definitions, create currencies, open up shops and much more.

## Requirements
Make sure you have `go` (atleast version 1.11), `make`, `protoc`, `docker` and `docker-compose` installed on your system.

## Setup
1. `make api`. This will generate Go bindings, a REST gateway and a Swagger JSON document.
2. `make build`. This will build the Go project.
3. `docker-compose up -d`. This will run a single node CockroachDB database.
4. `./bin/server/server`. Run the server. Too see all available arguments run `./bin/server/server --help`.

## Contributing

We're an open source project and welcome contributions. Read our [Contribution guidelines](CONTRIBUTING.md) for more information


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FGameComponent%2Feconomy-service.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FGameComponent%2Feconomy-service?ref=badge_large)
