# Interplanetary File Dumpster

An [imageboard](https://en.wikipedia.org/wiki/Imageboard), but images are stored in a peer-to-peer network

![Screenshot](contrib/screenshot.jpeg)

Features:

- Easy file sharing without registration and SMS. Supports images, video and audio files.
- Files are not stored on the disk, but are uploaded to the peer-to-peer [IPFS](https://ipfs.io) network instead.
- IPFS gateway is configurable. Use your own or any other [public gateway](https://ipfs.github.io/public-gateway-checker/) to host uploaded files.
- RSS feeds.
- Comments and up/down votes.
- Premoderation mode toggle.

## Requirements

- PostgreSQL database
- IPFS node
- Go 1.15.2 (only for building)
- A web server is recommended, i.e. nginx

## Build

- git clone repo
- run `make`

## Install

- copy (or symlink) `static/` and `templates/` directories to your location, i.e. `/etc/ipfd`
- copy `config.example` to your location, i.e. `/etc/ipfd/config`
- edit the config file
- install database schema from `migrations/xxx_init.up.sql` file
- optionally, configure nginx as a reverse proxy

Look inside `contrib/` directory for systemd and nginx config files. 

## Run

`./ipfd -config /etc/ipfd/config`

## Config options

- SITE\_URL, website url, i.e. https://example.org
- SITE\_NAME, website title
- BIND\_ADDRESS, bind ipfd web server to this address
- DATABASE\_URL, postgresql connection string
- TEMPLATES\_PATH, path to templates, i.e. `/etc/ipfd/templates`
- STATIC\_URL, URL to static files
- STATIC\_PATH, path to static files, i.e. `/etc/ipfd/static`
- SERVE\_STATIC, if you want to serve static with nginx, set this to false
- IPFS\_API, IPFS node settings, by default it should be /ip4/127.0.0.1/tcp/5001
- IPFS\_GATEWAY, URL of IPFS gateway
- IPFS\_PIN, files are not pinned in the IPFS repository if set to false
- MAX\_FILESIZE, maximum uploaded file size in bytes
- ALLOWED\_CONTENT\_TYPES, which MIME file types are allowed
- PREMODERATION, if true, posts are not displayed until admin approves them manually
