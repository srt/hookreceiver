Hookreceiver
============

Listens for SCM hooks from Bitbucket and executes a shell command when a notification is received.

Currently supported hooks:

* Bitbucket [POST Hook](https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management):
  Use `http://your-host:8080/hooks/bitbucket/whatever` for the URL.

Usage
-----

<pre>
$ hookreceiver -h
Usage of hookreceiver:
  -c="/etc/hookreceiver.conf.d": Config path (file or directory)
</pre>


Configuration File
------------------

Configuration files are JSON files. You define the adress/port hookreceiver will listen on (Addr) and 
an array of repository configurations. When a notifaction is received for one of the configured
repositories the given command is executed with `/bin/sh -c`. You can restrict the command to 
notifications that contain changes for a specific branch using the `Branch` property and set the
working directory for the command with `Dir`.

<pre>
{
  "Addr": ":8080",
  "Repositories": [
    {
      "Url": "https://bitbucket.org/srt/foo",
      "Branch": "master",
      "Command": "git pull",
      "Dir": "/var/www/foo"
    },
    {
      "Url": "https://bitbucket.org/srt/bar",
      "Branch": "develop",
      "Command": "git pull",
      "Dir": "/var/www/bar"
    }
  ]
}
</pre>

Configuration Directory
-----------------------

As an alternative you can use a configuration directory. All files in that directory will be merged to 
form the configuration. This makes it easy to use hookreceiver with puppet and similar systems.

00-main.conf:
<pre>
{
  "Addr": ":8080"
}
</pre>

01-foo.conf:
<pre>
{
  "Repositories": [
    {
      "Url": "https://bitbucket.org/srt/foo",
      "Branch": "master",
      "Command": "git pull",
      "Dir": "/var/www/foo"
    }
  ]
}
</pre>

02-bar.conf:
<pre>
{
  "Repositories": [
    {
      "Url": "https://bitbucket.org/srt/bar",
      "Branch": "develop",
      "Command": "git pull",
      "Dir": "/var/www/bar"
    }
  ]
}
</pre>

License
-------

Hookreceiver is released under the [MIT License](http://www.opensource.org/licenses/MIT).