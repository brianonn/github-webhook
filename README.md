# Github webhook listener to run a script

First commit.

This small golang app listens on port 5000 for calls to /webhook and
checks for pushes to a branch. When it sees a push to a specific
branch, it will start a shell script. The branch is currently
hardcoded to `staging`, the script is hardcoded to `push.sh`.  It is
intended that the script will contain the commands needed to get the
current release and push that to your webserver(s) and restart the
servers.  For example, the script can pull the current branch and run
rsync to send it to remote servers, or launch something on AWS, or
build a docker container and push that to a repo, and/or ssh to the
remote server to start up a web app, etc.

The _github-webapp_ application will read an environment variable
`GITHUB_WEBHOOK_SECRET` for a secret string, used to prevent random
internet user from causing the webhook to run. The same string should
also be set in the webhook secret field on the webhook page in Github.

The app uses a single semaphore to limit the number of simultaneously
running external scripts to one. So multiple pushes to the same branch
within a short time will cause all but the first one to wait for
access to the semaphore. Once the external script completes and the
semaphore is released, it will allow the next one to run. It is not
deterministic which one of multiple waiting pushes will succeed, and
there is no guaranteed ordering.


## TODO

- pass port and branch arguments to the script from the command line
- config file for mapping events on repos to scripts and args
- option to whitelist github servers to be the only ones that can
  originate a webhook.
- option to restrict specific git users only (maybe part of the config file too)
- use the form `<repo>.<event>.<branch>.sh` for finding the script to
  run. This will allow the same instance of `github-webhook` to serve
  different repos and different events
- use an in-order waiting queue, not a semaphore, so that pushes are
  held and executed in order, rather than randomly
- use cobra for command line arguments
- use viper for a proper config file
