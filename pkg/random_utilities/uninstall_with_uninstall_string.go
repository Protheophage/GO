// This module is Windows-specific.

package random_utilities

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func UninstallWithUninstallString(appName string) error {
	regPaths := []registry.Key{
		registry.LOCAL_MACHINE,
		registry.LOCAL_MACHINE,
	}
	subKeys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for i, regPath := range regPaths {
		key, err := registry.OpenKey(regPath, subKeys[i], registry.READ)
		if err != nil {
			continue
		}
		defer key.Close()

		subKeyNames, err := key.ReadSubKeyNames(-1)
		if err != nil {
			return err
		}

		for _, subKeyName := range subKeyNames {
			subKey, err := registry.OpenKey(key, subKeyName, registry.READ)
			if err != nil {
				continue
			}
			defer subKey.Close()

			displayName, _, err := subKey.GetStringValue("DisplayName")
			if err != nil || !strings.Contains(displayName, appName) {
				continue
			}

			uninstallString, _, err := subKey.GetStringValue("UninstallString")
			if err != nil {
				return fmt.Errorf("uninstall string not found for %s", appName)
			}

			if strings.Contains(uninstallString, "msiexec") {
				uninstallString = regexp.MustCompile(`msiexec\.exe .*{`).ReplaceAllString(uninstallString, "/Uninstall {")
				uninstallString += " /qn /norestart"
				fmt.Printf("Running: msiexec.exe %s\n", uninstallString)
				cmd := exec.Command("msiexec.exe", strings.Split(uninstallString, " ")...)
				cmd.Run()
			} else {
				fmt.Printf("Running: %s /S\n", uninstallString)
				cmd := exec.Command(uninstallString, "/S")
				cmd.Run()
			}

			fmt.Printf("Silent uninstallation process completed for %s.\n", appName)
			return nil
		}
	}

	return fmt.Errorf("application %s not found", appName)
}
