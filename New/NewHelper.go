package new

import (
	"strings"
)

// PrebuiltApps Data
type preBuiltApp struct {
	Label    string
	Name     string
	UrlSSH   string
	UrlHTTPS string
}

// PreBuiltApps App data
var PreBuiltApps = []preBuiltApp{
	{Label: "Blank", Name: "bindu-blank", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-blank.git"},
	{Label: "Basic Web", Name: "bindu-basic-web", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-basic-web.git"},
	{Label: "Basic Api", Name: "bindu-basic-api", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-basic-api.git"},
	{Label: "Blog", Name: "bindu-blog", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-blog.git"},
	{Label: "E-Commerce", Name: "bindu-e-commerce", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-e-commerce.git"},
	{Label: "GRPC Server", Name: "bindu-grpc-server", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-grpc-server.git"},
	{Label: "GRPC Client", Name: "bindu-grpc-client", UrlSSH: "", UrlHTTPS: "https://github.com/bindu-bindu/bindu-grpc-client.git"},
	{Label: "Download Third Party Project", Name: "download", UrlSSH: "", UrlHTTPS: ""},
}

// DbAdapters Predefine DB adapters
var DbAdapters = []string{"Mysql", "Sqlite", "PGSql", "MongoDB", "None"}

// bindu new Hello --app Blank --db 'AdapterName:HostName:Port:DbName:DbUserName: DbPass' --port 8080
// bindu new Hello --app blank --port 9999 --db Mysql:Localhost:3306:blog:root

// FindPrebuitAppIndex find prebuilt app index
func FindPrebuitAppIndex(keyword string) (int, bool) {
	for i, app := range PreBuiltApps {

		if strings.ToLower(app.Label) == strings.ToLower(keyword) || strings.ToLower(app.Name) == strings.ToLower(keyword) {
			return i, true
		}
	}
	if strings.Contains(strings.ToLower(keyword), "https://") {
		return len(PreBuiltApps) - 1, true
	}
	return 0, false
}
