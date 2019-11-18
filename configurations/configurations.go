package configurations

import (
	"fmt"
	"runtime"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/unrolled/secure"
)

type sessionConfiguration struct {
	Name    string
	Store   *sessions.CookieStore
	Options *sessions.Options
}

type whiteListAuthRoutes struct {
	Paths map[string]bool
}

type firewallSettings struct {
	LocalHost  string
	RemoteHost string
}

type timeout struct {
	Write time.Duration
	Read  time.Duration
}

type opsec struct {
	Options secure.Options
}

type Configurations struct {
	MysqlUrl    string
	MongoUrl    string
	MongoDb     string
	Port        string
	CRONThreads int
	CORS        string
	Timeout     timeout

	SlackToken   string
	SlackWebHook []string
	SlackChannel string

	BCryptSecret string

	Session  sessionConfiguration
	Security opsec

	WhiteListAuthRoutes whiteListAuthRoutes
	Templates           map[string]*pongo2.Template

	StaticPath      string
	StaticDir       string
	StaticPagesDir  string
	ResetHash       string
	UploadPath      string
	MaxSizeMbUpload int64

	FirewallSettings firewallSettings
}

var Configuration Configurations

func LoadConfig(c Configurations) {
	Configuration = c

	Configuration.Session.Store.Options = Configuration.Session.Options

	// Run in max of cpus
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func Load() {

	Configuration = Configurations{

		MysqlUrl: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			"root",         // User
			"rootpassword", // password
			"localhost",    // host
			"3311",         // port
			"blackwhale"),  // Database

		MongoUrl: "",
		MongoDb:  "",

		CRONThreads: 20,
		Port:        ":8990",
		CORS:        "*",

		Timeout: timeout{
			Write: 60 * time.Second,
			Read:  60 * time.Second,
		},

		ResetHash: "R3S3tM$g!c0",

		StaticPath:     "/static/",
		StaticDir:      "./views/public/",
		StaticPagesDir: "./views/pages/",
		UploadPath:     "./views/upload/",

		MaxSizeMbUpload: 10 << 55, // min << max

		BCryptSecret: "#1$eY)&E&0",

		// Session
		Session: sessionConfiguration{
			Name:  "A2%!#23g4$0$",
			Store: sessions.NewCookieStore([]byte("_-)(AS(&HSDH@ˆ@@#$##$*{{{$$}}}(U$$#@D)&#Y!)P(@M)(Xyeg3b321k5*443@@##@$!")),
			Options: &sessions.Options{
				Path:     "/",
				MaxAge:   3600 * 2, //86400 * 7,
				HttpOnly: true,
			},
		},

		Security: opsec{
			Options: secure.Options{
				BrowserXssFilter:   true,
				ContentTypeNosniff: false, // Da pau nos js
				SSLHost:            "locahost:443",
				SSLRedirect:        false,
			},
		},

		Templates: make(map[string]*pongo2.Template),

		// Slack
		SlackToken:   "",
		SlackWebHook: []string{"", ""},
		SlackChannel: "",

		// Firewall]
		FirewallSettings: firewallSettings{
			LocalHost:  "localhost:8080",
			RemoteHost: "localhosy:443",
		},
	}

	Configuration.Session.Store.Options = Configuration.Session.Options

	// Run in max of cpus
	runtime.GOMAXPROCS(runtime.NumCPU())
}
