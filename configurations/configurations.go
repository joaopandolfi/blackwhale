package configurations

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/unrolled/secure"
)

type SessionConfiguration struct {
	Name    string
	Store   *sessions.CookieStore
	Options *sessions.Options
}

type WhiteListAuthRoutes struct {
	Paths map[string]bool
}

type FirewallSettings struct {
	LocalHost  string
	RemoteHost string
}

type Timeout struct {
	Write time.Duration
	Read  time.Duration
}

type Opsec struct {
	Options       secure.Options
	Debug         bool
	TLSCert       string
	TLSKey        string
	BCryptCost    int // 10,11,12,13,14
	JWTSecret     string
	TokenValidity int
	DefaultPwd    string
}

type Configurations struct {
	Name        string
	MysqlUrl    string
	MongoUrl    string
	MongoDb     string
	MongoPool   int
	Port        string
	CRONThreads int
	CORS        string
	Timeout     Timeout

	SlackToken   string
	SlackWebHook []string
	SlackChannel string

	BCryptSecret string

	Session  SessionConfiguration
	Security Opsec

	WhiteListAuthRoutes WhiteListAuthRoutes
	Templates           map[string]*pongo2.Template

	StaticPath      string
	StaticDir       string
	StaticPagesDir  string
	ResetHash       string
	UploadPath      string
	MaxSizeMbUpload int64

	FirewallSettings FirewallSettings
}

var Configuration Configurations

// LoadJsonFile - Load file from Json config
func LoadJsonFile(fpath string) map[string]string {
	var confFile map[string]string
	file, err := os.Open(fpath)
	if err != nil {
		panic("Config file is not present")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&confFile)
	if err != nil {
		panic("Config file is not parseable")
	}
	return confFile
}

// LoadConfig from previous configurations
func LoadConfig(c Configurations) {
	Configuration = c

	Configuration.Session.Store.Options = Configuration.Session.Options

	if Configuration.Security.BCryptCost == 0 {
		Configuration.Security.BCryptCost = 14
	}
	if Configuration.MongoPool == 0 {
		Configuration.MongoPool = 5
	}

	// Run in max of cpus
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// LoadFromFile configurations
func LoadFromFile(path string) Configurations {
	fconf := LoadJsonFile(path)

	tkVal, _ := strconv.Atoi(fconf["TOKEN_VALIDITY_MINUTES"])
	mongoPool, _ := strconv.Atoi(fconf["MONGO_POOL"])
	timeout, _ := strconv.Atoi(fconf["SERVER_TIMEOUT"])
	bcryptCost, _ := strconv.Atoi(fconf["BCRYPT_COST"])

	return Configurations{
		Name: fconf["SERVER_NAME"],

		MysqlUrl: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			fconf["MYSQL_USER"],
			fconf["MYSQL_PASSWORD"],
			fconf["MYSQL_HOST"],
			fconf["MYSQL_PORT"],
			fconf["MYSQL_DB"]),

		MongoUrl:  fconf["MONGO_URL"],
		MongoDb:   fconf["MONGO_DB"],
		MongoPool: mongoPool,

		Port: ":" + fconf["SERVER_PORT"],
		CORS: fconf["SERVER_CORS"],

		Timeout: Timeout{
			Write: time.Duration(timeout) * time.Second,
			Read:  time.Duration(timeout) * time.Second,
		},

		MaxSizeMbUpload: 10 << 55, // min << max

		BCryptSecret: fconf["BCRYPT_SECRET"],
		ResetHash:    fconf["RESET_HASH"],

		// Session
		Session: SessionConfiguration{
			Name:  fconf["SESSION_NAME"],
			Store: sessions.NewCookieStore([]byte(fconf["SESSION_STORE"])),
			Options: &sessions.Options{
				Path:     "/",
				MaxAge:   3600 * 2, //86400 * 7,
				HttpOnly: true,
			},
		},

		Security: Opsec{
			Options: secure.Options{
				BrowserXssFilter:   true,
				ContentTypeNosniff: false, // Da pau nos js
				SSLHost:            "locahost:443",
				SSLRedirect:        false,
			},
			BCryptCost:    bcryptCost,
			Debug:         fconf["SERVER_DEBUG"] == "true",
			TLSCert:       fconf["TLS_CERT"],
			TLSKey:        fconf["TLS_KEY"],
			JWTSecret:     fconf["JWT_SECRET"],
			TokenValidity: tkVal,
			DefaultPwd:    fconf["SERVER_DEFAULT_PASSWORD"],
		},

		Templates: make(map[string]*pongo2.Template),

		// Slack
		SlackToken:   fconf["SLACK_TOKEN"],
		SlackWebHook: []string{fconf["SLACK_WEBHOOK_1"], fconf["SLACK_WEBHOOK_2"]},
		SlackChannel: fconf["SLACK_CHANNEL"],

		// Firewall]
		FirewallSettings: FirewallSettings{
			LocalHost:  "localhost:8080",
			RemoteHost: "localhosy:443",
		},
	}
}

// Load default configurations
func Load() {

	Configuration = Configurations{

		Name: "Blackwale - GO",
		MysqlUrl: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			"root",         // User
			"rootpassword", // password
			"localhost",    // host
			"3311",         // port
			"blackwhale"),  // Database

		MongoUrl:  "mongodb://127.0.0.1:27017",
		MongoDb:   "blackwhale",
		MongoPool: 5,

		CRONThreads: 20,
		Port:        ":8990",
		CORS:        "*",

		Timeout: Timeout{
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
		Session: SessionConfiguration{
			Name:  "A2%!#23g4$0$",
			Store: sessions.NewCookieStore([]byte("_-)(AS(&HSDH@ˆ@@#$##$*{{{$$}}}(U$$#@D)&#Y!)P(@M)(Xyeg3b321k5*443@@##@$!")),
			Options: &sessions.Options{
				Path:     "/",
				MaxAge:   3600 * 2, //86400 * 7,
				HttpOnly: true,
			},
		},

		Security: Opsec{
			Options: secure.Options{
				BrowserXssFilter:   true,
				ContentTypeNosniff: false, // Da pau nos js
				SSLHost:            "locahost:443",
				SSLRedirect:        false,
			},
			BCryptCost:    14,
			Debug:         true,
			TLSCert:       "",
			TLSKey:        "",
			JWTSecret:     "",
			TokenValidity: 60,
		},

		Templates: make(map[string]*pongo2.Template),

		// Slack
		SlackToken:   "",
		SlackWebHook: []string{"", ""},
		SlackChannel: "",

		// Firewall]
		FirewallSettings: FirewallSettings{
			LocalHost:  "localhost:8080",
			RemoteHost: "localhosy:443",
		},
	}

	Configuration.Session.Store.Options = Configuration.Session.Options

	// Run in max of cpus
	runtime.GOMAXPROCS(runtime.NumCPU())
}
