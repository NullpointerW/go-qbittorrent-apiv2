package qbt_apiv4

import (
	"fmt"
	"io"
	"testing"
)

func TestLogin(t *testing.T) {
	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#+v", cli)
}

func TestOpttoStringField(t *testing.T) {
	opt := optional{
		"count": 3,
		"name":  "tom",
		"size":  1.5}
	sm := opt.StringField()

	for k, v := range sm {
		fmt.Println(k + "|" + v)
	}
}

func TestAddTorrnet(t *testing.T) {
	mglink := `magnet:?xt=urn:btih:F224C0C47C1692008BE4391CE812B239813AD7F1&dn=%E3%80%90%E4%BB%96%E5%92%8C%E5%A5%B9%E7%9A%84%E5%AD%A4%E7%8B%AC%E6%83%85%E4%BA%8B%E3%80%91%E3%80%90%E9%AB%98%E6%B8%85720P%E7%89%88BD-RMVB.%E4%B8%AD%E5%AD%97%E3%80%91%E3%80%902014%E7%BE%8E%E5%9B%BD%E5%89%A7%E6%83%85%E5%A4%A7%E7%89%87%E3%80%91
	  `
	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		panic(err)
	}
	resp, err := cli.AddNewTorrent(mglink, "./download")
	if err != nil {
		panic(err)
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func TestTorrnetList(t *testing.T) {

	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		panic(err)
	}

	torrnet, err := cli.TorrentList(optional{
		"filter": "downloading",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(torrnet[0].Hash)
}

func TestGetTorrentProperties(t *testing.T) {

	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		panic(err)
	}
	torrnet, err := cli.TorrentList(optional{
		"filter": "downloading",
	})
	if err != nil {
		panic(err)
	}
	torrnetProp, err := cli.GetTorrentProperties(torrnet[0].Hash)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%#+v",torrnetProp)
	fmt.Println(torrnetProp.SavePath)
}
