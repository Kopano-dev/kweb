# reauth

Another authentication plugin for CaddyServer (v1). This is a forked version,
embedded into kweb with only a subset of the original backends.

See https://github.com/freman/caddy-reauth for the original upstream project.

## Contents

- [reauth](#reauth)
	- [Contents](#contents)
	- [Abstract](#abstract)
	- [Supported backends](#supported-backends)
	- [Supported failure handlers](#supported-failure-handlers)
	- [Configuration](#configuration)
		- [Spaces in configuration](#spaces-in-configuration)
	- [Backends](#backends)
		- [Simple](#simple)
		- [Upstream](#upstream)
		- [LDAP](#ldap)
	- [Failure handlers](#failure-handlers)
		- [HTTPBasic](#httpbasic)
		- [Redirect](#redirect)
		- [Status](#status)

## Abstract

Provides a common basis for various and multiple authentication systems. This came to be as we wanted to dynamically authenticate our
docker registry against gitlab-ci and avoid storing credentials in gitlab while still permitting users to log in with their own credentials.

## Supported backends

The following backends are supported.

* [Simple](#simple)
* [Upstream](#upstream)
* [LDAP](#ldap)


## Supported failure handlers

The following failure handlers are supported.

* [HTTPBasic](#httpbasic)
* [Redirect](#redirect)
* [Status](#status)

## Configuration

The core of the plugin supports the following arguments:

| Parameter-Name    | Description                                                                                        |
| ------------------|----------------------------------------------------------------------------------------------------|
| path              | the path to protect, may be repeated but be aware of strange interactions with `except` (required) |
| except            | sub path to permit unrestricted access to (optional, can be repeated)                              |
| failure           | what to do on failure (see failure handlers, default is [HTTPBasic](#httpbasic))                   |

Example:
```
	reauth {
		path /
		except /public
		except /not_so_secret
	}
```

Along with these two arguments you are required to specify at least one backend.

### Spaces in configuration

Through experimentation by [@mh720 (Mike Holloway)](https://github.com/mh720) it has been discovered that if you need spaces in your configuration that the best
bet is to use unicode escaping.

For example:
```
OU=GROUP\u0020NAME
```

I imagine this would allow you to escape any character you need this way including quotes.

## Backends

### Simple

This is the simplest plugin, taking just a list of username=password[,username=password].

Example:
```
	simple user1=password1,user2=password2
```

### Upstream

Authentication against an upstream http server by performing a http basic authenticated request and checking the response for a http 200 OK status code. Anything other than a 200 OK status code will result in a failure to authenticate.

Parameters for this backend:

| Parameter-Name    | Description                                                                              |
| ------------------|------------------------------------------------------------------------------------------|
| url               | http/https url to call                                                                   |
| skipverify        | true to ignore TLS errors (optional, false by default)                                   |
| timeout           | request timeout (optional 1m by default, go duration syntax is supported)                |
| follow            | follow redirects (disabled by default as redirecting to a login page might cause a 200)  |
| cookies           | true to pass cookies to the upstream server                                              |
| match             | used with follow, match string against the redirect url, if found then not logged in     |

Examples
```
	upstream url=https://google.com,skipverify=true,timeout=5s
  upstream url=https://google.com,skipverify=true,timeout=5s,follow=true,match=login
```

### LDAP

Authenticate against a specified LDAP server - for example a Microsoft AD server.

Parameters for this backend:

| Parameter-Name   | Description                                                                                                              |
| ------------------|-------------------------------------------------------------------------------------------------------------------------|
| url              | url, required - i.e. ldap://ldap.example.com:389                                                                         |
| tls              | should StartTLS be used? (default false)                                                                                 |
| username         | (read-only) bind username - i.e. ldap-auth                                                                               |
| password         | the password for the bind username                                                                                       |
| insecure         | true to ignore TLS errors (optional, false by default)                                                                   |
| timeout          | request timeout (optional 1m by default, go duration syntax is supported)                                                |
| base             | Search base, for example "OU=Users,OU=Company,DC=example,DC=com"                                                         |
| filter           | Filter the users, eg "(&(memberOf=CN=group,OU=Users,OU=Company,DC=example,DC=com)(objectClass=user)(sAMAccountName=%s))" |
| principal_suffix | suffix to append to usernames (eg: @example.com)                                                                         |
| pool_size        | size of the connection pool, default is 10                                                                               |

Example
```
	ldap url=ldap://ldap.example.com:389,timeout=5s,base="OU=Users,OU=Company,DC=example,DC=com",filter="(&(memberOf=CN=group,OU=Users,OU=Company,DC=example,DC=com)(objectClass=user)(sAMAccountName=%s))"
```

## Failure handlers

### HTTPBasic

This is the default failure handler and is by default configured to send the requested host as the realm

Parameters for this handler:

| Parameter-Name    | Description                                                                              |
| ------------------|------------------------------------------------------------------------------------------|
| realm             | name of the realm to authenticate against - defaults to host                             |

Example
```
	failure  basicauth realm=example.org
```

### Redirect

Redirect the user, perhaps to a login page?

Parameters for this handler:

| Parameter-Name    | Description                                                                              |
| ------------------|------------------------------------------------------------------------------------------|
| target            | target url for the redirection, supports {uri} for redirection (required)                |
| code              | the http status code to use, defaults to 302                                             |

Example
```
	failure redirect target=example.org,code=303
```

Example with uri
```
	failure redirect target=/auth?redir={uri},code=303
```

### Status

Simplest possible failure handler, return http status $code

Parameters for this handler:

| Parameter-Name    | Description                                                                              |
| ------------------|------------------------------------------------------------------------------------------|
| code              | the http status code to use, defaults to 401                                             |

Example
```
	failure status code=418
```
