# Image Server

[![Build Status](https://magnum.travis-ci.com/wanelo/image-server.svg?token=xxYxjHDAXkDK41qZ1dqA&branch=master)](https://magnum.travis-ci.com/wanelo/image-server)


## Image Generation

### Posting New Images

An image needs to be uploaded to a namespace.

    Namespaces allow to group image types. This allows different groups of images to have different dimensions and proccessings. For example avatars will require different image sizes than product images.

Uploading an image requires a source
```
POST http://localhost:7000/p?source=http://example.com/image.jpg
```

It is possible to pre-generate images and save them to the configured file store by passing the outputs when posting the image.

Example with curl
```shell
curl -X POST http://localhost:7000/p?outputs=x300.jpg,x300.webp\&source=http://example.com/image.jpg
```

The request returns the *"Image Information"* after the original image is saved to the file store.
Image outputs are generated after the request is complete. The response includes properties of the image, and the image hash to be used to retrieve it in the future.

### Image Information

Image properties can be retrieved by visiting the info page. The response is the same as the one returned when creating the image.
```
GET http://localhost:7000/p/f84/0ee/339d264d4bab1b169a653b1a91/info.json
```

```json
{
	"hash": "f840ee339d264d4bab1b169a653b1a91",
	"partitionedHash": "f84/0ee/339d264d4bab1b169a653b1a91"
	"height": "520",
	"width": "400"
}
```

### Sample images

**Dimensions**

    Maximum Width
    GET http://localhost:7000/p/f84/0ee/339d264d4bab1b169a653b1a9/w500.jpg

    Square
    GET http://localhost:7000/p/f84/0ee/339d264d4bab1b169a653b1a9/x600.jpg

    Rectangle
    GET http://localhost:7000/p/f84/0ee/339d264d4bab1b169a653b1a9/300x400.jpg

**Quality**

The default compression of the image can modified by appending `-q` and the desired quality `1-100`.

    Square with quality 50
    GET http://localhost:7000/p/f84/0ee/339d264d4bab1b169a653b1a9/x600-q50.jpg

## Manta

Command for manta task:
```shell
image --remote_base_path public/images --outputs x300.webp,x300.jpg process $MANTA_INPUT_FILE
```

## Error Handling

Few errors will cause the server to return error pages

- Source image is not found: NotFound (404)
- Image requested is larger than maximum_width: NotAcceptable (406)

## Development

Set up the environment:

```bash
mkdir -p $GOPATH/src/github.com/wanelo/
git clone git@github.com:wanelo/image-server.git $GOPATH/src/github.com/wanelo/image-server
ln -s $GOPATH/src/github.com/wanelo/image-server ~/workspace/image-server
cd ~/workspace/image-server
```

Install dependencies:

```bash
brew bundle
make deps
```

Install go from source:

```
hg clone https://code.google.com/p/go $HOME/go.trunk
cd $HOME/go.trunk/src && ./all.bash
cd $HOME/go.trunk/src && ./make.bash
cd $HOME/go.trunk/src && GOOS=solaris GOARCH=amd64 ./make.bash
```

Set up editor:

  - Atom.io package `go plus`

Compile the app:

To build the executables under `./bin`

```bash
make build
```

## Development Usage

There are few `make` helpers that start the development server. They all translate environment variables into flags.

### S3
Required ENV variables: `IMG_OUTPUTS`, `AWS_BUCKET`, `IMG_REMOTE_BASE_PATH`, `IMG_REMOTE_BASE_URL`
```
make dev-server-s3
```

### Manta
Required ENV variables: `IMG_OUTPUTS`, `MANTA_URL`, `MANTA_USER`, `MANTA_KEY_ID`, `SDC_IDENTITY`, `IMG_MANTA_BASE_PATH`
```
make dev-server-manta
```

### No uploader, only store images locally
Required ENV variables: `IMG_OUTPUTS`

```bash
make dev-server
```

## Tests

S3 tests make real HTTP calls when you have S3 ENV variables set.

```
make tests
```

## Configuration

All configuration is passed by flags


## Deploy

Make sure you increase the version number in core/version.go

```bash
make release
```

## Graphite Events

A local cache was not found and the image was processed. This also tracks count of images sent to remote store.
```
stats.image_server.image_request
```

In addition, the format is tracked (jpg, gif, webp)
```
stats.image_server.image_request.jpg
```

Request failed to return an image
```
stats.image_server.image_request_fail
```

Every download from original source, and a 404 was returned
```
stats.image_server.original_downloaded
```

The original image is not available, and a 404 was returned
```
stats.image_server.original_unavailable
```

## Profiling

On the production server

```
curl http://localhost:6060/debug/pprof/heap > images.pprof
```

on development machine
```
scp user@example.com:images.pprof .
go tool pprof bin/solaris/images --inuse_objects images.pprof
```

## Benchmarks

Make sure your computer can handle enough simultaneous connections. MacOS X by default allows 128. Need a lot more!

```shell
$ sudo sysctl -w kern.ipc.somaxconn=2048
```

Also need to increase the limit of maximum open files

To find out the limits on your computer:
```shell
launchctl limit
```

Increase the limits!
```shell
launchctl limit maxfiles 400000 1000000
```

to increase them [permanently](https://coderwall.com/p/lfjoaq)

### After Release

- Zero-downtime restart: http://rcrowley.org/talks/strange-loop-2013.html#27
- Configuration reload with signal
- Status page
  - current images processing count
  - current original download count
