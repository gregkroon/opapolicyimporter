# README for the OPA import Tool

## Overview
This tool is designed to automate the process of fetching `.rego` policy files from a specified GitHub repository and creating corresponding policy entries in Harness. It utilizes environment variables for configuration, making it flexible and easy to integrate into different environments.

*** WARNING this is a raw MVP use at your own risk  ***

## Features
- **GitHub Integration**: Connects to a specified GitHub repository to retrieve `.rego` files.
- **Harness Integration**: Creates policies in Harness using the contents of the fetched files.
- **Environment Configuration**: Utilizes environment variables for dynamic configuration.

## Requirements
- Go programming environment (Go 1.x or higher)
- Access to a GitHub repository with `.rego` files
- Access to a Harness account with necessary permissions

## Configuration
Set the following environment variables in your environment:

- `HARNESSORG`: Harness organization identifier
- `HARNESSPROJECT`: Harness project identifier
- `HARNESSMODULE`: Harness module identifier
- `HARNESSACCOUNTID`: Harness account identifier
- `HARNESSAPIKEY`: API key for accessing Harness
- `GITHUBTOKEN`: GitHub token for accessing GitHub repositories
- `GITHUBUSER`: GitHub username
- `GITHUBREPO`: GitHub repository name

## Installation
1. Ensure Go is installed and properly configured on your system.
2. Clone this repository or download the source files into a local directory.
3. Set the required environment variables as described in the configuration section.

## Usage
Execute the program by running the following command from the terminal:

```bash
go run main.go
