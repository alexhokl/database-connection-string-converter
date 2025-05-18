package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/alexhokl/database-connection-string-converter/dadbod"
	"github.com/alexhokl/database-connection-string-converter/sqls"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// dadbodToSqlsCmd represents the dadbodToSqls command
var dadbodToSqlsCmd = &cobra.Command{
	Use:   "dadbod-to-sqls",
	Short: "Convert dadbod config (from stdin) to sqls config (to stdout)",
	RunE:  runDadbodToSqls,
}

func init() {
	rootCmd.AddCommand(dadbodToSqlsCmd)
}

func runDadbodToSqls(cmd *cobra.Command, args []string) error {
	// Read JSON data from stdin
	var dadbodConnections []dadbod.Connection
	json.NewDecoder(os.Stdin).Decode(&dadbodConnections)

	// Transform the connections
	sqlsConnections := []sqls.Connection{}

	for _, connection := range dadbodConnections {
		sqlsConnection, err := ParseDadbodConnection(connection.URL)
		if err != nil {
			return fmt.Errorf("error parsing dadbod connection string: %w", err)
		}
		if sqlsConnection != nil {
			sqlsConnection.Alias = connection.Name
			sqlsConnections = append(sqlsConnections, *sqlsConnection)
		}
	}

	// Create YAML configuration
	ymlConfig := sqls.Config{
		LowercaseKeywords: false,
		Connections:       sqlsConnections,
	}

	// Encode YAML data to stdout
	yamlData, err := yaml.Marshal(ymlConfig)
	if err != nil {
		return fmt.Errorf("error encoding YAML for sqls: %w", err)
	}
	fmt.Println(string(yamlData))
	return nil
}

func ParseDadbodConnection(connURL string) (*sqls.Connection, error) {
	u, err := url.Parse(connURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
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

	// Create the SqlsConnection struct
	sqlsConn := &sqls.Connection{
		Driver:  u.Scheme,
		Proto:   "tcp",
		User:    username,
		Passwd:  password,
		Host:    host,
		Port:    port,
		DBName:  database,
		Options: options,
	}
	return sqlsConn, nil
}
