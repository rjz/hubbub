# Github integration

Describe policies for repositories managed on github

## Configuration

To use the github integration, [obtain an access token][github-token] for your service and add it to your environment.

    $ export HUBBUB_GITHUB_ACCESS_TOKEN=<your token>

## Goals

### `github_file`

Manage a file within an existing [git ref](https://git-scm.com/book/en/v2/Git-Internals-Git-References) ([API documentation](https://git-scm.com/book/en/v2/Git-Internals-Git-References)).

#### Parameters

  key        | type     | description
  ---------- | -------- | ----------------------------------
  `state`    | `string` | one of `"absent"` OR `"present"`
  `ref`      | `string` | a valid ref (e.g. `"heads/master"` for the master branch)
  `name`     | `string` | the filename within the repo
  `content`  | `string` | (optional) the content
  `filename` | `string` | (optional) the local file to copy to the repo

**NOTE**: Specifying both `content` and `filename` is ambiguous and will
cause an error.

#### Example

    "github_file": {
      "state": "present",
      "ref": "heads/master",
      "name": ".gitignore",
      "content":"npm-debug.log\nhumans.txt"
    }

### `github_webhook`

Manage a github webhook ([API documentation](https://developer.github.com/webhooks/)).

#### Parameters

  key        | type            | description
  ---------- | --------------- | ----------------------------------
  `state`    | `string`        | one of `"absent"` OR `"present"`
  `active`   | `boolean`       | whether the hook should be enabled
  `events`   | `array[string]` | events to apply the hook to (see [full list][gh-hook-events])
  `name`     | `string`        | (optional) default `"web"`; override for [service hooks][gh-service-hooks]
  `config`   | `object`        | settings for the hook; format varies by service (see [docs][gh-webhook-config])

**NOTE**: hook "uniqueness" is currently determined by `config.url`, which may
not be available for all third-party services.

#### Example

Example configuration for a simple web (i.e., non-service) hook:

    "github_webhook": {
      "state": "present",
      "config": {
        "url": "https://my-service.com/hooks/github",
        "content_type": "json",
        "insecure_ssl": 0,
        "secret": "abc123"
      },
      "active": true,
      "events": [
        "pull_request"
      ]
    }

[github-token]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
[gh-service-hooks]: https://developer.github.com/webhooks/#service-hooks
[gh-hook-events]: https://developer.github.com/webhooks/#events
[gh-webhook-config]: https://developer.github.com/v3/repos/hooks/#parameters
