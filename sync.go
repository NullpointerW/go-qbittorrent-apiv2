// Sync
// Sync API implements requests for obtaining changes since the last request. All Sync API methods are under "sync"
// e.g.: /api/v2/sync/{methodName}
package qbt_apiv2

import (
	"encoding/json"
	"io"
)

// Sync holds the sync response struct which contains
// the server state and a map of infohashes to Torrents
type Sync struct {
	Categories map[string]struct {
		Name     string `json:"name"`
		SavePath string `json:"savePath"`
	} `json:"categories"`
	FullUpdate  bool `json:"full_update"`
	Rid         int  `json:"rid"`
	ServerState struct {
		ConnectionStatus  string `json:"connection_status"`
		DhtNodes          int    `json:"dht_nodes"`
		DlInfoData        int    `json:"dl_info_data"`
		DlInfoSpeed       int    `json:"dl_info_speed"`
		DlRateLimit       int    `json:"dl_rate_limit"`
		Queueing          bool   `json:"queueing"`
		RefreshInterval   int    `json:"refresh_interval"`
		UpInfoData        int    `json:"up_info_data"`
		UpInfoSpeed       int    `json:"up_info_speed"`
		UpRateLimit       int    `json:"up_rate_limit"`
		UseAltSpeedLimits bool   `json:"use_alt_speed_limits"`
	} `json:"server_state"`
	Torrents map[string]Torrent `json:"torrents"`
}

func (c *Client) getMainData(rid int) (Sync, error) {
	resp, err := c.postXwwwFormUrlencoded("sync/maindata", optional{
		"rid": rid,
	})
	err = RespOk(resp, err)
	if err != nil {
		return Sync{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Sync{}, err
	}
	s := new(Sync)
	err = json.Unmarshal(b, s)
	if err != nil {
		return Sync{}, err
	}
	return *s, nil
}

func (c *Client) GetMainData() (Sync, error) {
	s, err := c.getMainData(c.rid)
	if err != nil {
		return Sync{}, err
	}
	c.rid = s.Rid
	return s, nil
}
