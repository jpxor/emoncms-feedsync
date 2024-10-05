# emoncms-feedsync
Synchronize data feeds between two instances of EmonCMS (eg. self-hosted --> emoncms.org).

**Note: EmonCMS has a built in feature to synchronize data between two instances. Try to use that first!**

This project is meant to be used when the built-in feature fails. Example, for an EmonCMS instance running in docker,
some of its services fail to start/run including the sync service.  I built this to run on my home server alongside 
the docker EmonCMS instance.

# What is EmonCMS?
https://emoncms.org/site/home

In their own words:
> Emoncms is a powerful open-source web-app for processing, logging and visualising energy, temperature and other environmental data

# Install and run on Linux server

Build from source or download and extract release from tar file:
```
# Download release: emoncms-feedsync-v###.tar.gz
# extract:
tar -xf emoncms-feedsync-v###.tar.gz

# or build:
go build ./cmd/emoncms-feedsync
```
Then update `sample-config.yaml` to add your EmonCMS api keys, and update the local host address and port number if needed:
```
local:
    host: localhost:8081
    apikey: [read-only-key]
remote:
    host: emoncms.org
    apikey: [read-write-key]
interval: 600
```
Then run the `sample-linux-install.sh` script, which will copy the binary and config files into `/opt/emoncms-feedsync` and start the systemd service. This will also ensure the service is restarted on boot:
```
chmod +x ./sample-linux-install.sh
./sample-linux-install.sh
```
Update all sample files to suit your needs.
