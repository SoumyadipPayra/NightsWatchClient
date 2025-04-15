package osquery

import nwPB "github.com/SoumyadipPayra/NightsWatchProtobufs/gogenproto/nightswatch"

type AppInfo struct {
	AppName    string
	AppVersion string
}

// SystemInfo represents the system information extracted from osquery
type SystemInfo struct {
	InstalledApps  []*AppInfo
	OSVersion      string
	OsqueryVersion string
}

// ToPB converts the SystemInfo struct to a DeviceDataRequest
func (s SystemInfo) ToPB() *nwPB.DeviceDataRequest {
	return &nwPB.DeviceDataRequest{
		OsqueryVersion: &nwPB.OSQueryVersion{
			Version: s.OsqueryVersion,
		},
		OsVersion: &nwPB.OSVersion{
			Version: s.OSVersion,
		},
		InstalledApps: &nwPB.InstalledApps{
			Apps: func() []*nwPB.App {
				apps := make([]*nwPB.App, len(s.InstalledApps))
				for i, app := range s.InstalledApps {
					apps[i] = &nwPB.App{
						Name:    app.AppName,
						Version: app.AppVersion,
					}
				}
				return apps
			}(),
		},
	}
}
