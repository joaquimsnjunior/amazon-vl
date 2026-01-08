package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"amazon-vl/internal/auth"
	"amazon-vl/internal/server"
)

var (
	version = "1.1.0"
	commit  = "dev"
)

func main() {
	// Define flags
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help information")
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("amazon-vl version %s (commit: %s)\n", version, commit)
		os.Exit(0)
	}

	// Handle help flag or missing arguments
	if *showHelp || flag.NArg() != 2 {
		printUsage()
		if *showHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Get positional arguments
	dir := flag.Arg(0)
	port := flag.Arg(1)

	// Validate directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("ERROR: Directory does not exist: %s", dir)
	}

	// Create server configuration
	cfg := server.Config{
		Dir:  dir,
		Port: port,
		Auth: auth.DefaultConfig(),
	}

	// Create and run server
	srv := server.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("ERROR: Server failed: %v", err)
	}
}

func printUsage() {
	fmt.Print(`
Amazon-VL - Secure Log File Server

USAGE:
    amazon-vl [OPTIONS] <directory> <port>

ARGUMENTS:
    <directory>    Path to the directory containing files to serve
    <port>         Port number to listen on (e.g., 8080, 9000)

OPTIONS:
    --help         Show this help message
    --version      Show version information

ENVIRONMENT VARIABLES:
    AUTH_USER      Username for authentication (default: joaquim)
    AUTH_HASH      MD5 crypt hash of password (default: hash for 'amazon')
    AUTH_REALM     Authentication realm (default: amazon-server-logs.com)

EXAMPLES:
    # Serve logs on port 9000
    amazon-vl /var/log 9000

    # With custom credentials
    AUTH_USER=admin AUTH_HASH='$1$xyz...' amazon-vl /var/log 8080

    # Generate password hash
    openssl passwd -1 -salt "$(openssl rand -base64 6)" "your_password"
`)
}
