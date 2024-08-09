// Torrent management
// All Torrent management API methods are under "torrents",
// e.g.: /api/v2/torrents/methodName.
package qbt_apiv2

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

// Torrent holds a basic torrent object from qbittorrent
// which is `sync/maindata` ,`torrents/info` returned
type Torrent struct {
	AddedOn           int     `json:"added_on"`
	AmountLeft        int     `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      float64 `json:"availability"`
	Category          string  `json:"category"`
	Completed         int     `json:"completed"`
	CompletionOn      int     `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DLLimit           int     `json:"dl_limit"`
	Dlspeed           int     `json:"dlspeed"`
	DownloadPath      string  `json:"download_path"`
	Downloaded        int     `json:"downloaded"`
	DownloadedSession int     `json:"downloaded_session"`
	Eta               int     `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	InfohashV1        string  `json:"infohash_v1"`
	InfohashV2        string  `json:"infohash_v2"`
	LastActivity      int     `json:"last_activity"`
	MagnetURI         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int     `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumIncomplete     int     `json:"num_incomplete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float64 `json:"progress"`
	Ratio             float64 `json:"ratio"`
	RatioLimit        float64 `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeedingTime       int     `json:"seeding_time"`
	SeedingTimeLimit  int     `json:"seeding_time_limit"`
	SeenComplete      int     `json:"seen_complete"`
	SeqDL             bool    `json:"seq_dl"`
	Size              int     `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int     `json:"time_active"`
	TotalSize         int     `json:"total_size"`
	Tracker           string  `json:"tracker"`
	TrackersCount     int     `json:"trackers_count"`
	UpLimit           int     `json:"up_limit"`
	Uploaded          int     `json:"uploaded"`
	UploadedSession   int     `json:"uploaded_session"`
	Upspeed           int     `json:"upspeed"`
}

// TorrentProp holds a torrent object from qbittorrent
// with more information than BasicTorrent
type TorrentProp struct {
	AdditionDate           int     `json:"addition_date"`
	Comment                string  `json:"comment"`
	CompletionDate         int     `json:"completion_date"`
	CreatedBy              string  `json:"created_by"`
	CreationDate           int     `json:"creation_date"`
	DlLimit                int     `json:"dl_limit"`
	DlSpeed                int     `json:"dl_speed"`
	DlSpeedAvg             int     `json:"dl_speed_avg"`
	Eta                    int     `json:"eta"`
	LastSeen               int     `json:"last_seen"`
	NbConnections          int     `json:"nb_connections"`
	NbConnectionsLimit     int     `json:"nb_connections_limit"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PieceSize              int     `json:"piece_size"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Reannounce             int     `json:"reannounce"`
	SavePath               string  `json:"save_path"`
	SeedingTime            int     `json:"seeding_time"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	ShareRatio             float64 `json:"share_ratio"`
	TimeElapsed            int     `json:"time_elapsed"`
	TotalDownloaded        int     `json:"total_downloaded"`
	TotalDownloadedSession int     `json:"total_downloaded_session"`
	TotalSize              int     `json:"total_size"`
	TotalUploaded          int     `json:"total_uploaded"`
	TotalUploadedSession   int     `json:"total_uploaded_session"`
	TotalWasted            int     `json:"total_wasted"`
	UpLimit                int     `json:"up_limit"`
	UpSpeed                int     `json:"up_speed"`
	UpSpeedAvg             int     `json:"up_speed_avg"`
}

// Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg      string `json:"msg"`
	NumPeers int    `json:"num_peers"`
	Status   string `json:"status"`
	URL      string `json:"url"`
}

// WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

// TorrentFile holds a torrent file object from qbittorrent
type TorrentFile struct {
	IsSeed       bool    `json:"is_seed"`
	Name         string  `json:"name"`
	Priority     int     `json:"priority"`
	Progress     float64 `json:"progress"`
	Size         int     `json:"size"`
	PieceRange   []int   `json:"piece_range"`
	Availability float64 `json:"availability"`
}

type File struct {
	Availability float64 `json:"availability"`
	Index        int     `json:"index"`
	Name         string  `json:"name"`
	PieceRange   []int   `json:"piece_range"`
	Priority     int     `json:"priority"`
	Progress     float64 `json:"progress"`
	Size         int     `json:"size"`
	IsSeed       bool    `json:"is_seed"`
}

func (c *Client) AddNewTorrent(opt Optional) error {
	resp, err := c.postMultipartData("torrents/add", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	if err = RespBodyOk(resp.Body, ErrAddTorrnetfailed); err != nil {
		return err
	}
	return nil
}

func (c *Client) AddNewTorrentViaUrl(url, path string, tags ...string) error {
	ap, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("cannot conv abs_path: %s: %w ", path, err)
	}
	opt := Optional{
		"urls":     url,
		"savepath": ap,
	}
	if len(tags) > 0 {
		var ts string
		for _, t := range tags {
			ts += t + ","
		}
		ts = ts[:len(ts)-1]
		opt["tags"] = ts
	}
	err = c.AddNewTorrent(opt)
	return err
}

func (c *Client) TorrentList(opt Optional) ([]Torrent, error) {
	resp, err := c.postXwwwFormUrlencoded("torrents/info", opt)

	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bt := new([]Torrent)
	err = json.Unmarshal(b, bt)
	if err != nil {
		return nil, err
	}
	return *bt, nil
}

func (c *Client) GetTorrentProperties(hash string) (TorrentProp, error) {
	resp, err := c.postXwwwFormUrlencoded("torrents/properties", Optional{
		"hash": hash,
	})
	err = RespOk(resp, err)
	if err != nil {
		return TorrentProp{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return TorrentProp{}, err
	}
	t := new(TorrentProp)
	err = json.Unmarshal(b, t)
	if err != nil {
		return TorrentProp{}, err
	}
	return *t, nil
}

func (c *Client) GetTorrentContents(hash string, indexes ...int) ([]TorrentFile, error) {
	opt := Optional{
		"hash": hash,
	}
	if len(indexes) > 0 {
		var idxes string
		for _, idx := range indexes {
			idxes += strconv.Itoa(idx) + "|"
		}
		idxes = idxes[:len(idxes)-1]
		opt["indexes"] = idxes
	}

	resp, err := c.postXwwwFormUrlencoded("torrents/files", opt)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tf := new([]TorrentFile)
	err = json.Unmarshal(b, tf)
	if err != nil {
		return nil, err
	}
	return *tf, nil
}

func (c *Client) DelTorrents(delfile bool, hashes ...string) error {
	hs := strings.Join(hashes, "|")
	opt := Optional{
		"hashes":      hs,
		"deleteFiles": delfile,
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/delete", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) DelTorrentsFs(hashes ...string) error {
	return c.DelTorrents(true, hashes...)
}

func (c *Client) DelTags(tags ...string) error {
	ts := strings.Join(tags, ",")
	opt := Optional{
		"tags": ts,
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/deleteTags", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) RenameFile(hash, old, new string) error {
	opt := Optional{
		"hash":    hash,
		"oldPath": old,
		"newPath": new,
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/renameFile", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) SetLocation(location string, hashes ...string) error {
	hs := strings.Join(hashes, "|")
	opt := Optional{
		"hashes":   hs,
		"location": location,
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/setLocation", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) RenameFolder(hash, old, new string) error {
	opt := Optional{
		"hash":    hash,
		"oldPath": old,
		"newPath": new,
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/renameFolder", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) AddCategory(categoryName, savePath string) error {
	opt := Optional{
		"category": categoryName,
	}
	if savePath != "" {
		opt["savePath"] = savePath
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/createCategory", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) RmCategoies(categories ...string) error {
	categs := strings.Join(categories, "\n")
	opt := Optional{
		"categories": categs,
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/removeCategories", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) Files(hash string, indexs ...string) ([]File, error) {
	idxs := ""
	if len(indexs) != 0 {
		idxs = strings.Join(indexs, "|")
	}

	opt := Optional{
		"hash": hash,
	}
	if idxs != "" {
		opt["indexes"] = idxs
	}
	resp, err := c.postXwwwFormUrlencoded("torrents/files", opt)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var fs []File
	err = json.Unmarshal(b, &fs)
	if err != nil {
		return nil, err
	}
	return fs, nil
}
