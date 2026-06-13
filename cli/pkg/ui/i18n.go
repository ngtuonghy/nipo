package ui

import "nipo/pkg/config"

// Translation defines the localized string resources used across the TUI.
type Translation struct {
	Title                string
	Version              string
	Status               string
	Traffic              string
	Session              string
	SessionStarting      string
	Forwarding           string
	HttpRequests         string
	PressQuit            string
	InitProxies          string
	DownloadBinary       string
	EstablishTunnel      string
	MapSubdomain         string
	Starting             string
	ONLINE               string
	Waiting              string
	Conflict             string
	SubdomainInUse       string
	TunnelPortInfo       string
	UseRandomSubdomain   string
	ExitOption           string
	NavigateSelect       string
	ConfigErrUnsupported string
	ConfigLangUpdated    string
	ConfigHeader         string
	ConfigLangLabel      string
	ConfigBackendLabel   string
	ConfigHint           string
	RootShort            string
	RootLong             string
	RootExample          string
	CmdConfigShort       string
	CmdHttpShort         string
	CmdStartShort        string
	CmdStatusShort       string
	CmdVersionShort      string
	CmdCompletionShort   string
	CmdHelpShort         string
	FlagHelpUsage        string
}

// Locales maps language codes to their respective Translation bundles.
var Locales = map[string]Translation{
	"en": {
		Title:                "NIPO TUNNEL by ngtuonghy",
		Version:              "Version",
		Status:               "Status",
		Traffic:              "Traffic",
		Session:              "Session",
		SessionStarting:      "9h (starts when online)",
		Forwarding:           "Forwarding",
		HttpRequests:         "HTTP Requests",
		PressQuit:            "[Press ctrl+c to quit]",
		InitProxies:          " Initializing local proxy engines...",
		DownloadBinary:       " Downloading tunnel binary...",
		EstablishTunnel:      " Establishing secure tunnel connection...",
		MapSubdomain:         " Mapping subdomains with Nipo API...",
		Starting:             " Starting...",
		ONLINE:               "ONLINE",
		Waiting:              "Waiting...",
		Conflict:             "CONFLICT",
		SubdomainInUse:       "Subdomain '%s' is already in use. Please choose another one.",
		TunnelPortInfo:       "Tunnel: %s (port %d)",
		UseRandomSubdomain:   "Use a random subdomain",
		ExitOption:           "Exit",
		NavigateSelect:       "[↑↓ navigate  Enter select]",
		ConfigErrUnsupported: "Error: unsupported language. Supported: 'en', 'vi'",
		ConfigLangUpdated:    "Language updated to '%s' globally",
		ConfigHeader:         "Nipo Configuration:",
		ConfigLangLabel:      "  Language:    %s",
		ConfigBackendLabel:   "  Backend URL: %s",
		ConfigHint:           "\nTo change interface language, use:\n  nipo config --language [en|vi]",
		RootShort:            "Nipo is a production-ready tunneling platform",
		RootLong:             "Nipo allows you to create secure HTTP tunnels from localhost to the internet.",
		RootExample:          "  # Start an HTTP tunnel on port 3000\n  nipo http 3000\n\n  # Start a tunnel with custom subdomain\n  nipo http 3000 --subdomain myapp\n\n  # Start all tunnels defined in local configuration file (nipo.yml)\n  nipo start\n\n  # Start specific tunnels from the config file\n  nipo start web api\n\n  # View and configure interface language\n  nipo config\n  nipo config --language vi",
		CmdConfigShort:       "Configure Nipo CLI settings",
		CmdHttpShort:         "Start an HTTP tunnel",
		CmdStartShort:        "Start tunnels defined in config",
		CmdStatusShort:       "Check tunnel status",
		CmdVersionShort:      "Print the version number of Nipo",
		CmdCompletionShort:   "Generate the autocompletion script for the specified shell",
		CmdHelpShort:         "Help about any command",
		FlagHelpUsage:        "help for nipo",
	},
	"vi": {
		Title:                "NIPO TUNNEL phát triển bởi ngtuonghy",
		Version:              "Version",
		Status:               "Status",
		Traffic:              "Traffic",
		Session:              "Session",
		SessionStarting:      "9h (bắt đầu khi online)",
		Forwarding:           "Forwarding",
		HttpRequests:         "HTTP Requests",
		PressQuit:            "[Nhấn Ctrl+C để thoát]",
		InitProxies:          " Đang khởi tạo proxy...",
		DownloadBinary:       " Đang tải tunnel binary...",
		EstablishTunnel:      " Đang thiết lập kết nối tunnel...",
		MapSubdomain:         " Đang đăng ký subdomain với Nipo API...",
		Starting:             " Đang khởi động...",
		ONLINE:               "ONLINE",
		Waiting:              "Đang chờ...",
		Conflict:             "CONFLICT",
		SubdomainInUse:       "Subdomain '%s' đã được sử dụng.",
		TunnelPortInfo:       "Tunnel: %s (port %d)",
		UseRandomSubdomain:   "Sử dụng subdomain ngẫu nhiên",
		ExitOption:           "Thoát",
		NavigateSelect:       "[↑↓ di chuyển  Enter chọn]",
		ConfigErrUnsupported: "Lỗi: ngôn ngữ không được hỗ trợ. Hỗ trợ: 'en', 'vi'",
		ConfigLangUpdated:    "Đã cập nhật ngôn ngữ thành '%s'",
		ConfigHeader:         "Nipo Config:",
		ConfigLangLabel:      "  Language:    %s",
		ConfigBackendLabel:   "  Backend URL: %s",
		ConfigHint:           "\nĐể thay đổi ngôn ngữ, sử dụng:\n  nipo config --language [en|vi]",
		RootShort:            "Nipo - nền tảng tunnel chuyên nghiệp",
		RootLong:             "Nipo cho phép bạn tạo HTTP tunnel bảo mật từ localhost ra internet.",
		RootExample:          "  # Khởi động HTTP tunnel trên port 3000\n  nipo http 3000\n\n  # Khởi động tunnel với subdomain tuỳ chọn\n  nipo http 3000 --subdomain myapp\n\n  # Khởi động tất cả tunnel trong config file (nipo.yml)\n  nipo start\n\n  # Khởi động tunnel cụ thể từ config file\n  nipo start web api\n\n  # Xem và thay đổi config\n  nipo config\n  nipo config --language vi",
		CmdConfigShort:       "Cấu hình Nipo CLI",
		CmdHttpShort:         "Khởi động HTTP tunnel",
		CmdStartShort:        "Khởi động tunnel từ config file",
		CmdStatusShort:       "Kiểm tra trạng thái tunnel",
		CmdVersionShort:      "Hiển thị version của Nipo",
		CmdCompletionShort:   "Tạo autocompletion script cho shell",
		CmdHelpShort:         "Hiển thị trợ giúp cho lệnh",
		FlagHelpUsage:        "Trợ giúp cho nipo",
	},
}

// GetT returns the Translation for the currently configured language.
// Falls back to English if the language is unsupported.
func GetT() Translation {
	t, ok := Locales[config.AppConfig.Lang]
	if !ok {
		return Locales["en"]
	}
	return t
}
