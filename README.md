# Image Server

[![Build Status](https://magnum.travis-ci.com/wanelo/image-server.svg?token=xxYxjHDAXkDK41qZ1dqA&branch=master)](https://magnum.travis-ci.com/wanelo/image-server)


## Image Generation

### Posting New Images

An image needs to be uploaded to a namespace.

Namespaces allow to group image types. For example avatars will require different image sizes than product images.

Uploading an image requires a source
```shell
curl -X POST http://localhost:7000/p?source=http://example.com/image.jpg
```

A binary file might be uploaded as well. The contents of the image need to be included in the body of the request.
```shell
> curl --data-binary "@./test/images/wine.jpg" -X POST http://localhost:7000/p
{
  "hash": "6e0072682e66287b662827da75b244a3",
  "height": 600,
  "width": 800,
  "content_type": "image/jpeg"
}
```

It is possible to pre-generate images by passing the outputs when posting the image.

```shell
> curl --data-binary "@./test/images/wine.jpg" -X POST http://localhost:7000/p?outputs=x300.jpg,x300.webp
{
  "hash": "6e0072682e66287b662827da75b244a3",
  "height": 496,
  "width": 574,
  "content_type": "image/jpeg"
}
```

The request returns the *"Image Information"* after the original image is saved to the file store.
Image outputs are generated after the request is complete. The response includes properties of the image, and the image hash to be used to retrieve it in the future.

### Image Information

Image properties can be retrieved by visiting the info page. The response is the same as the one returned when creating the image.
```shell
> curl http://localhost:7000/p/6e0/072/682/e66287b662827da75b244a3/info.json
{
  "hash": "6e0072682e66287b662827da75b244a3",
  "height": 496,
  "width": 574,
  "content_type": "image/jpeg"
}
```

### Image processing

**Dimensions**

    By Width
    GET http://localhost:7000/p/6e0/072/682/e66287b662827da75b244a3/w200.jpg

![Image](test/images/wine/w200.jpg?raw=true)


    Square
    GET http://localhost:7000/p/6e0/072/682/e66287b662827da75b244a3/x200.jpg

![Image](test/images/wine/x200.jpg?raw=true)

    Rectangle
    GET http://localhost:7000/p/6e0/072/682/e66287b662827da75b244a3/300x200.jpg

![Image](test/images/wine/300x200.jpg?raw=true)


**Quality**

The default compression of the image can be modified by appending `-q` and the desired quality `1-100`.

    Square with quality 50
    GET http://localhost:7000/p/6e0/072/682/e66287b662827da75b244a3/x200-q30.jpg

![Image](test/images/wine/x200-q30.jpg?raw=true)


### Cloud Storage

Images will be uploaded to either Amazon S3 or Joyent's Manta

To store images in S3 the following flags need to be set
```shell
--aws_access_key_id $AWS_ACCESS_KEY_ID --aws_secret_key $AWS_SECRET_KEY --aws_bucket $AWS_BUCKET
```

For Manta the following flags are required
```shell
--manta_url $MANTA_URL --manta_user $MANTA_USER --manta_key_id $MANTA_KEY_ID --sdc_identity $SDC_IDENTITY --remote_base_path $IMG_MANTA_BASE_PATH
```

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
git clone https://github.com/wanelo/image-server.git $GOPATH/src/github.com/wanelo/image-server
ln -s $GOPATH/src/github.com/wanelo/image-server ~/workspace/image-server
cd ~/workspace/image-server
```

Install dependencies:

Go needs to be installed with cross compilation. Imagemagick will require giflib and webp support.

On Mac
```bash
brew install --force go --with-cc-all
brew install --force giflib
brew install --force imagemagick --with-webp
make deps
```

Set up editor:

  - Atom.io package [go-plus](https://github.com/joefitzgerald/go-plus)

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

The server allows to be profiled when started with the profile flag

```
bin/images --profile server
```

The profiling information is available on `localhost:6060`

It is important to run the profiler

You will need the profiled data from the server, and analize it with the same executable file used on the server.

You will need to download the profiled data from the server.
```
ssh example.com "curl http://localhost:6060/debug/pprof/heap" > images.pprof
```

Use `go tool pprof` to analize the profile. Remember to use the same executable file as the one on production.
```
go tool pprof --inuse_objects bin/solaris/images images.pprof
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
