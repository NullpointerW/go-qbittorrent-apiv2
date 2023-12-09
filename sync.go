// Sync
// Sync API implements requests for obtaining changes since the last request. All Sync API methods are under "sync"
// e.g.: /api/v2/sync/{methodName}
package qbt_apiv2

import (
	"encoding/json"
	"io"
)

type ServerState struct {
	AllTimeDownload      int64  `json:"alltime_dl"`
	AllTimeUpload        int64  `json:"alltime_ul"`
	AverageTimeQueue     int64  `json:"average_time_queue"`
	ConnectionStatus     string `json:"connection_status"`
	DHTNodes             int64  `json:"dht_nodes"`
	DLInfoData           int64  `json:"dl_info_data"`
	DLInfoSpeed          int64  `json:"dl_info_speed"`
	DLRateLimit          int64  `json:"dl_rate_limit"`
	FreeSpaceOnDisk      int64  `json:"free_space_on_disk"`
	GlobalRatio          string `json:"global_ratio"`
	QueuedIOJobs         int64  `json:"queued_io_jobs"`
	Queueing             *bool  `json:"queueing"`
	ReadCacheHits        string `json:"read_cache_hits"`
	ReadCacheOverload    string `json:"read_cache_overload"`
	RefreshInterval      int64  `json:"refresh_interval"`
	TotalBuffersSize     int64  `json:"total_buffers_size"`
	TotalPeerConnections int64  `json:"total_peer_connections"`
	TotalQueuedSize      int64  `json:"total_queued_size"`
	TotalWastedSession   int64  `json:"total_wasted_session"`
	UpInfoData           int64  `json:"up_info_data"`
	UpInfoSpeed          int64  `json:"up_info_speed"`
	UpRateLimit          int64  `json:"up_rate_limit"`
	UseAltSpeedLimits    *bool  `json:"use_alt_speed_limits"`
	UseSubcategories     *bool  `json:"use_subcategories"`
	WriteCacheOverload   string `json:"write_cache_overload"`
}

// Sync holds the sync response struct which contains
// the server state and a map of info hashes to Torrents
type Sync struct {
	Categories map[string]struct {
		Name     string `json:"name"`
		SavePath string `json:"savePath"`
	} `json:"categories"`
	CategoriesRemoved []string            `json:"categories_removed"`
	FullUpdate        bool                `json:"full_update"`
	Rid               int                 `json:"rid"`
	ServerState       ServerState         `json:"server_state"`
	Torrents          map[string]Torrent  `json:"torrents"`
	TorrentsRemoved   []string            `json:"torrents_removed"`
	Tags              []string            `json:"tags"`
	TagsRemoved       []string            `json:"tags_removed"`
	Trackers          map[string][]string `json:"trackers"`
}

func (c *Client) getMainData(rid int) (Sync, error) {
	resp, err := c.postXwwwFormUrlencoded("sync/maindata", Optional{
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
