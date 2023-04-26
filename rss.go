package qbt_apiv2

import (
	"encoding/json"
	"io"
)

// map type for `rss/items` responed json schema
type RssItem map[string]Item

// get rss item via rss url
// if the specified URL does not exist in these items, the returned bool value is false
// otherwise it is true
func (m RssItem) GetWithUrl(url string) (Item, bool) {
	for _, it := range m {
		if it.Url == url {
			return it, true
		}
	}
	return Item{}, false
}

// RSS Schema
type Item struct {
	Articles      []Article `json:"articles"`
	HasError      bool      `json:"hasError"`
	IsLoading     bool      `json:"isLoading"`
	LastBuildDate string    `json:"lastBuildDate"`
	Title         string    `json:"title"`
	Uid           string    `json:"uid"`
	Url           string    `json:"url"`
}

type Article struct {
	Author      string `json:"author"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Id          string `json:"id"`
	Link        string `json:"link"`
	Title       string `json:"title"`
	TorrentURL  string `json:"torrentURL"`
	IsRead      bool   `json:"isRead,omitempty"`
}

type AutoDLRule struct {
	Enabled                   bool     `json:"enabled"`
	MustContain               string   `json:"mustContain"`
	MustNotContain            string   `json:"mustNotContain"`
	UseRegex                  bool     `json:"useRegex"`
	EpisodeFilter             string   `json:"episodeFilter"`
	SmartFilter               bool     `json:"smartFilter"`
	PreviouslyMatchedEpisodes []string `json:"previouslyMatchedEpisodes"`
	AffectedFeeds             []string `json:"affectedFeeds"`
	IgnoreDays                int      `json:"ignoreDays"`
	LastMatch                 string   `json:"lastMatch"`
	AddPaused                 bool     `json:"addPaused"`
	AssignedCategory          string   `json:"assignedCategory"`
	SavePath                  string   `json:"savePath"`
}

// RSS All RSS API methods are under "rss", e.g.: /api/v2/rss/methodName.
func (c *Client) AddFolder(path string) error {
	resp, err := c.postXwwwFormUrlencoded("rss/addFolder", optional{
		"path": path,
	})
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) AddFeed(url, path string) error {
	opt := optional{
		"url": url,
	}
	if path != "" {
		opt["path"] = path
	}
	resp, err := c.postXwwwFormUrlencoded("rss/addFeed", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) RemoveItem(path string) error {
	resp, err := c.postXwwwFormUrlencoded("rss/removeItem", optional{
		"path": path,
	})
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) MoveItem(dst, src string) error {
	resp, err := c.postXwwwFormUrlencoded("rss/moveItem", optional{
		"itemPath": src,
		"destPath": dst,
	})
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) GetAllItems(withData bool) (RssItem, error) {
	opt := optional{}
	if withData {
		opt["withData"] = true
	}
	resp, err := c.postXwwwFormUrlencoded("rss/items", opt)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ri := new(RssItem)
	json.Unmarshal(b, ri)
	return *ri, nil
}

func (c *Client) MarkAsRead(itemPath, articleId string) error {
	opt := optional{
		"itemPath": itemPath,
	}
	if articleId != "" {
		opt["articleId"] = articleId
	}
	resp, err := c.postXwwwFormUrlencoded("rss/markAsRead", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) RefreshItem(itemPath string) error {
	opt := optional{
		"itemPath": itemPath,
	}
	resp, err := c.postXwwwFormUrlencoded("rss/refreshItem", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) SetAoDLRule(ruleName string, ruleDef AutoDLRule) error {
	b, _ := json.Marshal(ruleDef)
	opt := optional{
		"ruleName": ruleName,
		"ruleDef":  string(b),
	}
	resp, err := c.postXwwwFormUrlencoded("rss/setRule", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}
