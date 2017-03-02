# BaconCobra API

Rest API to work against BaconCobra storage 

## Build status

* master:  [![CircleCI Master](https://circleci.com/gh/baconcobra/api/tree/master.svg?style=svg)](https://circleci.com/gh/baconcobra/api/tree/master)
* develop: [![CircleCI Develop](https://circleci.com/gh/baconcobra/api/tree/develop.svg?style=svg)](https://circleci.com/gh/baconcobra/api/tree/develop)

## Installation

```
make deps
make install
```

## Running Tests

```
make deps
make test
```

## Authentication

Authentication is handled by JWT. You must first authenticate via `/auth/` and use the returned web token as a header in all subsequent requests.

```
curl -i -X POST -d "username=something" -d "password=something" https://bc_instance/auth/
```

This will return the following json payload:

```json
{"token":"VALID-AUTH-TOKEN"}
```

This then can be used in subsequent requests, like so:

```
curl -i -H 'Authorization: Bearer VALID-AUTH-TOKEN' https://bc_instance/api/users/
```


## Endpoints

Supported endpoints are Tags, Actors, Videos and Tubes.


## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2016 supu.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
