# open-shark

CLI tool to configure the DPI and polling rate of an Attack Shark X11 mouse on Linux from a config file, and read its battery.

## Install

Requires Go with cgo and the `libudev` headers.

```sh
sudo apt install libudev-dev   # Debian/Ubuntu
sudo pacman -S systemd-libs    # Arch

git clone https://github.com/Guilherme-Costa-Dev/open-shark.git 
cd open-shark
go build -o ~/.local/bin/open-shark ./cmd/open-shark
```

## Usage

```sh
mkdir -p ~/.config/open-shark
cp config.example.toml ~/.config/open-shark/config.toml   # then edit it

sudo open-shark apply      # apply the config to the mouse
sudo open-shark battery    # read the battery level
sudo open-shark info       # device + battery
```

`sudo` is required because `/dev/hidraw*` is owned by root.

Edit `~/.config/open-shark/config.toml`:

```toml
[dpi]
stages = [800, 1600, 2400, 3200, 5000, 22000]   # 6 stages, 50-22000
active = 2                                      # active stage, 1-6

[polling]
rate = 1000                                     # 125, 250, 500 or 1000 Hz
```
