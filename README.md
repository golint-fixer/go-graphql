# gographql
[![Build Status](https://travis-ci.org/kmulvey/gographql.svg?branch=master)](https://travis-ci.org/kmulvey/gographql)
[![Coverage Status](https://coveralls.io/repos/github/kmulvey/gographql/badge.svg?branch=master)](https://coveralls.io/github/kmulvey/gographql?branch=master)

Generate a graphql schema in Go from an existing sql database.

This is currently alpha software that "works" for mysql but thats about it, [help is welcome](https://github.com/kmulvey/gographql/issues).


## Usage

`sql2graphql [options]`

### Options:
  
  - `--output`    - Directory to use when generating code *`(string [required])`*
  - `--schema`    - Schema name *`(string [required])`*
  - `--hostname`  - Hostname of database server *`(string [default: "localhost"])`*
  - `--port`      - Port number of database server *`(number)`*
  - `--username`  - Username to use when connecting *`(string [default: "root"])`*
  - `--password`  - Password to use when connecting *`(string [default: ""])`*
