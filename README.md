# NightsWatchClient

Go client to collect os and sys info of the device and push it to the server in a regular interval

# Client
- Grpc client to connect to the server
- Exposes two interfaces, one for registration during installation, another for login and sending device info
- Both the interfaces are implemented using the same impl class

# Enc_Dec
- Provides Basic encryption decryption functionality
- Uses AES to encrypt
- It is used to encrypt the password while storing in the device storage
- password loaded on the runtime, decrypted and a hash is generated and this hash is treated as the password in the server

# OsQuery
## Models
- Has the scehma deinition for the device data such as app_info, os_version and osquery_version
## Osquery
- Runs the osquery commands using `exec` and parse the retrieved info (maps) into the model structs

# Main
- The main package, that runs registration for the first time (installation)
- then it runs login(first time, followed by every token expiry)
- Then in a loop, it extracts device data in a regular interval and pushes it to the server

# Setup.sh
- bash file, can be treated as the installation file
- creates proper directories to store user name, password(encrypted), and installation status
- then builds the binary and launch it as service
- this service runs at system loading
- also keepalive is provided to start the program if halted by any
- and outputs and errors are captured in a log file
