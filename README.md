# Kopano Web

Flexible web server with integrated URL routing for Kopano services.

## Build dependencies

Make sure you have Go 1.12 or later installed. This assumes your GOPATH is `~/go` and
you have `~/go/bin` in your $PATH and you have [Dep](https://golang.github.io/dep/)
installed as well.

## Building from source

```
mkdir -p ~/go/src/stash.kopano.io/kgol
cd ~/go/src/stash.kopano.io/kgol
git clone <THIS-PROJECT> kweb
cd kweb
make
```
