[![Build Status](https://travis-ci.org/kovetskiy/eyed.svg?branch=master)](https://travis-ci.org/kovetskiy/eyed)

# eyed

![eye](https://cloud.githubusercontent.com/assets/8445924/9405831/a3d9e984-47ea-11e5-8f5a-d9a62a2530d3.png)

**eyed** it is the backend server for PAM module
[pam_eye](https://github.com/reconquest/pam_eye).
As mentioned in pam_eye documentation, pam module sends simple `GET` request to
**eyd** and uses local hostname as URL path.

## Installation

### Arch Linux
PKGBUILD and systemd unit available here
[archlinux](https://github.com/kovetskiy/eyed/tree/master/archlinux)

Building package:
```
$ git clone https://github.com/kovetskiy/eyed
$ cd archlinux
$ makepkg -fc
# pacman -U *.xz
```

Enabling and starting systemd service:
```
# systemctl enable eyed.service
# systemctl start eyed.service
```

### Other distros

**eyed** is go-getable, after running following command:
```
go get github.com/kovetskiy/eyed
```

**eyed** will be installed to `$GOPATH/bin`.


## Usage

**eyed** by default listening on `:80`, but you can set specified address
through passing argument `-l <listen>`.

**eyed** by default uses `/var/eyed/` folder as reports storage, for every
report, **eyed** appends the current date to a file with same name as hostname.
i.e. for request `GET http://eyed-host/node0.in.example.com` will be created (if not
exists) file `/var/eyed/node0.in.example.com`. Directory location can be
changed through passing argument `-d <directory>`.
