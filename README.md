# CloudFlare DNS Updater

This tool helps fools like me to cope with shitty ISP which does not provide a static IP to their customers.
It automatically updates DNS settings of any domain contained in your cloudflare account.
It will also monitors your public IP, if it changes, it will update your DNS records to reflect it.

The tools can update only a specific domain from a cloudflare account. Within this domain it can filter records
by type or by name. The information updated is always the public IP of the server it gets from a service named IPfy (https://api.ipify.org).

Logs can be outputed in stdout or syslog. Check the help for more information about how to enable it.

### Build
#### Requirements
- Docker >=1.9

#### Instructions
```sh
$ docker build -t cfdnsupdaterbuild ./environments/build
$ docker run -t -v <output_path_on_host>:/out cfdnsupdaterbuild -a <ARCH> [-b <BRANCH or TAG>]
```

### Installation
```sh
$ tar xvzf ./cfdnsupdater.<ARCH>.<VERSION>.tar.gz
```
Note: I also provide a systemd file you can manually install wherever in the appropriate directory ("/etc/systemd/system" for instance).

### Use
##### Inline config use##### File config use
```sh
$ cfdnsupdater fileconfig -c <path_to_configuration_file>
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