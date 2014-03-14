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

- save processed images into manta
- ability to update the source of the image
- status page
  - current images processing count
  - current original download count
- configuration options/file
  - port number
  - status page port number
  - max dimensions
  - environments
    - S3 url
  - whitelist image extensions

