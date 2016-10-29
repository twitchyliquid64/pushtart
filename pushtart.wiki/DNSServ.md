# DNSServ Extension
When enabled, DNSServ provides a simple DNS server. Records can be managed via commands, or be automatically added by a tart (see documentation about the tartconfig file).

DNSServ can also act as an upstream (caching) DNS server - That way you can use it as your nameserver!




## Enabling, Disabling, Configuring

_NOTE: Pushtart must be restarted for any commands mentioned in this section to take effect._

**Enable**: `./pushtart extension --extension DNSServ --operation enable`

**Disable**: `./pushtart extension --extension DNSServ --operation disable`

**Set Listener host/port**: `./pushtart extension --extension DNSServ --operation set-listener --listener ":53"`

**Enable Recursion **(let pushtart be your cache nameserver): `./pushtart extension --extension DNSServ --operation enable-recursion`

When operating in recursive mode, DNSServ has a fixed-size LRU cache with a 30-minute timeout. The size of this cache can be changed with the following command:

`./pushtart extension --extension DNSServ --cache-size 100`


## Listing domains and showing configuration

**Show configuration** (all plugins): `extension --operation show-config`

**List domains**: `ls-dns-domains`

## Adding/Deleting records

**Creating a DNS record** type A, resolving the domain crap.com to 192.168.1.1 with a cache timeout of 100 seconds:

`./pushtart extension --extension DNSServ --operation set-record --type A --domain crap.com --address 192.168.1.1 --ttl 100`


**Delete the type A record** for crap.com.

`./pushtart extension --extension DNSServ --operation delete-record --domain crap.com`

_NOTE: Only type A records can be set ATM. More record type (AAAA, MX, TXT) are planned, though let me know if you need it and I can expedite their development._