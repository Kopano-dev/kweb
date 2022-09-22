# CHANGELOG

## Unreleased



## v0.15.2 (2022-09-22)

- Fix Go 1.19 compatibility by updating h3 support
- Build with Go 1.19.1


## v0.15.1 (2022-04-27)

- Fix HTTP/3 request support with proxy backend


## v0.15.0 (2022-04-13)

- Do not run go mod tidy on build
- No longer run gofmt on all builds, it now needs to be run explicitly
- Bump required minimal Go version to 1.18
- Install Go build dependency tools with go install for Go 1.18 compatibility
- Build with and bump to Go 1.18


## v0.14.0 (2022-02-10)

- Update HTTP/3 support library
- Build with Go 1.17.6
- Raise checkstyle quality gate to 200
- Fix all reauth unit tests
- Properly set remote user when authenticating with reauth
- Hook up reauth plugin as embeeded plugin
- Import reauth plugin


## v0.13.1 (2021-12-09)

- Build with Go 1.17.5


## v0.13.0 (2021-10-08)

- Build with Go 1.17.2
- Fix proxy healthcheck support for unix upstreams
- Build with Go 1.16
- Bump golangci-linter to 1.41.1
- Update quic to support final HTTP3 specification
- Fix request_log folder creation check to avoid changing owner
- Mark mcu and admin API of kwmserver internal


## v0.12.4 (2020-09-29)

- Explicitly define backend read and write timeouts
- Limit automatic gzip file extensions


## v0.12.3 (2020-09-22)

- Fix alias directive redirect targets


## v0.12.2 (2020-09-21)

- Improve request log file creation logic


## v0.12.1 (2020-09-21)

- Ensure that default request log path is actually created
- Clarify ACME agree config file options
- Fix fastcgi2 to support SSL variables when TLSv1.3 is used


## v0.12.0 (2020-09-17)

- Add Referrer-Policy: no-referrer for WebApp
- Move static PWA CSP to configuration
- Make rate limiting flexible via configuration
- Remove builtin routes for obsolte kwmserver v1 API
- Change kweb config generator template syntax
- Refactor built in config to allow multiple hosts
- Hide KDAV version response header
- Add log rotation to requests log
- Use better example for request log file location
- Fix extra.d configuration support


## v0.11.1 (2020-09-15)

- Relax config file exist check
- Initialize process log on stdout
- Add marker around config dump output
- Fix passing of parameters for combined folders


## v0.11.0 (2020-09-15)

- Make configuration flexible and overridable


## v0.10.1 (2020-09-07)

- v0.10.1
- Build with Go 1.14.8


## v0.10.0 (2020-08-21)

- Use extra.d directory instead of extra.cfg
- Add support for extra configuration folder
- Add knob to dump the internal auto generated Caddyfile
- Change hostname default to * to accept all hostnames
- Fix compatibility with Go 1.15
- Move Jenkins warnings-ng to post phase
- Update 3rd party dependencies
- Add build args for Docker based builds
- Build with Go 1.14.7
- Allow chown to fail when running in Docker container
- Use non-deprecated warnings-ng plugin in Jenkins
- Add examples to test url routing and fastcgi timeouts


## v0.9.2 (2020-06-15)

- Update 3rd party dependencies
- Build with Go 1.14.4


## v0.9.1 (2020-03-09)

- Use vendor folder to generate 3rd party license file


## v0.9.0 (2020-03-09)

- Update 3rd party dependencies
- Build with Go 1.14
- Bump linter to v1.21.0
- Require unit tests to run in Jenkins pipeline
- Add test coverage reporting
- Add minimal version unit test
- Improve build process


## v0.8.3 (2019-11-07)

- Build with Go 1.13.4


## v0.8.2 (2019-10-28)

- Update third party dependencies to lates patch levels
- Build with Go 1.13.3


## v0.8.1 (2019-09-30)

- Build with Go 1.13.1
- Improve mod support in Makefile and linter
- Udpate Dockerfiles syntax
- Use Dockerfile.build in Jenkins
- Move to golangci-lint for Go linting
- Use vendor folder for build


## v0.8.0 (2019-09-10)

- Tell linker the buildid instead of patching binary
- Clear actionid part from buildid for truly reproducible builds
- Update license ranger to support Go modules
- Move from Dep to Go modules
- Build reproducible with Jenkins
- Use new trimpath feature to make builds reproducible
- Bump to Go 1.13
- Cleanup Dockerfile
- Add support for priviledged ports without root


## v0.7.0 (2019-07-30)

- Add additional commandline parameters for serve
- Update README
- Enable TLSv1.3 by default
- Update to Caddy 1.0 and dependencies
- Add import binder
- Install mailcap package
- Update Dockerfile
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

