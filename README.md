[![codecov](https://codecov.io/gh/gregbiv/news-api/branch/master/graph/badge.svg)](https://codecov.io/gh/gregbiv/news-api)
[![Build Status](https://travis-ci.org/gregbiv/news-api.svg?branch=master)](https://travis-ci.org/gregbiv/news-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/gregbiv/news-api)](https://goreportcard.com/report/github.com/gregbiv/news-api)

REST api for [News-Today](https://github.com/gregbiv/news-today).

- [Maintainers](#maintainers)
- [API Docs](#api-docs)
- [Development Environment](#development-environment)
	- [pre-commit](#pre-commit)
		- [Pre-commit error: Files were modified by this hook](#pre-commit-error-files-were-modified-by-this-hook)
- [Database Migrations](#database-migrations)
- [Test Suites](#test-suites)
	- [Unit tests (testing the code)](#unit-tests-testing-the-code)

## Maintainers

* Gregory Kornienko [@gregbiv](https://github.com/gregbiv)

[[table of contents]](#table-of-contents)

## API Docs

We document our API in [RAML](https://raml.org) format.

Please see the generated result [here](https://news-api.kornienko.site/docs/api.html)

RAML sources are located at `resources/docs` directory.
Feel free to edit them or notify the [maintainers](#maintainers) if you've found an error.

[[table of contents]](#table-of-contents)

## Development Environment

This project requires git, Go and Docker.

As soon as your development environment meets the aforementioned requirements, run the following commands:

1. `git clone git@github.com:gregbiv/news-api.git`
2. `cd news-api`
3. `make && make deps-dev`
4. `docker-compose up -d`
5. `make migrations-dev`

This repository has a set of automatic checks before any pull request is merged.
Make sure to read the remaining topics in this section to be able to contribute / commit to the repository.

[[table of contents]](#table-of-contents)

### pre-commit

In order to make sure that your code meets our requirements, please install `pre-commit`:
`brew install pre-commit`

Then `cd` to the project's directory, run `git add .` and then run `pre-commit run`.

Fix the issues according to an output of the program, **stage all the changes again** and then run `pre-commit` to see results.

Repeat the process until you fix all issues.

_Please note that pre-commit runs for **staged** changes only, so don't forget to `git add .` before you run it._

[[table of contents]](#table-of-contents)

#### Pre-commit error: Files were modified by this hook

The error message tells you that one of the hooks has modified a file, so most likely you have to stage everything and commit again.

For example, we have a hook which automatically generates a table of contents for the `README.md` file.
If you forgot to generate it manually, the hook will generate it for you and you'll get the aforementioned error because your README has been modified by the hook and the latest version is not staged.

Solution: `git add . && git commit -m "your message""`

[[table of contents]](#table-of-contents)

## Database Migrations

You can run database migrations as follows:

`make migrations-dev`

[[table of contents]](#table-of-contents)

## Test Suites

In this project you can test the code and the automation scripts.

### Unit tests (testing the code)

You can simply run `make test-unit`.

If you don't want to upload coverage statistics to Codecov, run: `make test-dev`.

[[table of contents]](#table-of-contents)
