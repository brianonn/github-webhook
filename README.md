# Github webhook listener to run a script

First commit.

This small golang app listens on port 5000 for calls to /webhook and
checks for pushes to a branch. The branch is currently hardcoded to
`staging`, the script is hardcoded to `push.sh`

The app will read an environment variable `GITHUB_WEBHOOK_SECRET` for
a secret string, used to prevent random internet uses from causing the
webhook to run. The same string should also be set in the webhook secret field
on the webhook page in Github.

The app uses a single semaphore to limit the number of requests to
one. So multiple pushes to the same branch within a short time will
wait for each to finish. It is not deterministic which one of multiple
pushes will succeed.


## TODO

- pass port and branch arguments to the script from the command line
- config file for mapping events on repos to scripts and args
- option to whitelist github servers to be the only ones that can
  originate a webhook.
- option to restrict specific git users only (maybe part of the config file too)
- cobra (command args)
- viper (config)
