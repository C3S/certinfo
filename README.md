<!--
title: "`certinfo`"
subtitle: Examine the health status of your TLS certificates
nodate: true
toc: true
lang: en
latex_engine: pdflatex
output:
  # pdf_document:
  bookdown::pdf_document2:
    toc: true
    latex_engine: pdflatex
    template: C3S_standard_doc.latex
    df_print: kable
shortname: "`certinfo`"
abstract: |
  Quickly get an overview of your TLS certificates. Sure there's tons of monitoring solutions out there,
  but we figured this was still a one-purpose tool that comes in handy.
-->

# `certinfo` -- Examine the health status of your TLS certificates

This small tool has one purpose: To tell you when the TLS certificates of your websites will expire. If you have
dozens of projects with individual certificates, you might find this helpful.

All it does is to try and open a TCP connection to a given domain, using timeout, port and protocol as defined.
If successful, it will show you the number of days until certificate expiration in one line (per domain).


## Installation

If missing, first install Go on your computer. On Ubuntu or Debian, for example, use:

```bash
sudo apt install golang-go
```

Then check out the source code:

```bash
git clone git@github.com:C3S/certinfo.git
cd certinfo
# fetch dependencies
go get .
```

`certinfo` depends on these additional modules:
* [color](https://github.com/fatih/color) for colored output
* [cobra](https://github.com/spf13/cobra) for providing commands and flags
* [viper](https://github.com/spf13/viper) for handling configuration files

You should now already be able to run it:

```bash
go run .
```

For permanent use, you should build it and run the binary:

```bash
go build
./certinfo
```

It's a single binary file, you can link or copy it to a directory in your `PATH` environment variable to be able to call it from anywhere.


## Usage

There's two feedback variants: The command `expires` will print the expiration date, whereas `bestbefore` also
compares that date to a given threshold, defaulting to 14 days. That is, if the certificate is still valid
beyond the threshold, you'll get a green `ok`, if it will expire before that you'll be told `please renew` in orange,
and if the certificate is already expired, you'll get a `red alert`.

This can be done in two modes: `expires` and `bestbefore` expect a single domain as additional argument, so they
do one check exactly. The commands `hosts_expire` and `hosts_bestbefore`, however, fetch domain names and HTTPS ports
from a YAML configuration file. This can either be a file called `config.yaml` in the current directory or `~/.config/certinfo/`,
or a file specified via the global `--config` flag. See the `config.yaml.example` file, it should explain itself.
In this mode, both IPv4 and IPv6 will be tried for each domain.

Partly depending on the feedback variant, you can adjust the IP `--protocol` (4 or 6, the default), the TCP `--port` (default is 443),
the number of `--days` used by `bestbefore`/`hosts_bestbefore` (default is 14), and the `--timeout` for dials in seconds (default is 5).

Connection errors are silently ignored by default. To see them, use the `--errors` flag.

## Examples

```console
# replace this with one of your domains
foo@bar:~$ MYDOMAIN="www.example.com"

# check the expiration date of your domain (IPv6, port 443)
foo@bar:~$ certinfo expires "${MYDOMAIN}"
expires: 02.10.2023 (IPv6)

# alternatively, see if the certificate is good for at least another 14 days
foo@bar:~$ certinfo bestbefore "${MYDOMAIN}"
will expire in 76 days (IPv6)  -- ok!

# now check for IPv4 and 7 days
foo@bar:~$ certinfo bestbefore "${MYDOMAIN}" --protocol 4 --days 7
expires in 6 days (IPv4)  -- please renew!

# flags can be shortened; since -p is already used for --port, the
# IP version was shortened to -i
foo@bar:~$ certinfo bestbefore "${MYDOMAIN}" -i 4 -d 7
expires in 6 days (IPv4)  -- please renew!
```

For the next examples, prepare a configuration file with some of your subdomains similar to this one:

```yaml
hosts:
    - example:
        url: www.example.com
        port: 443
    - another:
        url: another.example.com
        port: 443
```

It can then be used with `hosts_expire` and `hosts_bestbefore`:

```console
# if you save the config to ~/.config/certinfo/config.yaml or the current
# directory, it will be found automatically
foo@bar:~$ certinfo hosts_expire
www.example.com            expires: 02.10.2023 (IPv4)
www.example.com            expires: 02.10.2023 (IPv6)
another.example.com        expires: 02.10.2023 (IPv4)
another.example.com        expires: 19.05.2023 (IPv6)

# the same with hosts_bestbefore, but this time provide the hosts configuration
# via the command line flag
foo@bar:~$ certinfo --config ~/my_config.yaml hosts_bestbefore
www.example.com            (IPv4): expires 02.10.2023, in 76 days     -- ok!
www.example.com            (IPv6): expires 02.10.2023, in 76 days     -- ok!
another.example.com        (IPv4): expires 02.10.2023, in 76 days     -- ok!
another.example.com        (IPv6): expired 19.05.2023, -59 days ago   -- red alert!
```


## Contributing

To ask for help, report bugs, suggest feature improvements, or discuss the global
development, please use the issue tracker.


### Branches

Please note that all development happens in the `develop` branch. Pull requests against the `main`
branch will be rejected, as it is reserved for the current stable release.


## Licence

Copyright 2023 m.eik michalke <meik.michalke@c3s.cc>.

These scripts are free software: you can redistribute them and/or modify
them under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

These scripts are distributed in the hope that they will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with these scripts.  If not, see <https://www.gnu.org/licenses/>.
