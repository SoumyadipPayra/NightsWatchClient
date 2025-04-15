#!/bin/bash

# Stop the service if it's running
launchctl stop com.nightswatch.client

# Unload the service
launchctl unload ~/Library/LaunchAgents/com.nightswatch.client.plist

# Remove the plist file
rm -f ~/Library/LaunchAgents/com.nightswatch.client.plist

# Remove log files
rm -f /tmp/nightswatchclient.out /tmp/nightswatchclient.err

# Remove the binary
sudo rm -f /usr/local/bin/nightswatchclient

# Remove all configuration and metadata files
sudo rm -rf /usr/local/etc/nightswatch
sudo rm -rf .metadata

echo "NightsWatch client service, binary, and all configuration files have been removed."