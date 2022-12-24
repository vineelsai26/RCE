# RCE

RCE is a http API in Go with which you can run different programs written in different languages

## Getting Started

### With Docker

* From Docker Hub

```sh
docker run -it -e PORT=3000 -p 3000:3000 -v "/var/run/docker.sock:/var/run/docker.sock" -v "/usr/src/app/runs:/usr/src/app/runs" vineelsai/rce
```

* From GHCR

```sh
docker run -it -e PORT=3000 -p 3000:3000 -v "/var/run/docker.sock:/var/run/docker.sock" -v "/usr/src/app/runs:/usr/src/app/runs" ghcr.io/vineelsai26/rce:latest
```

***Note: Both images are the same***

### From command line

* Download the latest Binary from [Releases Page](https://github.com/vineelsai26/RCE/releases/latest) 

* Pull the required docker images

```sh
./rce-linux-amd64.bin --pull-images
```

* And run it with the command

```sh
./rce-linux-amd64.bin
```

## Usage

* To run your code just make a POST request to http://localhost:3000/run with body

```json
{
    "code": "print(\"Hello World\")",
    "language": "python"
}
```

* The code will get executed and return a JSON object

```json
{
    "output": "Hello World"
}
```

## Supported Languages

* python
* c
* cpp

## Help

* For help just run to know all the things you can do 

```sh
./rce-linux-amd64.bin --help
```
