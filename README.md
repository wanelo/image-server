# Image Server

## Usage

The default port number is 7000

To start the server under development
```bash
go run *.go
```

To run the compiled executable
```
./bin/image-server
```

## Configuration

Example of json configuration on `config/production.json`
```json
{
  "server_port": "7000",
  "status_port": "7001",
  "source_domain": "http://cdn-s3-2.wanelo.com",
  "whitelisted_extensions": ["jpg", "png", "webp"],
  "maximum_width": 1000
}

```

### Sample images

**Image Types**

    User Avatar
    http://localhost:7000/user/avatar/3589782/w500.jpg

    Product
    http://localhost:7000/product/image/10855050/x400.jpeg


**Dimensions**

    Maximum Width
    http://localhost:7000/user/avatar/3589782/w500.jpg

    Square
    http://localhost:7000/user/avatar/3589782/x600.jpg

    Rectangle
    http://localhost:7000/user/avatar/3589782/300x400.jpg


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

## Pending
- accept flags with
  - environment `-e production`, default will be `development`
- save processed images into manta
- ~~ability to update the source of the image [done]~~
- only allow whitelisted formats
- status page
  - current images processing count
  - current original download count
- graphite events (https://github.com/marpaia/graphite-golang)
  - completed
  - completed with errors
  - downloaded source from s3
  - failed downloading from s3
  - downloaded source from manta
  - failed downloading from manta
  - extension
- ~~configuration options/file~~
  - ~~port number [done]~~
  - ~~status page port number [done]~~
  - ~~max dimensions [done]~~
  - ~~source domain [done]~~
  - ~~whitelist image extensions [done]~~
- keep track of the image dimension statistics by image type, dimension, and extension (product image, user avatar). When an image is requested, other popular sizes can be generated in the background after the request. Images created on the background should not count towards statistics.
