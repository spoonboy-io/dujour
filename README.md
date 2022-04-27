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
consumable API is a requirement.

### Features

- Automatic self-signed TLS certificate (or use your own)
- Data is served from memory cache
- Supports any number of *.json and *.csv files
- Hot reload, new or edited data, with no restarts