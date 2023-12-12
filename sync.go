// Sync
// Sync API implements requests for obtaining changes since the last request. All Sync API methods are under "sync"
// e.g.: /api/v2/sync/{methodName}
package qbt_apiv2

import (
	"encoding/json"
	"io"
)

type ServerState struct {
	AllTimeDownload      int64  `json:"alltime_dl,omitempty"`
	AllTimeUpload        int64  `json:"alltime_ul,omitempty"`
	AverageTimeQueue     int64  `json:"average_time_queue,omitempty"`
	ConnectionStatus     string `json:"connection_status,omitempty"`
	DHTNodes             int64  `json:"dht_nodes,omitempty"`
	DLInfoData           int64  `json:"dl_info_data,omitempty"`
	DLInfoSpeed          int64  `json:"dl_info_speed,omitempty"`
	DLRateLimit          int64  `json:"dl_rate_limit,omitempty"`
	FreeSpaceOnDisk      int64  `json:"free_space_on_disk,omitempty"`
	GlobalRatio          string `json:"global_ratio,omitempty"`
	QueuedIOJobs         int64  `json:"queued_io_jobs,omitempty"`
	Queueing             *bool  `json:"queueing,omitempty"`
	ReadCacheHits        string `json:"read_cache_hits,omitempty"`
	ReadCacheOverload    string `json:"read_cache_overload,omitempty"`
	RefreshInterval      int64  `json:"refresh_interval,omitempty"`
	TotalBuffersSize     int64  `json:"total_buffers_size,omitempty"`
	TotalPeerConnections int64  `json:"total_peer_connections,omitempty"`
	TotalQueuedSize      int64  `json:"total_queued_size,omitempty"`
	TotalWastedSession   int64  `json:"total_wasted_session,omitempty"`
	UpInfoData           int64  `json:"up_info_data,omitempty"`
	UpInfoSpeed          int64  `json:"up_info_speed,omitempty"`
	UpRateLimit          int64  `json:"up_rate_limit,omitempty"`
	UseAltSpeedLimits    *bool  `json:"use_alt_speed_limits,omitempty"`
	UseSubcategories     *bool  `json:"use_subcategories,omitempty"`
	WriteCacheOverload   string `json:"write_cache_overload,omitempty"`
}

type Categories struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// Sync holds the sync response struct which contains
// the server state and a map of info hashes to Torrents
type Sync struct {
	Categories        map[string]Categories `json:"categories"`
	CategoriesRemoved []string              `json:"categories_removed"`
	FullUpdate        bool                  `json:"full_update"`
	Rid               int                   `json:"rid"`
	ServerState       ServerState           `json:"server_state"`
	Torrents          map[string]Torrent    `json:"torrents"`
	TorrentsRemoved   []string              `json:"torrents_removed"`
	Tags              []string              `json:"tags"`
	TagsRemoved       []string              `json:"tags_removed"`
	Trackers          map[string][]string   `json:"trackers"`
	TrackersRemoved   []string              `json:"trackers_removed"`
}

type MainData struct {
	ServerState ServerState
	Torrents    map[string]Torrent
	Categories  map[string]Categories
	Tags        []string
	Trackers    map[string][]string
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
	if s.FullUpdate {
		c.mainData = new(MainData)
		c.mainData.ServerState,
			c.mainData.Torrents,
			c.mainData.Categories,
			c.mainData.Trackers,
			c.mainData.Tags = s.ServerState, s.Torrents, s.Categories, s.Trackers, s.Tags
	} else {
		updateMainData(c.mainData, s)
	}
	return s, err
}
func (c *Client) GetMainDataFull() (MainData, error) {
	_, err := c.GetMainData()
	if err != nil {
		return MainData{}, err
	}
	return *c.mainData, nil
}

func updateMainData(m *MainData, s Sync) {
	raw, _ := json.Marshal(s.ServerState)
	_ = json.Unmarshal(raw, &m.ServerState)
	for _, k := range s.TorrentsRemoved {
		delete(m.Torrents, k)
	}
	for _, k := range s.CategoriesRemoved {
		delete(m.Categories, k)
	}
	for _, k := range s.TrackersRemoved {
		delete(m.Trackers, k)
	}
	it := make(map[string]struct{})
	if len(s.TrackersRemoved) > 0 {
		for _, k := range s.TagsRemoved {
			it[k] = struct{}{}
		}
		var nTag []string
		for _, t := range m.Tags {
			if _, e := it[t]; !e {
				nTag = append(nTag, t)
			}
		}
		m.Tags = nTag
	}
	for k, v := range s.Torrents {
		m.Torrents[k] = v
	}

	for k, v := range s.Categories {
		m.Categories[k] = v
	}

	for k, v := range s.Trackers {
		m.Trackers[k] = v
	}
	m.Tags = append(m.Tags, s.Tags...)

}
