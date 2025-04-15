package osquery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func Initialize() error {
	// Check if osqueryi is available
	_, err := exec.LookPath("osqueryi")
	if err != nil {
		return fmt.Errorf("osqueryi not found: %v", err)
	}
	return nil
}

// GetSystemInfo retrieves system information using osqueryi
func GetSystemInfo() (*SystemInfo, error) {
	// Get app name and version
	appInfo, err := runOsqueryQuery("SELECT name, bundle_short_version as version FROM apps")
	if err != nil {
		return nil, fmt.Errorf("failed to get app info: %v", err)
	}

	// Get OS version
	osInfo, err := runOsqueryQuery("SELECT version FROM os_version")
	if err != nil {
		return nil, fmt.Errorf("failed to get OS version: %v", err)
	}

	// Get osquery version
	osqueryInfo, err := runOsqueryQuery("SELECT version FROM osquery_info")
	if err != nil {
		return nil, fmt.Errorf("failed to get osquery version: %v", err)
	}

	return &SystemInfo{
		InstalledApps:  parseAppInfo(appInfo),
		OSVersion:      parseOSVersion(osInfo),
		OsqueryVersion: parseOsqueryVersion(osqueryInfo),
	}, nil
}

// runOsqueryQuery executes a query using osqueryi and returns the results
func runOsqueryQuery(query string) ([]map[string]string, error) {
	cmd := exec.Command("osqueryi", "--json", query)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run osquery query: %v", err)
	}

	var results []map[string]string
	if err := json.Unmarshal(out.Bytes(), &results); err != nil {
		return nil, fmt.Errorf("failed to parse osquery output: %v", err)
	}

	return results, nil
}

// parseAppInfo parses the app info from osquery output
func parseAppInfo(results []map[string]string) []*AppInfo {
	apps := make([]*AppInfo, len(results))
	for i, result := range results {
		apps[i] = &AppInfo{
			AppName:    result["name"],
			AppVersion: result["version"],
		}
	}
	return apps
}

// parseOSVersion extracts OS version from osquery output
func parseOSVersion(results []map[string]string) string {
	if len(results) > 0 {
		return results[0]["version"]
	}
	return ""
}

// parseOsqueryVersion extracts osquery version from osquery output
func parseOsqueryVersion(results []map[string]string) string {
	if len(results) > 0 {
		return results[0]["version"]
	}
	return ""
}
