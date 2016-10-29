## Configuration overview

You should configure your tart in two spots.

Firstly, your tart probably needs to get general information about its running environment (for instance, you may need to give it DB address/passwords or the like). This is accomplished by setting environment variables using the commands below.

Secondly, you may need to configure pushtart-specific parameters, like access controls for the tart, its name, whether it logs or not. These are also configured using commands below however are not accessible to the tart when it runs.

#### Storing configuration in your respositories

To avoid the overhead of running a heap of commands manually everytime you setup a tart, you can put the commands into a file in your repository, and these commands will be run on every `git push`.

Simply create a file in your project root `tartconfig`. Commands which require the `--tart` argument can omit `--tart`, we will populate that argument for you. Lastly, you can use your tart's environment variables in your tartconfig files in the same manner as bash: `$varname, or ${varname}`.

If you are using environment variables and `tartconfig`, consider creating your tart prior to a `git push`, and setting up the environment variables then. That will mean your `tartconfig` is run with the correct environment variable values at first run. See below for how to 'precreate' your tart.


## Command reference

#### List all tarts on the system

This command also prints their configuration.

```shell
ls-tarts
```

#### Set the  name of your tart

```shell
edit-tart --tart <pushURL> --name "Tart Bob"
```

#### Start / Stop a tart

```shell
start-tart --tart <pushURL>
stop-tart --tart <pushURL>
```

#### Give tart management permissions to other users

By default, only the user who created a tart can edit/start/stop it. Give access to other users to allow them to run these commands.

```shell
tart-add-owner --tart <pushURL> --username <username>
tart-remove-owner --tart <pushURL> --username <username>
```

#### Make the tart log its stdout/stderr

```shell
edit-tart --tart <pushURL> --log-stdout yes
edit-tart --tart <pushURL> --log-stdout no
```


#### Set/delete a tart's environment variables.

When a tart runs, it runs with these environment variables set.

```shell
edit-tart --tart <pushURL> --set-env "variable_name=variable_value"
edit-tart --tart <pushURL> --delete-env "variable_name"
```

#### Reparse the tart's `tartconfig` file.

If environment variables have changed and your tarts `tartconfig` file makes use of them, you may wish to re-execute all of the commands.

```shell
digest-tartconfig --tart <pushURL>
```

#### Setup a tart to automatically restart when it stops

Always specify a reasonable lull-period - this will prevent resource wastage if your tart is stuck in a crashloop.

```shell
tart-restart-mode --tart <pushURL> --enabled yes/no --lull-period <seconds>
```

#### Precreate a tart

You should only use this feature to set environment variables prior to your first `git push`. Make sure the pushURLs will match.

```shell
new-tart --tart <pushURL>
```