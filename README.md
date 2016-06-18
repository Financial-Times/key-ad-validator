# Key-AD-validator service

Will vet all public keys via provided HTTP endpoint and expose on `/authorized_keys` via HTTP.

## Architecture

Current:

```
+------------+  GET    +---------------------------------------+
|coco cluster| ----->  |github.com ssh-keys/authorized_keys    |
+------------+         +---------------------------------------+
```

Future:

```
+------------------+   GET   +------------------------------------+    GET    +----------------------------------------+
|coco cluster VMs  | ------> | key-ad-validator                   | --------> | github.com ssh-keys/authorizid_keys    |
+------------------+         +------------------------------------+           +----------------------------------------+
                                                |
                                                | ldap-search each line       +----------------------------------------+
                                                \---------------------------> | AD server                              |
                                                                              +----------------------------------------+
```

Hosted `authorized_keys` file must conain user's email in the comments field:

```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJBZGih0OxlRiHkxzLCog6sJIQgqDYpfuPRBKEXAMPLE first.last@domain.com
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJBZGih0OxlRiHkxzLCog6sJIQgqDYpfuPRBKEXAMPLE first.last@domain.com
...
```

that will be used to lookup the user in LDAP.

Service verifies account is valid by checking `userAccountControl` - looking for value of `512`: http://www.selfadsi.org/ads-attributes/user-userAccountControl.htm

Verified keys are then exposed on `/authorized_keys` HTTP endpoint.
