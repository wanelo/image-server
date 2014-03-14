# Image Server

## Development

Prerequisites

```bash
brew install imagemagick --with-webp
```

To download dependencies 
```bash
make deps
```

To build
```bash
make
```

To run server
```bash
go run *.go
```

## Pending
Configuration options

- max dimensions
- environments
  - S3 url
- whitelist image extensions

- status page
  - current images processing count
  - current original download count

