# Hookreceiver

[![Build Status](https://travis-ci.org/cweagans/hookreceiver.svg?branch=master)](https://travis-ci.org/cweagans/hookreceiver)

Listens for webhooks from code hosting services and executes one or more shell commands when a notification is received.

Currently supported hooks:

* Bitbucket [POST Hook](https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management):
  Use `http://your-host:8080/hooks/bitbucket/repo_name` for the URL.
* Gitlab [Push Hook](http://doc.gitlab.com/ce/web_hooks/web_hooks.html):
  Use `http://your-host:8080/hooks/gitlab/repo_name` for the URL.
* Stash [POST service webhook](https://confluence.atlassian.com/display/STASH/POST+service+webhook+for+Stash):
  Use `http://your-host:8080/hooks/stash/repo_name` for the URL.
* Generic Hook that does not parse the payload:
  Use `http://your-host:8080/hooks/generic/repo_name` for the URL.

## Installation

### Manual installation

Download the binary from the latest release on Github and put it in your `$PATH` somewhere.

### Docker

`docker run -v ./conf:/etc/hookreceiver.conf.d cweagans/hookreceiver:latest`

Make sure that `./conf` contains one or more configuration files.

### Build from source

* Install the Go toolchain
* Clone this repo
* Run `go build` (no external dependencies are needed)

## Usage

```bash
$ hookreceiver -h
Usage of hookreceiver:
  -c="/etc/hookreceiver.conf.d": Config path (file or directory)
```

## Configuration format

Configuration files are JSON files. You define the adress/port hookreceiver will listen on (Addr) and an array of repository configurations. When a notifaction is received for one of the configured repositories identified by `URL` the given commands are executed with `/bin/sh -c` in the working directory set with `Dir`.

`URL` indicates the canonical URL of the repository. For example when using Bitbucket this is `https://bitbucket.org/user/repo`.

Instead of using the `URL` property you can use `Name`. This is the `repo_name` part in the hookreciever URL. If you configure your repository to call `http://your-host:8080/hooks/bitbucket/bar` the `Name` would be `bar`. Using the name allows to configure repository providers that do not include the URL in their payload like Stash or the Generic provider.  

```json
{
  "Addr": ":8080",
  "Repositories": [
    {
      "URL": "https://bitbucket.org/srt/foo",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/foo"
    },
    {
      "Name": "bar",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/bar"
    }
  ]
}
```

## Configuration Directory

As an alternative you can use a configuration directory. All files in that directory will be merged to 
form the configuration. This makes it easy to use hookreceiver with puppet and similar systems.

`00-main.conf`:
```json
{
  "Addr": ":8080"
}
```

`01-foo.conf`:
```json
{
  "Repositories": [
    {
      "URL": "https://bitbucket.org/srt/foo",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/foo"
    }
  ]
}
```

`02-bar.conf`:
```json
{
  "Repositories": [
    {
      "Name": "bar",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/bar"
    }
  ]
}
```


## Provider-specific notes

### Bitbucket

Use `http://your-host:8080/hooks/bitbucket/repo_name` for the URL.

Bitbucket supports the `Name` and `URL` configuration properties. The URL is of the form `https://bitbucket.org/user/repo`.

You can also restrict the command to notifications that contain changes for a specific branch using the `Branch` property. You may miss some commits as Bitbucket only provides detailed information like affected files and branches if pushes do not exceed a certain size limit. Thus using the `Branch` property is generally discouraged.

Example:
```json
{
  "Repositories": [
    {
      "URL": "https://bitbucket.org/srt/foo",
      "Branch": "develop",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/foo"
    }
  ]
}
```

See also: https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management

### Gitlab

Use `http://your-host:8080/hooks/gitlab/repo_name` for the URL.

Gitlab supports the `Name` and `URL` configuration properties. The URL is of the form `https://gitlab.example.com/group/repo`. This is the value `Homepage` in the JSON payload.

See also: http://doc.gitlab.com/ce/web_hooks/web_hooks.html

### Stash

Use `http://your-host:8080/hooks/stash/repo_name` for the URL.

Stash does not include the canonical URL in the payload so hookreciever is unable to match a stash repository via URL. Use `Name` instead of `URL` to configure Stash hooks.

See also: https://confluence.atlassian.com/display/STASH/POST+service+webhook+for+Stash

### Generic

Use `http://your-host:8080/hooks/generic/repo_name` for the URL.

The Generic provider does not parse the payload so it can be used to support arbitray systems. Use `Name` instead of `URL` to configure Generic hooks.

## License

Hookreceiver is released under the [MIT License](http://www.opensource.org/licenses/MIT).

## Acknowledgements

* Thanks to @srt for writing the original version of hookreceiver.
