
Pushtart - The worlds easiest PaaS.
=======================================

Pushtart runs persistantly on any *nix box you own. You `git push` projects to it, and pushtart saves them on the box and runs the repositories' `startup.sh`. Its that simple!

There is also a simple (but fully featured) user management system, as well as the ability to set environment variables on each of your deployments (keep sensitive information out of your repositories!).

## Getting started

To get started, we need to build the program, run `pushtart make-config` to generate a configuration file, and setup a user with `pushtart make-user`. Then we are ready to run!

Assuming you have Go >1.6 installed on your system:

1. `git clone https://github.com/twitchyliquid64/pushtart`
2. `cd pushtart`
3. `export GOPATH=$PWD`
4. `go build`
5. `./pushtart make-config`
6. `./pushtart make-user --username bob --password hi --allow-ssh-password yes`
7. `./pushtart run`

You are now ready to `git push` your projects!

```
git remote add pushtart_server ssh://localhost:2022/test
git push pushtart_server master
```
:DD

#### Management interface

You can actually SSH into pushtart, and run the same commands you can on the unix command line (except make-config, and importing a SSH key).

`ssh <hostname> -p 2022`

Of course, you can change the port and bind-host in the configuration file.

#### Setting up a SSH key

Most sane people prefer to use SSH keys instead of passwords. To setup a key with an existing user, simply run this command before you start the server:

`./pushtart import-ssh-key --username <insert-pushtart-username-here> --pub-key-file ~/.ssh/id_rsa.pub `

#### Terminology

Everything makes sense, except I decided to call a repository in pushtart a 'tart' :). Tarts can be in running or stopped states - this is all controllable through commands.

### USAGE

```
USAGE: pushtart <command> [--config <config file>] [command-specific-arguments...]
if no config file is specified, config.json will be used.
SSH server keys, user information, tart status, and other (normally external) information is stored in the config file.
Commands:
	run (Not available from SSH shell)
	make-config (Not available from SSH shell)
	import-ssh-key --username <username> [--pub-key-file <path-to-.pub-file>] (Not available from SSH shell)
	make-user --username <username [--password <password] [--name <name] [--allow-ssh-password yes/no]
	edit-user --username <username [--password <password] [--name <name] [--allow-ssh-password yes/no]
	ls-users
	ls-tarts
	start-tart --tart <pushURL>
	stop-tart --tart <pushURL>
	edit-tart --tart <pushURL>[--name <name>] [--set-env "<name>=<value>"] [--delete-env <name>]
```

### TODO

 - [ ] Lock configuration file (.lock file? when pushtart is running)
 - [ ] Implement way to load a users ssh key when the server is running
 - [ ] Implement access controls to prevent different users from touching tarts they didnt create
 - [ ] Logging to file / console?