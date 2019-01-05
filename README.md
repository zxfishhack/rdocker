# rdocker

rdocker is a docker image for provide http endpoint from `unix:///var/run/docker.sock`.

## Installation

```bash
dep ensure -vendor-only
go build
docker build -t rdocker .
docker run -d -p 8080 -v /var/run/docker.sock:/var/run/docker.sock rdocker
```

## Usage

1. without environment
    ```
    docker -H tcp://<host>:8080 [command]
    ```
1. with environment
    
    set `DOCKER_HOST` to `tcp://<host>:8080`  
    ```
    docker [command]
    ```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)