# GoDriveLog Pi4 Fyne Kiosk Setup

Version: 0.1  
Target: Raspberry Pi 4 running a minimal Linux install that boots directly into the GoDriveLog Fyne display.

---

## 1. Goal

Set up a Raspberry Pi 4 as a single-purpose GoDriveLog in-vehicle display.

The target stack is:

```text
Raspberry Pi OS Lite 64-bit
  -> minimal X11 environment
  -> Go + Fyne build dependencies
  -> GoDriveLog daemon as a systemd service
  -> GoDriveLog Fyne display launched automatically on tty1
```

This avoids a full desktop environment while still supporting Fyne now.

---

## 2. Recommended OS

Use:

```text
Raspberry Pi OS Lite 64-bit
```

Reason:

- Official Raspberry Pi support.
- Minimal base image.
- Good Pi4 display, input, Bluetooth, Wi-Fi, USB and Mesa support.
- Debian-based, so package installation and systemd setup are predictable.
- Less yak-shaving than Alpine, Armbian, or custom kiosk distributions.

Avoid installing the full Raspberry Pi desktop unless you specifically want panels, menus, background services, and other tiny gremlins in waistcoats.

---

## 3. Update the Pi

Run:

```bash
sudo apt update
sudo apt full-upgrade -y
sudo reboot
```

After reboot:

```bash
sudo apt update
```

---

## 4. Install Go, Fyne and minimal X11 dependencies

Install the build tools, Fyne dependencies, Mesa/OpenGL support, and a minimal X11 kiosk environment:

```bash
sudo apt install -y \
  golang gcc git make pkg-config \
  libgl1-mesa-dev libgl1-mesa-dri mesa-utils \
  xorg-dev libxkbcommon-dev \
  xserver-xorg xinit openbox dbus-x11 \
  x11-xserver-utils unclutter
```

What this gives you:

- `golang`, `gcc`, `pkg-config`: Go/Fyne build chain.
- `libgl1-mesa-dev`, `libgl1-mesa-dri`, `mesa-utils`: graphics/OpenGL support.
- `xorg-dev`, `libxkbcommon-dev`: Fyne/Linux input and display dependencies.
- `xserver-xorg`, `xinit`: enough X11 to run the GUI.
- `openbox`: tiny window manager.
- `unclutter`: hides the mouse cursor.

No full desktop. No taskbar. No “helpful” notification daemon waving from the dashboard.

---

## 5. Create the GoDriveLog runtime user

Create a dedicated user:

```bash
sudo useradd -r -m -s /bin/bash godrivelog || true
```

Add it to the groups needed for serial, display, rendering and input access:

```bash
sudo usermod -aG dialout,video,input,render,tty godrivelog
```

Notes:

- `dialout` is usually needed for USB serial OBD adapters.
- `video` and `render` are needed for graphics access.
- `input` may be needed for touchscreen or input devices.
- `tty` is useful for console display/session behaviour.

Reboot after group changes:

```bash
sudo reboot
```

---

## 6. Prepare GoDriveLog directories

Create config and data directories:

```bash
sudo mkdir -p /etc/godrivelog
sudo mkdir -p /var/lib/godrivelog
sudo chown -R godrivelog:godrivelog /var/lib/godrivelog
```

Create or copy your config file:

```bash
sudo nano /etc/godrivelog/config.yaml
```

Adjust the config path later if your app expects a different file.

---

## 7. Build GoDriveLog with Fyne support

From your GoDriveLog repository on the Pi:

```bash
go mod tidy
go build -tags fyne -o bin/godrivelog ./cmd/godrivelog
```

Install the binary:

```bash
sudo install -m 0755 bin/godrivelog /usr/local/bin/godrivelog
```

Check it exists:

```bash
/usr/local/bin/godrivelog --help
```

Expected display command once the Fyne renderer exists:

```bash
/usr/local/bin/godrivelog display fyne --config /etc/godrivelog/config.yaml
```

Important: this Pi setup makes the system Fyne-ready now. The `display fyne` command still needs to exist in the GoDriveLog codebase.

---

## 8. Create the daemon systemd service

Create the service file:

```bash
sudo tee /etc/systemd/system/godrivelog.service >/dev/null <<'EOF'
[Unit]
Description=GoDriveLog OBD logger daemon
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/godrivelog daemon --config /etc/godrivelog/config.yaml
Restart=always
RestartSec=5
User=godrivelog
Group=godrivelog
WorkingDirectory=/var/lib/godrivelog

[Install]
WantedBy=multi-user.target
EOF
```

Enable and start it:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now godrivelog.service
```

Check status:

```bash
systemctl status godrivelog.service
```

View logs:

```bash
journalctl -u godrivelog.service -f
```

Design rule: keep the daemon separate from the display. The daemon owns OBD and data capture. The display shows state. Do not let the GUI become the octopus holding all the spanners.

---

## 9. Create the Fyne kiosk startup script

Create `.xinitrc` for the `godrivelog` user:

```bash
sudo -u godrivelog tee /home/godrivelog/.xinitrc >/dev/null <<'EOF'
#!/bin/sh

xset s off
xset -dpms
xset s noblank

unclutter -idle 0.2 -root &

exec openbox-session &
sleep 1

exec /usr/local/bin/godrivelog display fyne --config /etc/godrivelog/config.yaml
EOF

sudo chmod +x /home/godrivelog/.xinitrc
```

This disables screen blanking and starts the GoDriveLog Fyne display inside X11.

---

## 10. Auto-start X on tty1

Add this to the `godrivelog` user profile:

```bash
sudo -u godrivelog tee -a /home/godrivelog/.profile >/dev/null <<'EOF'

if [ -z "$DISPLAY" ] && [ "$(tty)" = "/dev/tty1" ]; then
    startx -- -nocursor
fi
EOF
```

---

## 11. Enable auto-login on tty1

Edit the systemd override:

```bash
sudo systemctl edit getty@tty1
```

Paste:

```ini
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin godrivelog --noclear %I $TERM
```

Apply:

```bash
sudo systemctl daemon-reload
sudo systemctl restart getty@tty1
```

On next boot, the Pi should:

1. Boot to tty1.
2. Auto-login as `godrivelog`.
3. Run `startx`.
4. Launch the GoDriveLog Fyne display.

---

## 12. Test manually first

Before relying on auto-start, test as the `godrivelog` user:

```bash
sudo -iu godrivelog
startx
```

If X starts but the app fails, inspect:

```bash
cat /home/godrivelog/.xsession-errors
journalctl -xe
```

Also test graphics:

```bash
glxinfo | grep -i renderer
```

If `glxinfo` is missing:

```bash
sudo apt install -y mesa-utils
```

---

## 13. Useful troubleshooting commands

Check daemon:

```bash
systemctl status godrivelog.service
journalctl -u godrivelog.service -f
```

Check boot/login:

```bash
systemctl status getty@tty1
```

Check groups:

```bash
id godrivelog
```

Check OBD adapter:

```bash
ls -l /dev/ttyUSB* /dev/ttyACM* 2>/dev/null
dmesg | grep -i tty
```

Check display environment:

```bash
echo $DISPLAY
ps aux | grep -E 'Xorg|openbox|godrivelog'
```

---

## 14. Optional hardening later

Do not do these until the app is working reliably. Premature hardening is how a five-minute job becomes a small religion.

Later improvements:

- Read-only root filesystem.
- Separate writable `/var/lib/godrivelog`.
- Watchdog reboot.
- Automatic log rotation.
- Power-loss-safe shutdown strategy.
- Kiosk healthcheck service.
- Dedicated installer script.
- Systemd unit for display instead of `.profile` auto-start.
- Bluetooth OBD pairing automation if using Bluetooth.
- Hardware RTC if the vehicle is parked without network access.

---

## 15. Suggested repo branch for Fyne work

For the GoDriveLog code changes, branch from latest `main`:

```bash
git checkout main
git pull
git checkout -b add-fyne-display
```

Scope for that branch:

- Add the Fyne renderer.
- Keep it behind the `fyne` build tag.
- Add `godrivelog display fyne`.
- Do not mix in daemon rewrites, OBD changes, or packaging work.

Keep the branch small. Small branches are like sharp chisels: useful, controllable, and less likely to remove a thumb.

---

## 16. References

- Raspberry Pi OS downloads: https://www.raspberrypi.com/software/operating-systems/
- Fyne Linux setup: https://docs.fyne.io/started/quick/
