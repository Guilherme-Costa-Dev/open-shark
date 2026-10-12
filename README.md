# open-shark

CLI tool to configure the DPI and polling rate of an Attack Shark X11 mouse on Linux from a config file, and read its battery.

## Install

**Arch** — this repo ships a PKGBUILD. Build and install it locally:

```sh
makepkg -si
```

**Other distros** (needs Go with cgo and the `libudev` headers):

```sh
go build -o open-shark ./cmd/open-shark
sudo install -m 755 open-shark /usr/local/bin/open-shark
```

`makepkg -si` also installs a udev rule, so the tool runs without root. If you
install manually, add the rule yourself (see `PKGBUILD`) or run with `sudo`.

## Usage

```sh
mkdir -p ~/.config/open-shark
cp config.example.toml ~/.config/open-shark/config.toml   # then edit it

open-shark apply      # apply the config to the mouse
open-shark battery    # read the battery level
open-shark info       # device + battery
```

Edit `~/.config/open-shark/config.toml`:

```toml
[dpi]
stages = [800, 1600, 2400, 3200, 5000, 22000]   # 6 stages, 50-22000
active = 2                                      # active stage, 1-6

[polling]
rate = 1000                                     # 125, 250, 500 or 1000 Hz
```
## Disclaimer
This was all vibe coded, i dont really care about this code or what its doing, but it fixes my problem ig.
