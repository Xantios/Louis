package OPNSense

import (
	"net/http"
	"os"
)

// Backup pulls a backup from OPNSense
func (c *Client) Backup() error {
	// Eccentric URL found in a script on Codeberg by SweetGood
	// https://codeberg.org/SWEETGOOD/andersgood-opnsense-scripts/src/branch/main/backup-opnsense-via-api.sh
	resp, body, err := c.Get("/core/backup/download/this")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	// Check if a backup directory exists
	if _, err := os.Stat(c.backupPath); os.IsNotExist(err) {
		err = os.MkdirAll(c.backupPath, 0755)
		if err != nil {
			return err
		}
	}

	// Store backup
	// @TODO: Add date to filename
	fName := "backup.xml"
	err = os.WriteFile(c.backupPath+"/"+fName, []byte(body), 0644)
	if err != nil {
		return err
	}

	return nil
}
