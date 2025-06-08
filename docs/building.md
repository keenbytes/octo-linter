# Building

## Binary
Run `go build -o octo-linter` to compile the binary.

Use `GOOS` and `GOARCH` environment variables to build binary for a specific platform.  More information
can be found in the [Go docs](https://go.dev/doc/install/source#environment).

### Docker image
To build the docker image, use the following command.

````
docker build -t octo-linter .
````

If an image is built on a different platform than `linux/amd64`, an additional `--platform=linux/amd64`
argument is necessary.  See [command reference](https://docs.docker.com/reference/cli/docker/buildx/build/#platform)
for `docker build`.
