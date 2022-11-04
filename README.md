# Budget 2 Sheets
![Go](https://github.com/jbleduigou/budget2sheets/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbleduigou/budget2sheets)](https://goreportcard.com/report/github.com/jbleduigou/budget2sheets)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=jbleduigou/budget2sheets)](https://dependabot.com)


I created this project for automating some tasks I was doing manually when taking care of my personal finances.  
Long story short, this is an ETL for tracking my expenses.  
The project consists of two lambda functions : [budgetcategorizer](https://github.com/jbleduigou/budgetcategorizer) and [budget2sheets](https://github.com/jbleduigou/budget2sheets)  
The transform part is performed by **budgetcategorizer** and has two main responsibilities : sanitize the expense description and assign a category to the expense.  
The load part is performed by **budget2sheets** which is going to upload the transactions (i.e. expenses) to Google Sheets.

## Overall Architecture

![Architecture Diagram](architecture_diagram.png)

## Getting Started

Clone the repo inside the following directory:

```bash
~/go/src/github.com/jbleduigou/

```

If you want to fork the repo, replace the latest path element with your GitHub handle.

### Prerequisites

You will need to have Go installed on your machine.  
Currently the project uses version 1.19

### Building
You will find a Makefile at the root of the project.  
To run the full build and have zip file ready for AWS Lambda use:

```bash
make zip
```

If you only want to run the unit tests:

```bash
make test
```

## Deployment

For now deployment is made manually.  
It would be nice to have a cloudformation template at some point.

## Improvements / Remaining Work

* extract logic to dedicated classes and write unit tests
* improve error handling
* create cloud formation template
* fix weird behaviour for GitHub Actions upload artifact : https://github.com/actions/upload-artifact/issues/39 ?

## Configuration

Configuration in this application is based on the concepts exposed by [The Twelve Factor App](https://12factor.net/).  
The idea is to strictly separate config from code by using environment variables.  
Please read the page on [12 Factor Configuration](https://12factor.net/config) for more details.

### Environment Variables

The following environment variables should be declared within you lambda:

| Name                         | Description                                 | Sample Value                                                 |
| ---------------------------- |:-------------------------------------------:| :-----------------------------------------------------------:|
| GOOGLE_ACCESS_TOKEN          | OAuth 2.0 to Access Google APIs             | Swac9MxkndN1elrl3y7Gk6XWizKC97gs48eJ3p7O                     |
| GOOGLE_CLIENT_ID             | OAuth 2.0 to Access Google APIs             | 769149424942-sUsY928FXCeQB15Dpot0.apps.googleusercontent.com |
| GOOGLE_CLIENT_SECRET         | OAuth 2.0 to Access Google APIs             | Swac9MxkndN1elrl3y7G                                         |
| GOOGLE_PROJECT_ID            | OAuth 2.0 to Access Google APIs             | budget2sheets-633863                                         |
| GOOGLE_REFRESH_TOKEN         | OAuth 2.0 to Access Google APIs             | XgukkG2720EZOZuiYiHZ                                         |
| GOOGLE_SPREADSHEET_ID        | ID of spreadsheet to be populated           | 6JIJnDjHKcXwF9EhCTPLp0BJiZ03tMPvbVn36sj4                     |
| GOOGLE_SPREADSHEET_RANGE     | Range of cells where data should be appended| Suivi Dépenses Janvier!A2:F2                                 |

## Project Structure

The project structure was inspired by the project [Go DDD](https://github.com/marcusolsson/goddd).  
The entry point is located in folder cmd/budgetcategorizer.  
What it does is instantiating all the dependencies for the command.  
The business logic was separated by concerns and placed in dedicated folders.  
Interfaces were introduced to avoid tight coupling and therefore facilitate unit testing (amongst other benefits).  

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

* **Jean-Baptiste Le Duigou** - *Initial work* - [jbleduigou](https://github.com/jbleduigou)

See also the list of [contributors](https://github.com/jbleduigou/budget2sheets/contributors) who participated in this project.

## License

Licensed under the Apache License, Version 2.0.  
See [LICENSE.txt](LICENSE.txt) for more details.  
Copyright 2020 Jean-Baptiste Le Duigou
