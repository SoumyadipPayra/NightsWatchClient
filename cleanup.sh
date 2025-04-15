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

# Remove the config files
sudo rm -f /usr/local/etc/nightswatch/username /usr/local/etc/nightswatch/password

# Remove the metadata directory
sudo rm -rf /usr/local/etc/nightswatch

echo "NightsWatch client service and binary have been removed."