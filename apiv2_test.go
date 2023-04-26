package qbt_apiv2

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
		"b":     []byte("tom"),
		"name":  "tom",
		"size":  1.5}
	sm := opt.StringField()

	for k, v := range sm {
		fmt.Println(k + "|" + v)
	}
}

func TestAddTorrnet(t *testing.T) {
	link := `magnet:?xt=urn:btih:7827e38d4b7eac848829fadd8a3c6c28561d0f2c&tr=http%3a%2f%2ft.nyaatracker.com%2fannounce&tr=http%3a%2f%2ftracker.kamigami.org%3a2710%2fannounce&tr=http%3a%2f%2fshare.camoe.cn%3a8080%2fannounce&tr=http%3a%2f%2fopentracker.acgnx.se%2fannounce&tr=http%3a%2f%2fanidex.moe%3a6969%2fannounce&tr=http%3a%2f%2ft.acg.rip%3a6699%2fannounce&tr=https%3a%2f%2ftr.bangumi.moe%3a9696%2fannounce&tr=udp%3a%2f%2ftr.bangumi.moe%3a6969%2fannounce&tr=http%3a%2f%2fopen.acgtracker.com%3a1096%2fannounce&tr=udp%3a%2f%2ftracker.opentrackr.org%3a1337%2fannounce`
	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		panic(err)
	}
	resp, err := cli.AddNewTorrent(link, "./")
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

func TestGetMainData(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 3; i++ {
		sync, err := cli.GetMainData()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v \n ============================== \n", sync)

	}
}

func TestGetTorrnetContent(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	tf, err := cli.GetTorrentContents("7827e38d4b7eac848829fadd8a3c6c28561d0f2c", 0, 1, 2)
	if err != nil {
		fmt.Printf("%+v", err)
		t.FailNow()
	}
	fmt.Printf("%+v", tf)
}

func TestGetAllRssItem(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	item, err := cli.GetAllItems(true)
	if err != nil {
		fmt.Printf("%+v", err)
		t.FailNow()
	}
	fmt.Println(item)
	// i,f:=item.GetWithUrl("http://www.kisssub.org/rss-%E6%94%BE%E5%AD%A6%E5%90%8E%E5%A4%B1%E7%9C%A0%E7%9A%84%E4%BD%A0+%E5%96%B5%E8%90%8C%E5%A5%B6%E8%8C%B6%E5%B1%8B.xml")
	// if f{
	// 	fmt.Println(i)
	// }
}

func TestSetAoDLRule(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.SetAoDLRule("testing2", AutoDLRule{
		Enabled:       false,
		UseRegex:      false,
		AffectedFeeds: []string{"http://www.kisssub.org/rss-%E6%94%BE%E5%AD%A6%E5%90%8E%E5%A4%B1%E7%9C%A0%E7%9A%84%E4%BD%A0+%E5%96%B5%E8%90%8C%E5%A5%B6%E8%8C%B6%E5%B1%8B.xml"},
		SavePath:      "D:\\",
	})
	if err != nil {
		fmt.Printf("%+v", err)
		t.FailNow()
	}
}

func TestLsAoDLRule(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	m, err := cli.LsAutoDLRule()
	if err != nil {
		fmt.Printf("%+v", err)
		t.FailNow()
	}
	fmt.Println(m)
}

func TestLsArtMatchRlue(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	m,err:=cli.LsArtMatchRlue("testing")
	if err != nil {
		fmt.Printf("%+v", err)
		t.FailNow()
	}
	fmt.Println(m)
}
