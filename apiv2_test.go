package qbt_apiv2

import (
	"fmt"
	
	"testing"
	
)

func TestLogin(t *testing.T) {
	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Printf("%#+v", cli)
}

func TestOpttoStringField(t *testing.T) {
	opt := Optional{
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
	link := `magnet:?xt=urn:btih:16abb2f5bcb405b8ac9d952345f87c87a6af85cc&tr=http://open.acgtracker.com:1096/announce`
	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		panic(err)
	}
	err = cli.AddNewTorrentViaUrl(link, "./", "subject251")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}

}

func TestTorrnetList(t *testing.T) {
	cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		panic(err)
	}

	torrnet, err := cli.TorrentList(Optional{
		"filter": "all",
		// "tag":    "subject251",
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
	torrnet, err := cli.TorrentList(Optional{
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
	cli, err := NewCli("http://localhost:8080")
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
	full, err := cli.GetMainDataFull()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n ====================???========== \n", full)
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
	err = cli.SetAutoDLRule("testing2", AutoDLRule{
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

func TestLsArtMatchRule(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	m, err := cli.LsArtMatchRule("testing")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}
	fmt.Println(m)
}

func TestAddFeeds(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.AddFeed("http://www.kisssub.org/rss-%E4%B8%9C%E4%BA%AC%E7%8C%AB%E7%8C%AB.xml", "")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}
}

func TestDelTorr(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.DelTorrentsFs("79d4e6885d8c796c114ce912b1e612c0a97b01e9", "940c46c2ba144ba90fa95278f8dbc12dd52036c0")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}
}

func TestDelTags(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.DelTags("123", "456")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}
}

func TestRenameFile(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.RenameFile("385191f125783e4dc16689f0ed7b5cf00961155d",
		"[UHA-WINGS][Tengoku Daimakyou][06][x264 1080p][CHS].mp4",
		"es103.mp4")
	if err != nil {
		t.Error(err)
	}
}

func TestSetLocation(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.SetLocation("/upload/anime",
		"385191f125783e4dc16689f0ed7b5cf00961155d")
	if err != nil {
		t.Error(err)
	}
}

func TestAddCategory(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.AddCategory("test", "")
	if err != nil {
		t.Error(err)
	}
}

func TestRmCategoies(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	err = cli.RmCategoies("test")
	if err != nil {
		t.Error(err)
	}
}

func TestFiles(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	fs, err := cli.Files("5ec234d089a09b381e5bc4f7b82241689b5457fd")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fs)
}

func TestPreferences(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	cfg, err := cli.GetPreferences()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(cfg)
}

func TestSetPreferences(t *testing.T) {
	fmt.Println("TestSetPreferences")
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	cfg, err := cli.GetPreferences()
	fmt.Printf("%+v\n", cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cfg.RSSAutoDownloadingEnabled = true
	cfg.RSSProcessingEnabled = true
	cfg.RSSMaxArticlesPerFeed = 50
	cfg.RSSRefreshInterval = 25
	cfg.ProxyType = HttpA
	cfg.ProxyIP = "localhost"
	cfg.ProxyPort = 7890
	err = cli.SetPreferences(cfg)
	// b, err := json.Marshal(cfg)
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	// fmt.Println(string(b))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

func TestGetVersion(t *testing.T) {
	cli, err := NewCli("http://localhost:8991")
	if err != nil {
		panic(err)
	}
	ver, err := cli.GetVersion()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(ver)
}

func TestGetApiVersion(t *testing.T) {
	cli, err := NewCli("http://localhost:10086")
	if err != nil {
		panic(err)
	}
	ver, err := cli.GetApiVersion()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(ver)
}

func TestRenameFolder(t *testing.T) {
	cli, err := NewCli("http://localhost:10086")
	if err != nil {
		panic(err)
	}
	err = cli.RenameFolder("385191f125783e4dc16689f0ed7b5cf00961155d",
		"[UHA-WINGS][Tengoku Daimakyou]",
		"es103")
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveItem(t *testing.T) {
	cli, err := NewCli("http://localhost:10086")
	if err != nil {
		panic(err)
	}
	err = cli.RemoveItem("385191f125783e4dc16689f0ed7b5cf00961155d")
	if err != nil {
		t.Error(err)
	}
}
func TestVersionInt(t *testing.T) {
	fmt.Println(versionInt("v4.3.1"))
	fmt.Println(versionInt("V4.3.1beta"))
	fmt.Println(versionInt("v4.3.1-beta"))
	fmt.Println(versionInt("v4.3.1-rc"))
	fmt.Println(versionInt("v4.3.1-rc0"))
	fmt.Println(versionInt("v4.3.1rc0"))
}