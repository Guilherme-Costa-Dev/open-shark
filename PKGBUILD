# Local package: builds from this directory. Not intended for the AUR.
#
#   makepkg -si
#
pkgname=open-shark
pkgver=1.0.0
pkgrel=1
pkgdesc="Configure an Attack Shark X11 mouse (DPI, polling rate, battery) from a config file"
arch=('x86_64')
url="https://github.com/Guilherme-Costa-Dev/open-shark"
license=('unknown')
depends=('glibc' 'gcc-libs' 'systemd-libs')
makedepends=('go')
install="$pkgname.install"
source=()

build() {
  cd "$startdir"
  export CGO_ENABLED=1
  export GOFLAGS="-trimpath -buildvcs=false"
  go build -o open-shark ./cmd/open-shark
}

package() {
  cd "$startdir"
  install -Dm755 open-shark "$pkgdir/usr/bin/$pkgname"
  install -Dm644 config.example.toml "$pkgdir/usr/share/doc/$pkgname/config.example.toml"

  # udev rule: grant the logged-in user access to the mouse's configuration
  # interface so no root is needed at runtime.
  install -d "$pkgdir/usr/lib/udev/rules.d"
  cat > "$pkgdir/usr/lib/udev/rules.d/99-attack-shark-x11.rules" <<'EOF'
# Attack Shark X11 configuration interface.
# fa60 = 2.4 GHz dongle, fa55 = X11 wired, fa61 = R1 wired.
SUBSYSTEM=="hidraw", ATTRS{idVendor}=="1d57", ATTRS{idProduct}=="fa60", TAG+="uaccess"
SUBSYSTEM=="hidraw", ATTRS{idVendor}=="1d57", ATTRS{idProduct}=="fa55", TAG+="uaccess"
SUBSYSTEM=="hidraw", ATTRS{idVendor}=="1d57", ATTRS{idProduct}=="fa61", TAG+="uaccess"
EOF
}
