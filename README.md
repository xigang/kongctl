## Overview

A CLI For Kong API Gateway &amp; Service Mesh

## Installation

Make sure you have a working Go environment. Go version 1.2+ is supported. [See the install instructions for Go.](https://golang.org/doc/install)

To install cli, simply run:

```
$ go get github.com/xigang/kongctl
$ make
```
Move kongctl binary to your `PATH`

## QuickStart

- [Getting-started](docs/getting-started.md)

## Features

- Support for CURD of upstream, target, service, route, consumer, plugin objects.
- Supports Basic Authentication and Statsd plugins.

## LICENSE

- [Apache License](LICENSE)