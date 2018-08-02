# Boot to RasPi stage intercom

This is a simple overview to scratch install RasPi stage intercom (based on [talkiepi](https://github.com/dchote/talkiepi/)) on your Raspberry Pi, and have it start on boot. 
This document assumes that you have raspbian-jessie-lite installed on your SD card, and that the distribution is up to date. 
This document also asumes that you have already configured network/wifi connectivity on your Raspberry Pi. 
(hint: add something like " ip=192.168.1.200::192.168.1.1:255.255.255.0:rpi:eth0:on" to cmdline.txt and create an empty file called "ssh" to get started.) 

By default it will run without any arguments, it will autogenerate a username and then connect to the Talkipi mumble server. 
You should change this behavior by appending commandline arguments of e.g. `-server 192.168.1.1:64738`,`-altserver 192.168.1.1:64738` , `-username CAM_A` to the ExecStart line in `/etc/systemd/system/mumble.service` once installed. 

You can set -server and  -altserver to the same address. If server can not be reached, altserver is tried without any delay.
A common setup is for -server to be the Ethernet-address and -altserver to be the address of the same server in Wifi.
-server will always be tried first.

It will also accept arguments for `-password`, `-insecure`, `-certificate` and `-channel`, all defined in `cmd/talkiepi/main.go`, if you run your own mumble server, these will be self explanatory. 


## TL:DR


There is a script [INSTALL_INTERCOM](INSTALL_INTERCOM) that does nearly all of these steps.  
In addition INTALL_VOLUME_KEYS allowy you to also control the ALSA output level of the headphone via 2 buttons on GPIO-pins. 
It assumes a file /boot/INTERCOM_SERVER.txt with the content of /conf/systemd/mumble.service already adapted to your setup (server IP, username, password).

When you are done, see [here](https://learn.adafruit.com/read-only-raspberry-pi/) about a nice script to make the file system read-only (unless a jumper is set at boot time).
GPIO21 (last in the row) looks like a good choice for the R/W jumper.

## Create a user

As root on your Raspberry Pi (`sudo -i`), create a mumble user:
```
adduser --disabled-password --disabled-login --gecos "" mumble
usermod -a -G cdrom,audio,video,plugdev,users,dialout,dip,input,gpio mumble
```

## Install

As root on your Raspberry Pi (`sudo -i`), install golang and other required dependencies, then build talkiepi:
```
apt-get install golang libopenal-dev libopus-dev git

su mumble

mkdir ~/gocode
mkdir ~/bin

export GOPATH=/home/mumble/gocode
export GOBIN=/home/mumble/bin

cd $GOPATH

# formerly: go get github.com/layeh/gopus
go get layeh.com/gopus
go get github.com/dchote/gopus
go get github.com/MarcusWolschon/RasPi_stage_intercom

cd $GOPATH/src/github.com/MarcusWolschon/RasPi_stage_intercom

go build -o /home/mumble/bin/talkiepi cmd/talkiepi/main.go 
```


## Start on boot

As root on your Raspberry Pi (`sudo -i`), copy mumble.service in to place:
```
cp /home/mumble/gocode/src/github.com/MarcusWolschon/RasPi_stage_intercom/conf/systemd/mumble.service /etc/systemd/system/mumble.service

systemctl enable mumble.service
```

## Create a certificate

This is optional, mainly if you want to register your talkiepi against a mumble server and apply ACLs.
```
su mumble
cd ~

openssl genrsa -aes256 -out key.pem
```

Enter a simple passphrase, its ok, we will remove it shortly...

```
openssl req -new -x509 -key key.pem -out cert.pem -days 1095
```

Enter your passphrase again, and fill out the certificate info as much as you like, its not really that important if you're just hacking around with this.

```
openssl rsa -in key.pem -out nopasskey.pem
```

Enter your password for the last time.

```
cat nopasskey.pem cert.pem > mumble.pem
```

Now as root again (`sudo -i`), edit `/etc/systemd/system/mumble.service` appending `-username USERNAME_TO_REGISTER -certificate /home/mumble/mumble.pem` at the end of `ExecStart = /home/mumble/bin/talkiepi`

Run `systemctl daemon-reload` and then `service mumble restart` and you should be set with a tls certificate!


## Install the USB soundcard

Because the raspberry pi does not come with a microphone port, we use a common USB soundcard to do that job.
You will need to change the default system sound device.
As root on your Raspberry Pi (`sudo -i`), find your device by running `aplay -l`, take note of the index of the device (likely 1) and then edit the alsa config (`/usr/share/alsa/alsa.conf`), changing the following:
```
defaults.ctl.card 1
defaults.pcm.card 1
```
_1 being the index of your device_


If your headset is too quiet, you can adjust the volume using amixer as such:
```
amixer -c 1 set Headphone 60%
```
or 
```
amixer -c 1 set Speaker 60%
```
_1 being the index of your device_


I will be adding volume control settings in an upcoming push.

## install mumble server

First we need to get a mumble server. e.g. from https://wiki.natenom.de/mumble/mumble-herunterladen