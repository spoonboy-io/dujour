# Dujour

## A JSON/CSV Data File Server

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/spoonboy-io/dujour?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/spoonboy-io/dujour?style=flat-square)](https://goreportcard.com/report/github.com/spoonboy-io/dujour)
[![DeepSource](https://deepsource.io/gh/spoonboy-io/dujour.svg/?label=active+issues&token=uYY_4Kwjq9MnjT7TzykEyv-J)](https://deepsource.io/gh/spoonboy-io/dujour/?ref=repository-badge)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/spoonboy-io/dujour/Build?style=flat-square)](https://github.com/spoonboy-io/dujour/actions/workflows/build.yml)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/spoonboy-io/dujour/Unit%20Test/master?label=tests&style=flat-square)](https://github.com/spoonboy-io/dujour/actions/workflows/unit_test.yml)

[![GitHub Release Date](https://img.shields.io/github/release-date/spoonboy-io/dujour?style=flat-square)](https://github.com/spoonboy-io/dujour/releases)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/spoonboy-io/dujour/latest?style=flat-square)](https://github.com/spoonboy-io/dujour/commits)
[![GitHub](https://img.shields.io/github/license/spoonboy-io/dujour?label=license&style=flat-square)](LICENSE)

## About

Dujour is a JSON/CSV data file server. It supports any usecase in which simple transformation of JSON and CSV data to a web 
consumable JSON REST API is a requirement.

## Releases

You can find the [latest software here](https://github.com/spoonboy-io/dujour/releases/latest).

### Features

- Automatic self-signed TLS certificate (or use your own)
- Supports CSV files. Application will parse them to JSON
- Supports any number of JSON or CSV data files, memory being the only constraint
- Hot reload. New or edited data can be added with no server restart needed
- Data is served from memory. Fast

### Usage
Add `.json` and `.csv` data files to the `data` directory and Dujour will automatically load, validate and serve each data file at two REST API endpoints in JSON format.

In each data file, element/row data should contain an `id` key/column which should be unique in the dataset.

Data is loaded and served from an in-memory cache. No restart of the server is required when adding new data. Adding a new file of same name will cause the cache to be cleared and the data reloaded.

The filename of each data file determines the API endpoints which are created. For example, for a file named `users.json` (not case sensitive), Dujour will serve data at two endpoints:-

#### Get all users
This endpoint will retrieve all users:
```
GET $serverUrl:18651/users
```

#### Get a specific user
This endpoint will retrieve a specific user:
```
GET $serverUrl:18651/users/$id
```

### Installation
Grab the tar.gz or zip archive for your OS from the [releases page](https://github.com/spoonboy-io/dujour/releases/latest).

Unpack it to the target host, and then start the server.

```
./dujour
```

Or with nohup..

```
nohup ./dujour &
```

To update the application server, stop the server, replace the binary, then start the server.

### Limitations

- Dujour does not perform mutations on the data files. Only `GET` operations are supported.
- To  

### Development Opportunities

- Support CRUD operations on the data
- Run the application as a service

### License
Licensed under [Mozilla Public License 2.0](LICENSE)
