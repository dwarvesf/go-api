# Project Name

Short description of your project.

## Table of Contents

- [About](#about)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Running the Development Server](#running-the-development-server)
  - [Generating Swagger Documentation](#generating-swagger-documentation)
  - [Generating Mock Interfaces](#generating-mock-interfaces)
  - [Running Tests](#running-tests)
- [Database](#database)
  - [Starting the Development Database](#starting-the-development-database)
  - [Starting the Testing Database](#starting-the-testing-database)
  - [Migrations](#migrations)
- [Contributing](#contributing)
- [License](#license)

## About

Briefly describe your project, its purpose, and its main features.

## Getting Started

### Prerequisites

- Go (add specific version)
- Docker (add specific version)
- ...

### Installation

1. Clone the repository:

   ```bash
   git clone git@github.com:dwarvesf/go-api.git new-project
   cd new-project
   ```

2. Install project dependencies:

   ```bash
   make setup
   ```

## Usage

### Running the Development Server

Start the development server:

```bash
make dev
```

### Generating Swagger Documentation

Generate Swagger documentation:

```bash
make gen-swagger
```

### Generating Mock Interfaces

Generate mock interfaces:

```bash
make gen-mocks
```

### Running Tests

Run tests with the testing database:

```bash
make test
```

## Database

### Starting the Development Database

Start the development database container:

```bash
make pg-start-dev
```

### Starting the Testing Database

Start the testing database container:

```bash
make pg-start-test
```

### Migrations

Apply pending migrations:

```bash
make pg-migrate-up
```

Rollback the last migration:

```bash
make pg-migrate-down
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## License

This project is licensed under the [License Name](LICENSE).