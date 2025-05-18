package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/alexhokl/database-connection-string-converter/dadbod"
	"github.com/alexhokl/database-connection-string-converter/vscode"
	"github.com/spf13/cobra"
)

// dadbodToVscodeCmd represents the dadbodToVscode command
var dadbodToVscodeCmd = &cobra.Command{
	Use:   "dadbod-to-vscode",
	Short: "Convert dadbod config (from stdin) to VS Code config (to stdout)",
	RunE:  runDadbodToVscode,
}

func init() {
	rootCmd.AddCommand(dadbodToVscodeCmd)
}

func runDadbodToVscode(cmd *cobra.Command, args []string) error {
	// Read JSON data from stdin
	var dadbodConnections []dadbod.Connection
	if err := json.NewDecoder(os.Stdin).Decode(&dadbodConnections); err != nil {
		return fmt.Errorf("error decoding dadbod connections: %w", err)
	}

	// Transform the connections
	vscodeConnections := []vscode.Connection{}

	for _, connection := range dadbodConnections {
		vscodeConnection, err := ParseDadbodConnectionToVscode(connection.URL)
		if err != nil {
			return fmt.Errorf("error parsing dadbod connection string: %w", err)
		}
		if vscodeConnection != nil {
			vscodeConnection.ProfileName = connection.Name
			vscodeConnections = append(vscodeConnections, *vscodeConnection)
		}
	}

	// Encode JSON data to stdout
	jsonData, err := json.MarshalIndent(vscodeConnections, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding JSON for VS Code: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

// ParseDadbodConnectionToVscode parses a dadbod connection string and converts it to a VS Code connection object.
func ParseDadbodConnectionToVscode(connURL string) (*vscode.Connection, error) {
	u, err := url.Parse(connURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	if u.Scheme != "sqlserver" {
		return nil, nil
	}

	username := u.User.Username()
	password, _ := u.User.Password()
	host := u.Hostname()

	// Extract port
	portStr := u.Port()
	if portStr == "" {
		switch u.Scheme {
		case "mysql":
			portStr = "3306"
		case "postgresql":
			portStr = "5432"
		case "sqlserver":
			portStr = "1433"
		case "redis":
			portStr = "6379"
		case "sqlite":
			portStr = "0" // SQLite doesn't use a port

		default:
			return nil, fmt.Errorf("unsupported dialect: %s", u.Scheme)
		}
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	// Extract database name
	database := strings.TrimPrefix(u.Path, "/")

	// Extract options
	options := make(map[string]string)
	for key, values := range u.Query() {
		if len(values) > 0 {
			options[key] = values[0]
		}
	}

	// Create a new VS Code connection object
	vscodeConnection := &vscode.Connection{
		Server:                 fmt.Sprintf("%s,%d", host, port),
		Database:               database,
		AuthenticationType:     "SqlLogin",
		User:                   username,
		Password:               password,
		EmptyPasswordInput:     false,
		SavePassword:           false,
		TrustServerCertificate: true,
		Encrypt:                "Mandatory",
		ID:                     "",
		GroupID:                "",
	}

	return vscodeConnection, nil
}
