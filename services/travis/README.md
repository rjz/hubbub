# Travis CI integration

Describe policies for Travis CI repository settings

## Configuration

To use the travis-ci integration, [obtain an access token][travis-token] for your service and add it to your environment.

For Travis free (travis-ci.org),

    $ export HUBBUB_TRAVIS_ORG_TOKEN=<your token>

For Travis Pro (travis-ci.com),

    $ export HUBBUB_TRAVIS_PRO_TOKEN=<your token>

For policies applied across a mix of public and private repositories, simply
generate and set both tokens!

## Goals

### `travis_env_var`

Define Travis environment variables ([API
documentation](https://docs.travis-ci.com/api/#settings:-environment-variables)).

#### Parameters

  key     | type     | description
  ------- | -------- | ----------------------------------
  `state` | `string` | one of `"absent"` OR `"present"`
  `name`  | `string` | the name of the variable to set
  `value` | `string` | (optional) the variable's value

### `travis_repository_settings`

Update Travis repository settings ([API documentation](https://docs.travis-ci.com/api/#settings:-general)).

#### Parameters

  key                           | type      | description
  ----------------------------- | --------- | ----------------------------------
  `builds_only_with_travis_yml` | `boolean` | (optional) see API docs
  `build_pushes`                | `boolean` | (optional) see API docs
  `build_pull_requests`         | `boolean` | (optional) see API docs
  `maximum_number_of_build`     | `number`  | (optional) default: 0 (unlimited)

[travis-token]: https://blog.travis-ci.com/2013-01-28-token-token-token/
