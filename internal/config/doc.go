// Package config provides types and helpers for loading portwatch
// configuration from a JSON file.
//
// # File format
//
// Configuration is stored as a JSON object with the following fields:
//
//	{
//	  "host":     "localhost",   // host to scan (default: "localhost")
//	  "ports":    [22, 80, 443], // ports to monitor (required)
//	  "interval": "5s"           // scan interval (default: "5s")
//	}
//
// # Usage
//
//	cfg, err := config.Load("portwatch.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
package config
