Welcome to the pushtart wiki!

## Extensions

[DNSServ](https://github.com/twitchyliquid64/pushtart/wiki/DNSServ) - Full command reference

[HTTPProxy](https://github.com/twitchyliquid64/pushtart/wiki/HTTPProxy) - Full command reference

## Commands quick reference

#### Setting up a git repository for integration

```shell
# Server port is 2022 by default.
git remote add myPushtartRemote ssh://<server-address>:<server-port>/test
git push myPushtartRemote master
# The pushURL (you may need it in future commands) is /test in this example.
```

#### Accessing the management interface

_NOTE: Assumes the server is already running, and a user is already setup with either password or publickey authentication._

```shell
# Server port is 2022 by default.
ssh <hostname> -p 2022
```

#### Make a user

```shell
#Before the server is running
./pushtart make-user --username bob --password hi --allow-ssh-password yes
#Creates a user where: username=bob, password=hi, and user is allowed to SSH/git-push using his password

#While the server is running (run command in management interface)
 make-user --username bob --password hi --allow-ssh-password yes
#Creates a user where: username=bob, password=hi, and user is allowed to SSH/git-push using his password
#run ls-users to see your new user!
```

#### Setting up a users SSH key

If the server is running:

```shell
cat ~/.ssh/id_rsa.pub | ssh <server-address> -p <server-port> import-ssh-key --username <pushtart-username>
```

This assumes the user you are logged in as has an account on pushtart. If the username is different, use `<username>@<server-address>` in place of `<server-address>`. If you still cannot authenticate you will need to shut down the server and use the below command.

If the server is not running:

```shell
./pushtart import-ssh-key --username <pushtart-username> --pub-key-file ~/.ssh/id_rsa.pub
```


#### Configuring a tart (name, owners, environment variables etc)

Refer to the section: [Tart Configuration](https://github.com/twitchyliquid64/pushtart/wiki/Tart Configuration)


#### Viewing server logs without local access

From within the management interface: `logs`

From the command line (via SSH): `ssh <server-address> -p <server-port> logs`

_server-port is 2022 by default._