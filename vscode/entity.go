package vscode

type Connection struct {
	Server                 string `json:"server"`
	Database               string `json:"database"`
	AuthenticationType     string `json:"authenticationType"`
	User                   string `json:"user"`
	Password               string `json:"password"`
	EmptyPasswordInput     bool   `json:"emptyPasswordInput"`
	SavePassword           bool   `json:"savePassword"`
	ProfileName            string `json:"profileName"`
	ID                     string `json:"id"`
	Encrypt                string `json:"encrypt"`
	TrustServerCertificate bool   `json:"trustServerCertificate"`
	GroupID                string `json:"groupId"`
}
