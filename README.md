
go-qbittorrent-apiv2
-----------------------
### Description
---------
[qBittorrent web API(v2.8.3)](https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)) wrapper for go

Requires qBittorrent version ≥ v4.1.



Some of the code is reused from this repository [superturkey650/go-qbittorrent](https://github.com/superturkey650/go-qbittorrent).

#####TODO
At present, the API for [RSS](https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#rss-experimental) has been fully wrapped, as well as partially wrapping the API for torrent and sync. Welcome to submit PR to wrap all APIs.




###Installation
---------
 go get in your module.
```
$ go get github.com/NullpointerW/go-qbittorrent-apiv2
```

###Usage
---------
``` go
    import (
        qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
    )
    // When connecting to qBittorrent on the local network and 
    // 'Bypass from localhost' setting is active.
    // The parameters after 'host' can be ignored.
    // e.g.:NewCli("http://localhost:8991")
    cli, err := NewCli("http://localhost:8991", "admin", "123456")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
    torrnet, err := cli.TorrentList(optional{
		"filter": "downloading",
	})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	torrnetProp, err := cli.GetTorrentProperties(torrnet[0].Hash)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println(torrnetProp.SavePath)
    m, err := cli.LsAutoDLRule()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println(m)

```
