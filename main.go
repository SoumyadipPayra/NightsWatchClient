package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	nightswatchclient "github.com/SoumyadipPayra/NightsWatchClient/client"
	"github.com/SoumyadipPayra/NightsWatchClient/enc_dec"
	"github.com/SoumyadipPayra/NightsWatchClient/osquery"
	nwPB "github.com/SoumyadipPayra/NightsWatchProtobufs/gogenproto/nightswatch"
)

var (
	username string
	password string
)

func main() {

	os.Mkdir(".metadata", 0755)
	// Check if installation file exists
	if _, err := os.Stat(".metadata/installed"); os.IsNotExist(err) {
		// Create installation file with default value
		err := os.WriteFile(".metadata/installed", []byte("false"), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create installation file: %v\n", err)
			os.Exit(1)
		}
	}

	installed := false
	if _, err := os.Stat(".metadata/installed"); err == nil {
		// File exists, read the value
		content, err := os.ReadFile(".metadata/installed")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to read installation status: %v\n", err)
			os.Exit(1)
		}
		installed = strings.Contains(string(content), "true")
	}

	if !installed {
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "Error: Expected username and password as arguments\n")
			fmt.Fprintf(os.Stderr, "Usage: %s <username> <password>\n", os.Args[0])
			os.Exit(1)
		}

		username = os.Args[1]
		password = os.Args[2]
		// Run installation mode
		installationClient, err := nightswatchclient.NewNightsWatchInstallationClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create installation client: %v\n", err)
			os.Exit(1)
		}
		defer installationClient.Close()

		err = installationClient.Register(context.Background(), &nwPB.RegisterRequest{
			Name:     username,
			Password: enc_dec.GenerateHash(password),
		})
		if err != nil && strings.Contains(err.Error(), "SQLSTATE 23505") {
			fmt.Println("User already exists")
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to register: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("--------------------------------")
		fmt.Println("Registration successful")
		fmt.Println("--------------------------------")
		// Create installation file
		err = os.WriteFile(".metadata/installed", []byte("true"), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write installation status: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(".metadata/username"); os.IsNotExist(err) {
			err = os.WriteFile(".metadata/username", []byte(username), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to write username: %v\n", err)
				os.Exit(1)
			}
		}

		if _, err := os.Stat(".metadata/password"); os.IsNotExist(err) {
			// Encrypt the password before storing
			encryptedPassword, err := enc_dec.Encrypt(password)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to encrypt password: %v\n", err)
				os.Exit(1)
			}

			err = os.WriteFile(".metadata/password", []byte(encryptedPassword), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to write password: %v\n", err)
				os.Exit(1)
			}
		}
		return
	}

	err := readCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read credentials: %v\n", err)
		os.Exit(1)
	}

	// Run initialization mode
	initClient, err := nightswatchclient.NewNightsWatchInitClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create initialization client: %v\n", err)
		os.Exit(1)
	}
	defer initClient.Close()

	err = initClient.Login(context.Background(), &nwPB.LoginRequest{
		Name:     username,
		Password: enc_dec.GenerateHash(password),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to login: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("--------------------------------")
	fmt.Println("Login successful")
	fmt.Println("--------------------------------")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	err = osquery.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize osquery: %v\n", err)
		os.Exit(1)
	}

	for range ticker.C {
		deviceData, err := osquery.GetSystemInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get device data: %v\n", err)
			continue
		}

		err = initClient.SendDeviceData(context.Background(), username, deviceData.ToPB())
		if err != nil {
			if strings.Contains(err.Error(), "invalid token") {
				fmt.Println("--------------------------------")
				fmt.Println("Login expired, logging in again")
				fmt.Println("--------------------------------")
				err = initClient.Login(context.Background(), &nwPB.LoginRequest{
					Name:     username,
					Password: enc_dec.GenerateHash(password),
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: Failed to login: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("--------------------------------")
				fmt.Println("Login successful")
				fmt.Println("--------------------------------")
				continue
			}

			fmt.Fprintf(os.Stderr, "Error: Failed to send device data: %v\n", err)
			continue
		}
		fmt.Println("--------------------------------")
		fmt.Println("Device data sent successfully at", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("--------------------------------")
	}

}

func readCredentials() error {
	// Read username and password from files
	name, err := os.ReadFile(".metadata/username")
	if err != nil {
		return fmt.Errorf("failed to read username: %v", err)
	}
	username = string(name)

	encryptedPassword, err := os.ReadFile(".metadata/password")
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}

	// Decrypt the password
	pwd, err := enc_dec.Decrypt(string(encryptedPassword))
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %v", err)
	}
	password = pwd

	return nil
}
