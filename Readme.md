# Velopack Go
[![Go Reference](https://pkg.go.dev/badge/github.com/quaadgras/velopack-go.svg)](https://pkg.go.dev/github.com/quaadgras/velopack-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/quaadgras/velopack-go)](https://goreportcard.com/report/github.com/quaadgras/velopack-go)

It's easy to distribute Go desktop applications with automatic updates!

1. Add this module to your project:

$ `go get github.com/quaadgras/velopack-go`

2. Add automatic updates to your project:

Here's as a one-liner, that will download updates in the background and
apply them the next time the application starts up:
```go
go velopack.DownloadUpdatesInTheBackground("https://the.place/you-will-host/updates")
```

Alternatively, you can also develop your own update function:
```go
package main

import "github.com/quaadgras/velopack-go/velopack"

func init() {
  velopack.Run(velopack.App{
  	AutoApplyOnStartup: true,
  })
}

func update() error {
	manager, err := velopack.NewUpdateManager("https://the.place/you-will-host/updates")
	if err != nil {
		return err
	}
	latest, status, err := manager.CheckForUpdates();
	if err != nil {
		return err
	}
	if status == velopack.UpdateAvailable {
		if err := manager.DownloadUpdates(latest, func(progress uint){
			// show progress to the user
		}); err != nil {
			return err
		}
		if err := manager.ApplyUpdatesAndRestart(latest); err != nil {
			return err
		}
	}
	return nil
}
```

3. Follow the official Velopack guide to package your application for distribution.

https://docs.velopack.io/packaging/overview

**Quick Start**

`vpk pack --packId "MyCompany.MyApp" --packVersion "0.0.0" --packDir path/to/build --mainExe executable.name -o /local/place/for/updates/and/releases`
