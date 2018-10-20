# Peer Server

Works along with the [peer](#) Sublime Text plugin.

It serves as a relay for content for peer programming sessions.

**Important Note:** This is WIP. Not intended for large scale usage yet. The current solution works **in memory** and it
is a traditional REST API.

## Future Steps

- [ ] Add authentication.
- [ ] Research and implement best DB solution and move out of in-memory solution.
- [ ] Research and implement websockets and move out of traditional RESTful to get contents of topic.

These points must be complete to tag an official version `1.0.0`. Everything before that would be a dev only `0.*.*` version.

## Recommendations

- Install the binary behind some sort of reverse proxy server such as [Nginx](https://nginx.org/en/).
- Configure TLS in your server.

## Requirements

- [Go](https://golang.org/dl/) language.

## Installation

1. Get the project `go get github.com/jltorresm/peerserver`.
1. Go to the project folder `cd $GOPATH/src/github.com/jltorresm/peerserver/`.
1. Install dependencies by running `dep ensure`.
1. Build the project `go build .`
1. Run `./peerserver`.
