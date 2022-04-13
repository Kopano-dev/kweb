# Kopano Web

Flexible web server with integrated URL routing for Kopano services.

## Build dependencies

Make sure you have Go 1.18 or later installed. This project uses Go modules.

## Building from source

```
git clone <THIS-PROJECT> kweb
cd kweb
make
```

### Build with Docker

```
docker build -t kwebd-builder -f Dockerfile.build .
docker run -it --rm -u $(id -u):$(id -g) -v $(pwd):/build kwebd-builder
```

## License

See `LICENSE.txt` for licensing information of this project.
