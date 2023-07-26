package qbt_apiv2

import (
	"encoding/json"
	"io"
)

type proxyTyp int

const (
	Http    proxyTyp = iota + 1 // HTTP proxy without authentication
	Socks5                      // SOCKS5 proxy without authentication
	HttpA                       // HTTP proxy with authentication
	Socks5A                     // SOCKS5 proxy with authentication
	Socks4                      // SOCKS4 proxy without authentication
)

type Config struct {
	AddTrackers                        string         `json:"add_trackers"`
	AddTrackersEnabled                 bool           `json:"add_trackers_enabled"`
	AltDLLimit                         int            `json:"alt_dl_limit"`
	AltUpLimit                         int            `json:"alt_up_limit"`
	AlternativeWebuiEnabled            bool           `json:"alternative_webui_enabled"`
	AlternativeWebuiPath               string         `json:"alternative_webui_path"`
	AnnounceIP                         string         `json:"announce_ip"`
	AnnounceToAllTiers                 bool           `json:"announce_to_all_tiers"`
	AnnounceToAllTrackers              bool           `json:"announce_to_all_trackers"`
	AnonymousMode                      bool           `json:"anonymous_mode"`
	AsyncIoThreads                     int            `json:"async_io_threads"`
	AutoDeleteMode                     int            `json:"auto_delete_mode"`
	AutoTmmEnabled                     bool           `json:"auto_tmm_enabled"`
	AutorunEnabled                     bool           `json:"autorun_enabled"`
	AutorunOnTorrentAddedEnabled       bool           `json:"autorun_on_torrent_added_enabled"`
	AutorunOnTorrentAddedProgram       string         `json:"autorun_on_torrent_added_program"`
	AutorunProgram                     string         `json:"autorun_program"`
	BannedIPS                          string         `json:"banned_IPs"`
	BittorrentProtocol                 int            `json:"bittorrent_protocol"`
	BlockPeersOnPrivilegedPorts        bool           `json:"block_peers_on_privileged_ports"`
	BypassAuthSubnetWhitelist          string         `json:"bypass_auth_subnet_whitelist"`
	BypassAuthSubnetWhitelistEnabled   bool           `json:"bypass_auth_subnet_whitelist_enabled"`
	BypassLocalAuth                    bool           `json:"bypass_local_auth"`
	CategoryChangedTmmEnabled          bool           `json:"category_changed_tmm_enabled"`
	CheckingMemoryUse                  int            `json:"checking_memory_use"`
	ConnectionSpeed                    int            `json:"connection_speed"`
	CurrentInterfaceAddress            string         `json:"current_interface_address"`
	CurrentNetworkInterface            string         `json:"current_network_interface"`
	Dht                                bool           `json:"dht"`
	DiskCache                          int            `json:"disk_cache"`
	DiskCacheTTL                       int            `json:"disk_cache_ttl"`
	DiskIoReadMode                     int            `json:"disk_io_read_mode"`
	DiskIoType                         int            `json:"disk_io_type"`
	DiskIoWriteMode                    int            `json:"disk_io_write_mode"`
	DiskQueueSize                      int            `json:"disk_queue_size"`
	DLLimit                            int            `json:"dl_limit"`
	DontCountSlowTorrents              bool           `json:"dont_count_slow_torrents"`
	DyndnsDomain                       string         `json:"dyndns_domain"`
	DyndnsEnabled                      bool           `json:"dyndns_enabled"`
	DyndnsPassword                     string         `json:"dyndns_password"`
	DyndnsService                      int            `json:"dyndns_service"`
	DyndnsUsername                     string         `json:"dyndns_username"`
	EmbeddedTrackerPort                int            `json:"embedded_tracker_port"`
	EmbeddedTrackerPortForwarding      bool           `json:"embedded_tracker_port_forwarding"`
	EnableCoalesceReadWrite            bool           `json:"enable_coalesce_read_write"`
	EnableEmbeddedTracker              bool           `json:"enable_embedded_tracker"`
	EnableMultiConnectionsFromSameIP   bool           `json:"enable_multi_connections_from_same_ip"`
	EnablePieceExtentAffinity          bool           `json:"enable_piece_extent_affinity"`
	EnableUploadSuggestions            bool           `json:"enable_upload_suggestions"`
	Encryption                         int            `json:"encryption"`
	ExcludedFileNames                  string         `json:"excluded_file_names"`
	ExcludedFileNamesEnabled           bool           `json:"excluded_file_names_enabled"`
	ExportDir                          string         `json:"export_dir"`
	ExportDirFin                       string         `json:"export_dir_fin"`
	FilePoolSize                       int            `json:"file_pool_size"`
	HashingThreads                     int            `json:"hashing_threads"`
	IdnSupportEnabled                  bool           `json:"idn_support_enabled"`
	IncompleteFilesEXT                 bool           `json:"incomplete_files_ext"`
	IPFilterEnabled                    bool           `json:"ip_filter_enabled"`
	IPFilterPath                       string         `json:"ip_filter_path"`
	IPFilterTrackers                   bool           `json:"ip_filter_trackers"`
	LimitLANPeers                      bool           `json:"limit_lan_peers"`
	LimitTCPOverhead                   bool           `json:"limit_tcp_overhead"`
	LimitUTPRate                       bool           `json:"limit_utp_rate"`
	ListenPort                         int            `json:"listen_port"`
	Locale                             string         `json:"locale"`
	Lsd                                bool           `json:"lsd"`
	MailNotificationAuthEnabled        bool           `json:"mail_notification_auth_enabled"`
	MailNotificationEmail              string         `json:"mail_notification_email"`
	MailNotificationEnabled            bool           `json:"mail_notification_enabled"`
	MailNotificationPassword           string         `json:"mail_notification_password"`
	MailNotificationSender             string         `json:"mail_notification_sender"`
	MailNotificationSMTP               string         `json:"mail_notification_smtp"`
	MailNotificationSSLEnabled         bool           `json:"mail_notification_ssl_enabled"`
	MailNotificationUsername           string         `json:"mail_notification_username"`
	MaxActiveCheckingTorrents          int            `json:"max_active_checking_torrents"`
	MaxActiveDownloads                 int            `json:"max_active_downloads"`
	MaxActiveTorrents                  int            `json:"max_active_torrents"`
	MaxActiveUploads                   int            `json:"max_active_uploads"`
	MaxConcurrentHTTPAnnounces         int            `json:"max_concurrent_http_announces"`
	MaxConnec                          int            `json:"max_connec"`
	MaxConnecPerTorrent                int            `json:"max_connec_per_torrent"`
	MaxRatio                           int            `json:"max_ratio"`
	MaxRatioAct                        int            `json:"max_ratio_act"`
	MaxRatioEnabled                    bool           `json:"max_ratio_enabled"`
	MaxSeedingTime                     int            `json:"max_seeding_time"`
	MaxSeedingTimeEnabled              bool           `json:"max_seeding_time_enabled"`
	MaxUploads                         int            `json:"max_uploads"`
	MaxUploadsPerTorrent               int            `json:"max_uploads_per_torrent"`
	MemoryWorkingSetLimit              int            `json:"memory_working_set_limit"`
	OutgoingPortsMax                   int            `json:"outgoing_ports_max"`
	OutgoingPortsMin                   int            `json:"outgoing_ports_min"`
	PeerTos                            int            `json:"peer_tos"`
	PeerTurnover                       int            `json:"peer_turnover"`
	PeerTurnoverCutoff                 int            `json:"peer_turnover_cutoff"`
	PeerTurnoverInterval               int            `json:"peer_turnover_interval"`
	PerformanceWarning                 bool           `json:"performance_warning"`
	Pex                                bool           `json:"pex"`
	PreallocateAll                     bool           `json:"preallocate_all"`
	ProxyAuthEnabled                   bool           `json:"proxy_auth_enabled"`
	ProxyHostnameLookup                bool           `json:"proxy_hostname_lookup"`
	ProxyIP                            string         `json:"proxy_ip"`
	ProxyPassword                      string         `json:"proxy_password"`
	ProxyPeerConnections               bool           `json:"proxy_peer_connections"`
	ProxyPort                          int            `json:"proxy_port"`
	ProxyTorrentsOnly                  bool           `json:"proxy_torrents_only"`
	ProxyType                          proxyTyp       `json:"proxy_type"`
	ProxyUsername                      string         `json:"proxy_username"`
	QueueingEnabled                    bool           `json:"queueing_enabled"`
	RandomPort                         bool           `json:"random_port"`
	ReannounceWhenAddressChanged       bool           `json:"reannounce_when_address_changed"`
	RecheckCompletedTorrents           bool           `json:"recheck_completed_torrents"`
	RefreshInterval                    int            `json:"refresh_interval"`
	RequestQueueSize                   int            `json:"request_queue_size"`
	ResolvePeerCountries               bool           `json:"resolve_peer_countries"`
	ResumeDataStorageType              string         `json:"resume_data_storage_type"`
	RSSAutoDownloadingEnabled          bool           `json:"rss_auto_downloading_enabled"`
	RSSDownloadRepackProperEpisodes    bool           `json:"rss_download_repack_proper_episodes"`
	RSSMaxArticlesPerFeed              int            `json:"rss_max_articles_per_feed"`
	RSSProcessingEnabled               bool           `json:"rss_processing_enabled"`
	RSSRefreshInterval                 int            `json:"rss_refresh_interval"`
	RSSSmartEpisodeFilters             string         `json:"rss_smart_episode_filters"`
	SavePath                           string         `json:"save_path"`
	SavePathChangedTmmEnabled          bool           `json:"save_path_changed_tmm_enabled"`
	SaveResumeDataInterval             int            `json:"save_resume_data_interval"`
	ScanDirs                           map[string]int `json:"scan_dirs"`
	ScheduleFromHour                   int            `json:"schedule_from_hour"`
	ScheduleFromMin                    int            `json:"schedule_from_min"`
	ScheduleToHour                     int            `json:"schedule_to_hour"`
	ScheduleToMin                      int            `json:"schedule_to_min"`
	SchedulerDays                      int            `json:"scheduler_days"`
	SchedulerEnabled                   bool           `json:"scheduler_enabled"`
	SendBufferLowWatermark             int            `json:"send_buffer_low_watermark"`
	SendBufferWatermark                int            `json:"send_buffer_watermark"`
	SendBufferWatermarkFactor          int            `json:"send_buffer_watermark_factor"`
	SlowTorrentDLRateThreshold         int            `json:"slow_torrent_dl_rate_threshold"`
	SlowTorrentInactiveTimer           int            `json:"slow_torrent_inactive_timer"`
	SlowTorrentULRateThreshold         int            `json:"slow_torrent_ul_rate_threshold"`
	SocketBacklogSize                  int            `json:"socket_backlog_size"`
	SsrfMitigation                     bool           `json:"ssrf_mitigation"`
	StartPausedEnabled                 bool           `json:"start_paused_enabled"`
	StopTrackerTimeout                 int            `json:"stop_tracker_timeout"`
	TempPath                           string         `json:"temp_path"`
	TempPathEnabled                    bool           `json:"temp_path_enabled"`
	TorrentChangedTmmEnabled           bool           `json:"torrent_changed_tmm_enabled"`
	TorrentContentLayout               string         `json:"torrent_content_layout"`
	TorrentStopCondition               string         `json:"torrent_stop_condition"`
	UpLimit                            int            `json:"up_limit"`
	UploadChokingAlgorithm             int            `json:"upload_choking_algorithm"`
	UploadSlotsBehavior                int            `json:"upload_slots_behavior"`
	Upnp                               bool           `json:"upnp"`
	UpnpLeaseDuration                  int            `json:"upnp_lease_duration"`
	UseCategoryPathsInManualMode       bool           `json:"use_category_paths_in_manual_mode"`
	UseHTTPS                           bool           `json:"use_https"`
	UTPTCPMixedMode                    int            `json:"utp_tcp_mixed_mode"`
	ValidateHTTPSTrackerCertificate    bool           `json:"validate_https_tracker_certificate"`
	WebUIAddress                       string         `json:"web_ui_address"`
	WebUIBanDuration                   int            `json:"web_ui_ban_duration"`
	WebUIClickjackingProtectionEnabled bool           `json:"web_ui_clickjacking_protection_enabled"`
	WebUICSRFProtectionEnabled         bool           `json:"web_ui_csrf_protection_enabled"`
	WebUICustomHTTPHeaders             string         `json:"web_ui_custom_http_headers"`
	WebUIDomainList                    string         `json:"web_ui_domain_list"`
	WebUIHostHeaderValidationEnabled   bool           `json:"web_ui_host_header_validation_enabled"`
	WebUIHTTPSCERTPath                 string         `json:"web_ui_https_cert_path"`
	WebUIHTTPSKeyPath                  string         `json:"web_ui_https_key_path"`
	WebUIMaxAuthFailCount              int            `json:"web_ui_max_auth_fail_count"`
	WebUIPort                          int            `json:"web_ui_port"`
	WebUIReverseProxiesList            string         `json:"web_ui_reverse_proxies_list"`
	WebUIReverseProxyEnabled           bool           `json:"web_ui_reverse_proxy_enabled"`
	WebUISecureCookieEnabled           bool           `json:"web_ui_secure_cookie_enabled"`
	WebUISessionTimeout                int            `json:"web_ui_session_timeout"`
	WebUIUpnp                          bool           `json:"web_ui_upnp"`
	WebUIUseCustomHTTPHeadersEnabled   bool           `json:"web_ui_use_custom_http_headers_enabled"`
	WebUIUsername                      string         `json:"web_ui_username"`
}

func (c *Client) GetPreferences() (cfg Config, err error) {
	resp, err := c.postXwwwFormUrlencoded("app/preferences", nil)
	err = RespOk(resp, err)
	if err != nil {
		return cfg, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(b, &cfg)
	return cfg, err
}
func (c *Client) SetPreferences(cfg Config) (err error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	opt := Optional{
		"json": string(b),
	}
	resp, err := c.postXwwwFormUrlencoded("app/setPreferences", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) GetVersion() (ver string, err error) {
	resp, err := c.postXwwwFormUrlencoded("app/version", nil)
	err = RespOk(resp, err)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) GetApiVersion() (ver string, err error) {
	resp, err := c.postXwwwFormUrlencoded("app/webapiVersion", nil)
	err = RespOk(resp, err)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
