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

## Getting started

1. Specify storage settings in the `/etc/vgrepo.conf` configuration file. Parameter `path` is a 
directory that contains repositories with their metadata: name, versions and providers of VMs. 
Parameter `url` is used to discover your images and provides a permanent link to metadata.

    ```
    [storage]
    
      # Repository URL and port
      url: http://vagrant.example.tld
    
      # Repository path to store images and metadata
      path: /srv/storage
    ```
    
2. Create directory for the repository path `/srv/storage` and make sure that it is writable.

3. Add the image to the repository:

    ```
    vgrepo add /path/to/image.box powerbox 1.0.0 virtualbox
    ```
    
4. Configure NGINX to serve static files from `/srv/storage` directory.

    ```
    server {
        listen 80;
        
        server_name vagrant.example.tld;
        
        access_log off;
        error_log off;
        
        root /srv/storage;
        
        location / {
            autoindex on;
            expires -1;
        }
    }
    ```
 
Done. After adding changes you can specify URL `http://vagrant.example.tld/metadata/powerbox/powerbox.json` in 
the `config.vm.box_url` to force Vagrant checking updates every time you run command `vagrant up`. 

## Advanced

Imagine you have an image with the name `powerbox`. The standard path for metadata will be 
`http://vagrant.example.tld/metadata/powerbox/powerbox.json`, however it looks awful and unmemorable. 
You can use well-looking URL instead of direct link to JSON metadata file with the following NGINX 
configuration of the virtual host:  

```
server {
    listen 8080;
    
    server_name vagrant.example.tld;

    root /srv/storage;

    location ~ ^/r/([^\/]+)$ {
        return 301 $uri/;
    }

    location ~ ^/r/([^\/]+)/$ {
        index /metadata/$1/$1.json;
        try_files /metadata/$1/$1.json =404;
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

Now you are able to distribute your images for more than one machine over HTTP or HTTPS 
with a short and nice URLs with a format `http://vagrant.example.tld/r/powerbox`. 

## Usage

```
Usage: vgrepo {options} {command}

Commands

  add source name version provider    Add image to the Vagrant repository
  list                                Show the list of available images
  delete name version provider        Delete the image from the repository
  info name                           Display info of the particular repository
  render template output              Create index by given template file
  help                                Display the current help message

Options

  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  vgrepo add $HOME/powerbox-1.0.0.box powerbox 1.1.0 virtualbox
  Add image to the Vagrant repository

  vgrepo list
  Show the list of available repositories

  vgrepo delete powerbox 1.1.0
  Remove the image from the repository

  vgrepo info powerbox
  Show detailed info about the repository

  vgrepo render index.html /etc/vgrepo/templates/default.tpl
  Create index file by given template with output index.html
```

## License

[MIT](LICENSE)
