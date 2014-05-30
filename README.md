# Image Server

[![Build Status](https://magnum.travis-ci.com/wanelo/image-server.svg?token=xxYxjHDAXkDK41qZ1dqA&branch=master)](https://magnum.travis-ci.com/wanelo/image-server)

## Usage

The default port number is 7000

To start the server under development
```bash
make run
```

To run the compiled executable
```
./bin/image-server -e production
```

## Configuration

All configuration is passed by flags

## Image Generation

### Sample images

**Image Types**

    User Avatar
    GET http://localhost:7000/user/avatar/3589782/w500.jpg

    Product
    GET http://localhost:7000/product/image/10855050/x400.jpeg


**Dimensions**

    Maximum Width
    GET http://localhost:7000/user/avatar/3589782/w500.jpg

    Square
    GET http://localhost:7000/user/avatar/3589782/x600.jpg

    Rectangle
    GET http://localhost:7000/user/avatar/3589782/300x400.jpg

**Quality**

The default compression of the image can modified by appending `-q` and the desired quality `1-100`.

    Square with quality 50
    GET http://localhost:7000/user/avatar/3589782/x600-q50.jpg

### Multi Size Processing

This is useful for pre-generating images and saving them to the configured file store.
The request returs success after the original image is the file store.
Image outputs are generated after the request is complete.

    POST http://localhost:7000/user/avatar/3589782?outputs=x300.jpg,x300.webp&source=http://example.com/image.jpg

## CLI

Allows to create a range of images in parallel
```shell
images -namespace p -outputs x300.jpg,x300.webp -start 1000000 -end 1001000
```

## Manta CLI

Allows to download a range of images in parallel
```shell
images-manta -start 10000000 -end 10001000 -concurrency 100
```



## Error Handling

Few errors will cause the server to return error pages

- Source image is not found: NotFound (404)
- Image requested is larger than maximum_width: NotAcceptable (406)

## Development

Prerequisites

```bash
brew install imagemagick --with-webp
```

To download dependencies
```bash
make deps
```

## Compilation

`make` will build the executable under `./bin`
```bash
make
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

## Pending

### Required
- Strip metadata
- Default background color needs to be white (for transparent gifs, etc)

### After Release

- Zero-downtime restart: http://rcrowley.org/talks/strange-loop-2013.html#27
- Configuration reload with signal
- Status page
  - current images processing count
  - current original download count
