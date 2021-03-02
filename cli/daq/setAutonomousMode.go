package daq

import (
	"fmt"

	"cli.rc.ccm.dunescience.org/internal"
)

// SetAutonomousMode enables or disables autonomous mode on given daq-applications
func SetAutonomousMode(w internal.Writers, c *internal.RCConfig, enabled bool, names ...string) error {
	actionStr := "enabling"
	if !enabled {
		actionStr = "disabling"
	}
	for _, name := range names {
		fmt.Fprintf(w.Out, "%s autonomous mode for daq application %s... ", actionStr, name)
		app, err := internal.GetResource(c, internal.DAQAppKind, name)
		if err != nil {
			return err
		}
		if app.Spec["enabled"].(bool) == enabled {
			fmt.Fprintln(w.Out, "UNCHANGED")
		}

		app.Spec["enabled"] = enabled
		internal.Apply(w, c, *app)
		fmt.Fprintln(w.Out, "OK")
	}
	return nil
}
