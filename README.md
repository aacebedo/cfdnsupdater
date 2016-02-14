# CloudFlare DNS Updater

This tool helps fools like me to cope with shitty ISP which does not provide a 
static IP to their customers. It automatically updates DNS settings of any 
domain contained in your cloudflare account. It will also monitor the server's 
public IP, if it changes, it will update the domain's DNS records to reflect it.

The tool can update multiple domains from a cloudflare account. 
For each domain a filter can be applied on DNS records to update (filter by name
 or by record type)
The program will update the IP address of each record with the public IP 
it previously obtained from a service named IPfy (https://api.ipify.org).

Logs can be output in stdout or syslog. 
Check the help for more information about how to enable it.

### Build
#### Requirements
- Docker >=1.9
- Go >=1.5

#### Build with docker 
Execute the following commands:
```sh
$ docker build -t cfdnsupdaterbuild ./environments/build
$ docker run -t -v <output_path_on_host>:/build cfdnsupdaterbuild -a <ARCH> \
       [-b <BRANCH or TAG>] -p aacebedo/cfdnsupdater -bn cfdnsupdater -o /build 
```
This command will create a tar gz and debian file in the build directory 
of the container.

#### Build with go

Execute the following commands in a environment with go command availabled
```sh
$ go get github.com/aacebedo/cfdnsupdater ./... 
```
 
### Installation
There are multiple ways to install the software. Package for ARM and amd64
available for each methods.

#### Archive 
Get an archive on github and execute the following commands:
```sh
$ tar xvzf ./cfdnsupdater.<ARCH>.<VERSION>.tar.gz
```
Note: I also provide a systemd file you can manually install in the 
appropriate directory ("/etc/systemd/system" for instance).

#### Debian repository
Add the following debian repository to your apt source and install it through
apt-get:
```sh
$ "deb https://dl.bintray.com/aacebedo/cfdnsupdater <DISTRIBUTION> main" | \
      sudo tee -a /etc/apt/sources.list
$ sudo apt update
$ sudo  apt install cfdnsupdater
```

#### Docker registry
Pull the docker image and run it in a container:
```sh
docker pull aacebedo-docker-cfdsnupdater.bintray.io/cfdnsupdater-amd64:<VERSION> 
```

### Use
```sh
$ cfdnsupdater -c <path_to_configuration_file>
```
Here is an example of a configuration file (also included in the distribution):
```sh
foo.com:
  email: user@fake.com
  apikey: get_your_api_key_from_your_cloud_flare_a_account
  period: 60
  record_types: []
  record_names: []
```
##### More Help
```sh
$ cfdnsupdater -h
```

### Detailed design
Coming soon

### License
LGPL 3