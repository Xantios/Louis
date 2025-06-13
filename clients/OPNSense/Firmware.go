package OPNSense

import (
	"encoding/json"
)

// StatusUpdateResponse this struct is incomplete by design!
// The unmarshalling a bumload of packages is not relevant; however, it does break reliability
type StatusUpdateResponse struct {
	ApiVersion          string `json:"api_version"`
	Connection          string `json:"connection"`
	DownloadSize        string `json:"download_size"`
	LastCheck           string `json:"last_check"`
	NeedsReboot         string `json:"needs_reboot"`
	OsVersion           string `json:"os_version"`
	ProductId           string `json:"product_id"`
	ProductTarget       string `json:"product_target"`
	ProductVersion      string `json:"product_version"`
	ProductAbi          string `json:"product_abi"`
	Repository          string `json:"repository"`
	UpgradeMajorMessage string `json:"upgrade_major_message"`
	UpgradeMajorVersion string `json:"upgrade_major_version"`
	UpgradeNeedsReboot  string `json:"upgrade_needs_reboot"`
	Product             struct {
		ProductAbi   string `json:"product_abi"`
		ProductArch  string `json:"product_arch"`
		ProductCheck struct {
			ApiVersion          string `json:"api_version"`
			Connection          string `json:"connection"`
			DownloadSize        string `json:"download_size"`
			LastCheck           string `json:"last_check"`
			NeedsReboot         string `json:"needs_reboot"`
			OsVersion           string `json:"os_version"`
			ProductId           string `json:"product_id"`
			ProductTarget       string `json:"product_target"`
			ProductVersion      string `json:"product_version"`
			ProductAbi          string `json:"product_abi"`
			Repository          string `json:"repository"`
			UpgradeMajorMessage string `json:"upgrade_major_message"`
			UpgradeMajorVersion string `json:"upgrade_major_version"`
			UpgradeNeedsReboot  string `json:"upgrade_needs_reboot"`
		} `json:"product_check"`
		ProductCopyrightOwner string `json:"product_copyright_owner"`
		ProductCopyrightUrl   string `json:"product_copyright_url"`
		ProductCopyrightYears string `json:"product_copyright_years"`
		ProductEmail          string `json:"product_email"`
		ProductHash           string `json:"product_hash"`
		ProductId             string `json:"product_id"`
		ProductLatest         string `json:"product_latest"`
		ProductLog            int    `json:"product_log"`
		ProductMirror         string `json:"product_mirror"`
		ProductName           string `json:"product_name"`
		ProductNickname       string `json:"product_nickname"`
		ProductRepos          string `json:"product_repos"`
		ProductSeries         string `json:"product_series"`
		ProductTier           string `json:"product_tier"`
		ProductTime           string `json:"product_time"`
		ProductVersion        string `json:"product_version"`
		ProductWebsite        string `json:"product_website"`
	} `json:"product"`
	StatusMsg    string `json:"status_msg"`
	StatusReboot string `json:"status_reboot"`
	Status       string `json:"status"`
}

func (c *Client) Update() (bool, string, error) {

	// First, send a POST request to trigger the update
	poke, body, err := c.Post("/core/firmware/status")
	if err != nil {
		return false, "", err
	}

	if poke.StatusCode != 200 {
		return false, "", err
	}

	// Second, send a GET request to get the status
	resp, body, err := c.Get("/core/firmware/status")
	if err != nil {
		return false, "", err
	}

	if resp.StatusCode != 200 {
		return false, "", err
	}

	// Unmarshal the response into a StatusUpdateResponse struct
	jsb := &StatusUpdateResponse{}
	err = json.Unmarshal([]byte(body), jsb)
	if err != nil {
		return false, "", err
	}

	if c.debug {
		jsb2, _ := json.Marshal(jsb)
		PrintJson(string(jsb2))
	}

	updateAvailable := jsb.Status != "none"
	return updateAvailable, jsb.StatusMsg, nil
}
