# About

[![Build Status](https://travis-ci.org/gongled/vgrepo.svg?branch=master)](https://travis-ci.org/gongled/vgrepo)

HashiCorp company does [Vagrant](https://www.vagrantup.com) for managing the lifecycle of virtual machines. It is great, but they do not 
provide any open source tools for versioning and discovering your own images without a necessity to have an account 
on [HashiCorp Atlas](https://atlas.hashicorp.com/help/intro/features-list). 

`vgrepo` is a simple CLI tool for managing Vagrant repositories. In pair with HTTP server it provides 
simple way to distribute your images without worries about manual upgrading them on your team.

## Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `vgrepo` from scratch, make sure you have a working Go 1.5+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/gongled/vgrepo
```

If you want update `vgrepo` to latest stable release, do:

```
go get -u github.com/gongled/vgrepo
```

## Configuration

Specify storage settings in the `/etc/vgrepo.conf` configuration file. `path` is a directory that contains 
repositories with their metadata: name, versions and providers of VMs. `url` uses to discover your images 
and provides a permanent link to metadata.

```
[storage]

  url: http://localhost:8080

  path: /srv/vagrant

```

Run NGINX with the following configuration of the virtual host. 
It allows you distribute your images for more than one machine over HTTP or HTTPS. 

```
server {
    listen 8080;
    server_name localhost;

    root /srv/vagrant;

    location ~ ^/([^\/]+)/$ {
        index /metadata/$1.json;
        try_files /$1/metadata/$1.json =404;
    }

    location ~ \.json$ {
        add_header Content-Type application/json;
    }

    location ~ \.box$ {
        add_header Content-Type application/octet-stream;
    }

    location / {
	    autoindex off;
        expires -1;
    }
}
```

Done.

For example, you have an image with the name `powerbox`. Now after adding you can specify URL 
`http://localhost:8080/powerbox` in the `config.vm.box_url` to force Vagrant checking 
updates every time you run command `vagrant up`. Sounds great!

## Usage

```
Usage: vgrepo {command} {options}

Commands

  add       Add image to the Vagrant repository
  list      Show the list of available images
  delete    Delete the image from the repository
  info      Display info of the particular repository
  help      Display the current help message

Options

  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  vgrepo add $HOME/powerbox-1.0.0.box powerbox 1.1.0
  Add image to the Vagrant repository

  vgrepo list
  Show the list of available images

  vgrepo remove powerbox 1.1.0
  Remove the image from the repository

  vgrepo info powerbox
  Remove the image from the repository
```

## License

[MIT](LICENSE)
