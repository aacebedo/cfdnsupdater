# CloudFlare DNS Updater

This tool help fools like me to cope with shitty ISP which does not provide a static IP to their customers.
It automatically updates DNS settings of any domain contained in your cloudflare account.
It will also monitors your public IP, if it changes, it will update your DNS records to reflect it.

### Requirements
- Python >= 3.1
- python-argparse
- python-argcomplete >= 1.0.0
- requests >= 2.9.1

### Installation
```sh
$ ./setup.py install
or
$ pip3 install cfdnsupdater
```
Note: I also provide a systemd file you can manually install wherever you like.

### Use
There are two ways to use this script:
##### Inline config use
Starting with a configuration passeed in the command line:
```sh
$ cfdnsupdater inlineconfig -e <cloudflare_account_email> -a <cloudflare_apikey> -t <record_types> <domain> <record_content>
```
##### File config use
Starting with a file containing the configuration:
```sh
$ cfdnsupdater fileconfig -c <path_to_configuration_file>
```
An example of a configuration file is provided.

##### Help
```sh
$ cfdnsupdater -h
```
### License
LGPL 3