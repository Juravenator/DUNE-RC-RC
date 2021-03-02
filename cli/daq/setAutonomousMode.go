package daq

import (
	"fmt"
	"io"

	"cli.rc.ccm.dunescience.org/internal"
)

// SetAutonomousMode enables or disables autonomous mode on given daq-applications
func SetAutonomousMode(writer io.Writer, c *internal.RCConfig, enabled bool, names ...string) error {
	actionStr := "enabling"
	if !enabled {
		actionStr = "disabling"
	}
	for _, name := range names {
		fmt.Fprintf(writer, "%s autonomous mode for daq application %s... ", actionStr, name)
		app, err := internal.GetResource(c, internal.DAQAppKind, name)
		if err != nil {
			return err
		}
		if app.Spec["enabled"].(bool) == enabled {
			fmt.Fprintln(writer, "UNCHANGED")
		}

		app.Spec["enabled"] = enabled
		internal.Apply(writer, c, *app)
		fmt.Fprintln(writer, "OK")
	}
	return nil
}
