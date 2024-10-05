#!/bin/bash

# Create the directory structure
sudo mkdir -p /opt/emoncms-feedsync

# Copy the binary
sudo cp ./emoncms-feedsync-v* /opt/emoncms-feedsync/emoncms-feedsync
sudo cp ./sample-config.yaml  /opt/emoncms-feedsync/config.yaml

# Copy the service file
sudo cp ./sample-feedsync.service /etc/systemd/system/emoncms-feedsync.service

# Set appropriate permissions
sudo chmod 755 /opt/emoncms-feedsync/emoncms-feedsync
sudo chmod 644 /opt/emoncms-feedsync/config.yaml
sudo chmod 644 /etc/systemd/system/emoncms-feedsync.service

# Reload systemd to recognize the new service
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable emoncms-feedsync.service

# Start the service
sudo systemctl start emoncms-feedsync.service

echo "Installation complete. Emoncms-feedsync service is now running."
