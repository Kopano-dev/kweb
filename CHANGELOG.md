# CHANGELOG

## Unreleased

- Generate changelog with https://github.com/git-chglog/git-chglog
- Correctly use CGO_ENABLED variable bin cmd build target


## v0.6.1 (2019-03-18)

- Remove default X-Frame-Options header
- Remove kweb injected response headers for konnect


## v0.6.0 (2019-03-01)

- Pass default_redirect and legacy_reverse_proxy from cfg
- Add license modification notice to fastcgi
- Improve fastcgi index support
- Lookup fastcgi index file in rule root instead of vhost root
- Enable local fastcgi directive as fastcgi2
- Add without sub directive to fastcgi
- Import fastcgi from caddy upstream (72d0debde6bf01b5fdce0a4f3dc2b35cba28241a)
- Implement alias directive
- Show details and print OK for make check


## v0.5.0 (2019-01-24)

- Update 3rd party dependencies
- Make Jenkins fetch vendor dependencies in its own step
- Use verbose dep ensure


## v0.4.4 (2019-01-24)

- Show Go version in Jenkins
- Bump base copyright years to 2019


## v0.4.3 (2018-12-21)

- Fixup binscript to properly use http2 and extra parameters
- Add dummy Docker healtcheck for DOcker swarm support
- Use correct assets path in Dockerfile
- Add support for Dockerfile to set assets path via env variable
- Add Dockerfile and scripts to build a docker image


## v0.4.2 (2018-11-20)

- Fixup more typos in internal config


## v0.4.1 (2018-11-20)

- Fix internal config error


## v0.4.0 (2018-11-19)

- Use new location of golint
- Use defined version of Go dep tool
- Add routes to kapi kvs v1 API


## v0.3.2 (2018-11-08)

- Fixup sandboxing with systemd


## v0.3.1 (2018-11-08)

- Remove debug of internal config
- Use latest systemd sandboxing


## v0.3.0 (2018-11-08)

- Remove unused v key from default configjson
- Allow override of config.json files
- Add rewrite rule to proxy kwm v2 API


## v0.2.0 (2018-10-11)

- Add support to load other files in root of static pwa
- Improve CSP handling


## v0.1.3 (2018-10-08)

- Remove reload support since it is not tested/not supported for now
- Remove systemd proc limit
- Fix wrong parameter


## v0.1.2 (2018-10-08)

- Fixup missing variable
- Fixup systemd nproc limit spelling


## v0.1.1 (2018-10-08)

- Add robots.txt and favicon.ico
- Add path pattern for relative static pwa


## v0.1.0 (2018-10-08)

- Add config, bin script and systemd service


## v0.0.9 (2018-10-05)

- Improve CSP security and allow data: URLs for img-src
- Use forked caddy-prometheus


## v0.0.8 (2018-10-04)

- Use kweb specific assets path


## v0.0.7 (2018-10-04)

- Always set most of the security headers
- Add HSTS support with default


## v0.0.6 (2018-10-04)

- Add Jenkinsfile
- Update deps


## v0.0.5 (2018-10-04)

- Add support for wss:// in pwa CSP
- Add request-log configuration
- Make it a multi host config by default


## v0.0.4 (2018-10-04)

- Fixup config JSON for missing double quotes
- Add support for extra configuration file


## v0.0.3 (2018-10-04)

- Fixup pwa folderish redirect with query
- Add bind and default-redirect parameters


## v0.0.2 (2018-10-04)

- Log requests in combined format
- Make staticpwa paths relative to -root value
- Route client side urls to pwa index
- Add grapi API to config.json
- Implement legacy reverse proxy support


## v0.0.1 (2018-10-02)

- Ignore www folder
- Add commandline options for serve command
- Add make vet target
- Add serve command with standard Kopano services
- Add LICENSE.
- Set more build variables
- Add minimal README
- Implement basic launch code
- Add license ranger with dep support
- Initial web server implementation
- Initial commit
