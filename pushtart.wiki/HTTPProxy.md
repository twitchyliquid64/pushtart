HTTPProxy is an extension that provides a HTTP proxy, configurable from within pushtart and from your tarts' `tartconfig`. In addition, it allows you to setup authentication rules for each tart, using pushtarts user system to allow/deny access.


## Enabling, Disabling, Configuring

_NOTE: Pushtart must be restarted for any commands mentioned in this section to take effect._

**Enable**: `./pushtart extension --extension HTTPProxy --operation enable`

**Disable**: `./pushtart extension --extension HTTPProxy --operation disable`

**Set Listener host/port**: `./pushtart extension --extension HTTPProxy --operation set-listener --listener ":8080"`

**Set default service domain**: `./pushtart extension --extension HTTPProxy --operation set-default-domain --domain "localhost"`


## Listing proxy entries and showing configuration

**Show configuration** (all plugins): `extension --operation show-config`

**List domains**: `ls-domain-proxies`

## Adding / deleting domain proxies

To create a proxy which will proxy requests with the host `testdomain` to `golang.org` using http, enter:

```shell
extension --extension HTTPProxy --operation set-domain-proxy --domain testdomain --targetport 80 --scheme http --targethost golang.org
```

To delete the above proxy:

```shell
extension --extension HTTPProxy --operation delete-domain-proxy --domain testdomain
```

## Guarding use of your proxies using auth rules

Allowing/denying access is setup by creating 'rules'. When you first create a domain, it has no rules, so all requests are allowed without authentication. If any rules are created for the domain, by default requests will be denied to that domain unless a `ALLOW` rule matches. `DENY` rules are evaluated before `ALLOW` rules, allowing you to create specific blocks even if you have broad `ALLOW` rules.

The following rule types are currently implemented:

|      Type      |                  Description                                    | Additional parameters (required) |
| -------------- | --------------------------------------------------------------- | -------------------------------- |
| USR_ALLOW      | Allows access to the user specified by --username.              | `username`                       |
| USR_DENY       | Explicitly denies access to the user specified by --username.   | `username`                       |
| ALLOW_ANY_USER | Allows access to any pushtart user.                             |                                  |

The following example adds a rule that will allow access to `myusername` for the proxy `testdomain`. From this point forward, no anonymous users will be granted access, and other pushtart users will also be denied (unless other rules are created for them).

The second line deletes an existing rule that denies access to `myusername`. If both matchin deny and allow rules exist, deny will take precedence.

```shell
extension --extension HTTPProxy --operation add-authorization-rule --type USR_ALLOW --domain testdomain --username myusername
extension --extension HTTPProxy --operation remove-authorization-rule --domain testdomain --type USR_DENY --username myusername
```

## Status pages

If you enable HTTPProxy, you can access a few additional management pages. There are:

http(s)://<default-domain>:<http-port>/health - Returns 'Ok' if the server is healthy

http(s)://<default-domain>:<http-port>/status - Displays a summary of non-sensitive configuration.

More will be added in future, including pages which require authentication.
