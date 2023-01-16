# Ravelin Example 3DS Checkout Implementation

This is an example 3DS implementation using Ravelin's 3D Secure API.

It demonstrates how to integrate with Ravelin's 3D Secure API and offer 3DS Secure in a checkout scenario.

For more detail, see [Ravelin's 3D Secure documentation](https://developer.ravelin.com/guides/3d-secure/).

## Requirements

- In order to test against Ravelin's 3DS test infrastructure, the client IP address needs to be allowed in Ravelin's firewall. This will be done as part of the integration process.

## Build and Run

This project is built with Go. 

To build and run the project, first download and install Go from golang.org.

From the root of the repository:
```shell
go build
./ravelin-3ds-demo -ravelin-api-key=<replace-with-api-key>
```

**Alternatively the project can be run from a docker container.**

From the root of the repository:
```shell
docker build -t ravelin-3ds-demo .
docker run -p 8085:8085 ravelin-3ds-demo -ravelin-api-key=<replace-with-api-key>
```

### Command Line Arguments

| Command Line Argument | Description |
| --- | --- | 
| `-ravelin-api-key` | Your Ravelin Sandbox API Key, accessible from the Ravelin Dashboard. <br> Test cards only work with sandbox accounts. <br> See [documentation](https://developer.ravelin.com/apis/authentication/) for more details. |
| `-ravelin-api-url` | The URL of the Ravelin 3DS API. <br> Defaults to https://pci.ravelin.com. |
| `-merchant-api` | The hostname the example 3DS implementation project is using. <br> This is used for API calls between the front-end and the back-end. <br> Defaults to http://localhost:8085. |
