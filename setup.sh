#!/bin/bash

# Build the binary
echo "Building NightsWatch client..."
go build -o nightswatchclient

# Create necessary directories
echo "Creating necessary directories..."
sudo mkdir -p /usr/local/bin
sudo mkdir -p /usr/local/etc/nightswatch

# Install the binary
echo "Installing binary..."
sudo cp nightswatchclient /usr/local/bin/
sudo chmod +x /usr/local/bin/nightswatchclient

# Create a launchd plist file for nightswatchclient
echo "Creating launchd service..."
cat > ~/Library/LaunchAgents/com.nightswatch.client.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.nightswatch.client</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/nightswatchclient</string>
        <string>--config</string>
        <string>/usr/local/etc/nightswatch/config.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/tmp/nightswatchclient.err</string>
    <key>StandardOutPath</key>
    <string>/tmp/nightswatchclient.out</string>
    <key>WorkingDirectory</key>
    <string>/usr/local/etc/nightswatch</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>
EOF

# Load the service
echo "Loading service..."
sudo launchctl load ~/Library/LaunchAgents/com.nightswatch.client.plist

# Start the service
echo "Starting service..."
sudo launchctl start com.nightswatch.client

echo "NightsWatch client service has been installed and started."
echo "Logs can be found in /tmp/nightswatchclient.out and /tmp/nightswatchclient.err"
echo "Configuration files will be stored in /usr/local/etc/nightswatch/"
