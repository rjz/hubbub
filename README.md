# Hubbub

VCS-configurations-as-code: establish, apply, and maintain transparent policies
for your many source-controlled projects.

Teams and individuals might use it to:

  - Manage repository settings
  - Check in licenses or contributing guidelines
  - Configure webhooks and third-party services
  - Configure CI integrations

Repositories hosted on [Github][github] are supported out of the box;
contributions for integrating with other third-party integrations and VCS hosts
[are welcome][contributing]!

## Build

Build hubbub using go >= 1.4 and make:

    $ make

## Get started

  1. Define a policy
  2. Assign it to your repositories

### Define a policy

Policies are JSON documents that describes the desired state of the repository
in terms of sequential goals.

Let's create a policy that adds a 'hello_world.txt' file to subject
repositories:

    [
      {
        "github_file" : {
          "state": "present",
          "ref": "heads/master",
          "name": "hello_world.txt",
          "content":"'Hi!' --hubbub"
        }
      },
    ]

Save it as `hello_world.json`.

### Assign it to your repositories

Next, let's create a list of repos that will be subject to the policy.

    [
      { "url": "github.com/rjz/uno" },
      { "url": "github.com/rjz/dos" },
      { "url": "github.com/rjz/tres" }
    ]

Save it as `repos.json`.

### Apply the policy

In order to use the Github API, we'll need to [obtain][github-token] a valid
access token and add it to the environment:

    $ export HUBBUB_GITHUB_ACCESS_TOKEN=xyz

Finally, we can use `hubbub` to add `hello_world.txt` to each of the subject
repos.

    $ hubbub \
      -policy=hello_world.json \
      -repositories=repos.json

### Service Integrations

Check out each [service's README](services/).

## License

MIT

[github]: https://github.com
[github-token]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
[contributing]: CONTRIBUTING.md

