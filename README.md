<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# Janus

[![Build Status](https://travis-ci.org/hellofresh/janus.svg?branch=master)](https://travis-ci.org/hellofresh/janus)
[![GoDoc](https://godoc.org/github.com/hellofresh/janus?status.svg)](https://godoc.org/github.com/hellofresh/janus)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellofresh/janus)](https://goreportcard.com/report/github.com/hellofresh/janus)

> An API Gateway written in Go

This is a lightweight API Gateway and Management Platform that enables you to control who accesses your API,
when they access it and how they access it. API Gateway will also record detailed analytics on how your
users are interacting with your API and when things go wrong.

## Why Janus?

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Ianus, pronounced [ˈjaː.nus]) is the god of beginnings,
gates, transitions, time, doorways, passages, and endings. He is usually depicted as having two faces since he
looks to the future and to the past. [Wikipedia](https://en.wikipedia.org/wiki/Janus)

We thought it would be nice to name the project after the God of the Gates :smile:

## What is an API Gateway?

An API Gateway sits in front of your application(s) and/or services and manages the heavy lifting of authorisation,
access control and throughput limiting to your services. Ideally, it should mean that you can focus on creating
services instead of implementing management infrastructure. For example, if you have written a really awesome
web service that provides geolocation data for all the cats in NYC, and you want to make it public,
integrating an API gateway is a faster, more secure route than writing your own authorisation middleware.

## Key Features

This API Gateway offers powerful, yet lightweight features that allows fine gained control over your API ecosystem.

* **No dependency hell** - Single binary made with go
* **REST API** - Full programatic access to the internals makes it easy to manage your API users, keys and API Configuration from within your systems
* **Rate Limiting** - Easily rate limit your API users, rate limiting is granular and can be applied on a per-key basis
* **CORS Filter** - Enable cors for your API, or even for specific endpoints
* **API Versioning** - API Versions can be easily set and deprecated at a specific time and date
* **Multiple auth protocols** - Out of the box, we support JWT, OAuth 2.0 and Basic Auth access methods
* **Hot-reloading of configuration** - No need to restart the process
* **Graceful shutdown** - Graceful shutdown of http connections
* **Docker Image** - Small [official](https://quay.io/repository/hellofresh/janus) docker image included

## Installation

### Docker

The simplest way of installing janus is to run the docker image for it. Just check the [docker-compose.yml](ci/assets/docker-compose.yml)
example and then run it.

```sh
docker-compose up -d
```

Now you should be able to get a response from the gateway. 

Try the following command:

```sh
http http://localhost:8080/
```

### Manual

You can get the binary and play with it in your own environment (or even deploy it where ever you like).
Just go to the [releases](https://github.com/hellofresh/janus/releases) page and download the latest one for your platform.

Make sure you have the following dependencies installed in your environment:

 - Mongodb - For storing the proxies configurations

And then just define where your dependencies are located

```sh
export DATABASE_DSN="mongodb://localhost:27017/janus"
export SECRET="yourSecret"
export ADMIN_USERNAME="admin"
export ADMIN_PASSWORD="admin"
```

If you want you can have a `stastd` server so you can have some metrics about your gateway. For that just define:

```sh
export STATSD_DSN="statsd:8125"
export STATSD_PREFIX="janus."
```

## Getting Started

After you have *janus* up and running we need to setup our first proxy. Everything that we want to do on the gateway
we do it through a REST API, since all endpoints are protected, we need to login first.

```sh
http -v --json POST localhost:8080/login username=admin password=admin
```

The username and password are defined in an environmental variable called `ADMIN_USERNAME` and `ADMIN_PASSWORD`. It defaults to *admin*/*admin*.

<p align="center">
  <a href="http://g.recordit.co/dDjkyDKobL.gif">
    <img src="http://g.recordit.co/dDjkyDKobL.gif">
  </a>
</p>


### Creating a proxy

The main feature of the API Gateway is to proxy the requests to a different service, so let's do this.
Now that you are authenticated, you can send a request to `/apis` to create a proxy.

```
http -v POST localhost:8080/apis "Authorization:Bearer yourToken" "Content-Type: application/json" < examples/api.json
```

<p align="center">
  <a href="http://g.recordit.co/Hi7SX8s5IA.gif">
    <img src="http://g.recordit.co/Hi7SX8s5IA.gif">
  </a>
</p>

This will create a proxy to `https://jsonplaceholder.typicode.com/posts` when you hit the api gateway on `GET /posts`.
Now just make a request to `GET /posts`

```
http -v --json GET http://localhost:8080/posts/1
```
<p align="center">
  <a href="http://g.recordit.co/vufeMjwEfg.gif">
    <img src="http://g.recordit.co/vufeMjwEfg.gif">
  </a>
</p>

Done! You just made your first request through the gateway.

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation

* Janus Docs: https://godoc.org/github.com/hellofresh/janus
* Go lang: https://golang.org/