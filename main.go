package main

import (
        "fmt"
        "os"
        "os/user"

        "github.com/Sirupsen/logrus"
)

func init() {
        logrus.SetLevel(logrus.DebugLevel)
        logrus.SetOutput(os.Stderr)
        // Retrieve current user to get home directory
        usr, err := user.Current()
        if err != nil {
                fatalf("cannot get current user: %v", err)
        }

        // Get home directory for current user
        homeDir := usr.HomeDir
        if homeDir == "" {
                fatalf("cannot get current user home directory")
        }
}

func main() {
        fmt.Printf("Main")
}

func fatalf(format string, args ...interface{}) {
        fmt.Printf("* fatal: "+format+"\n", args...)
        os.Exit(1)
}
