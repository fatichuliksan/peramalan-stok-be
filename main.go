package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"peramalan-stok-be/src/helper/logger"
	"peramalan-stok-be/src/helper/postgre"
	responseHelper "peramalan-stok-be/src/helper/response"
	validatorHelper "peramalan-stok-be/src/helper/validator"
	viperHelper "peramalan-stok-be/src/helper/viper"
	"strconv"
	"time"

	"peramalan-stok-be/src/delivery/api"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

func main() {
	locationDefault := "Asia/Jakarta"
	config := initConfig()
	if config.GetString(`app.timezone`) != "" {
		os.Setenv("TZ", config.GetString(`app.timezone`))
		locationDefault = config.GetString(`app.timezone`)
	}

	timeLocation, err := time.LoadLocation(locationDefault)
	if err != nil {
		panic(err)
	}

	db := database(config, timeLocation)

	_ = sentry.Init(sentry.ClientOptions{
		Dsn: config.GetString(`glitchtip.dsn`),
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					logger.Default().Println(req)
				}
			}
			logger.Default().Println(event)
			return event
		},
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Debug:            true,
	})

	api := api.NewAPI{
		Echo:      echo.New(),
		Config:    config,
		Validator: validatorHelper.NewValidator(),
		Response:  responseHelper.NewResponse(),
		DB:        db,
	}

	api.Printer = message.NewPrinter(getDefaultLanguageTag(api.Config))

	api.Register()

	port := flag.Int("port", api.Config.GetInt(`app.port`), "port")
	host := flag.String("host", api.Config.GetString(`app.host`), "host")
	flag.Parse()
	api.Echo.Start(*host + ":" + strconv.Itoa(*port))

}

func initConfig() viperHelper.Interface {
	dir := flag.String("dir", "", "config-dir")
	fileType := flag.String("fileType", "json", "fileType")
	loadLanguage(*dir)
	flag.Parse()
	configDir := "config.json"
	if *dir != "" {
		configDir = *dir + "/config.json"
	}
	return viperHelper.NewViper(configDir, *fileType)
}

func loadLanguage(configDir string) {
	type msgs struct {
		Message     string
		Translation interface{}
	}

	languages := [...]string{
		"en", "id",
	}

	var messages []msgs
	for _, lang := range languages {
		tag := language.MustParse(lang)
		path := ""
		if configDir != "" {
			path = configDir + "/static/lang/" + lang + ".json"
		} else {
			path = "static/lang/" + lang + ".json"
		}
		jsonFile, err := os.Open(path)
		if err != nil {
			logger.Default().Println(err)
		} else {
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &messages)
			jsonFile.Close()
			for _, msg := range messages {
				switch m := msg.Translation.(type) {
				case string:
					message.SetString(tag, msg.Message, m)
				case map[string]interface{}:
					message.Set(tag, msg.Message, plural.Selectf(int(m["arg"].(float64)), m["format"].(string), m["case"].([]interface{})...))
				}
			}
		}
	}
}

func getDefaultLanguageTag(config viperHelper.Interface) language.Tag {
	lang := config.GetString("app.language")
	if lang == "" {
		lang = language.English.String()
	}
	tag, _, _ := message.DefaultCatalog.Matcher().Match(language.MustParse(lang))
	return tag
}

func database(config viperHelper.Interface, timeLocation *time.Location) *gorm.DB {
	var appName = config.GetString("app.name")
	logger.Default().Println("App name:", appName)
	logger.Default().Println("Time location:", timeLocation.String())

	logger.Default().Println("Let's try connect to DB")
	db := postgre.NewPostgre(config.GetString("database.postgre.username"),
		config.GetString("database.postgre.password"),
		config.GetString("database.postgre.host"),
		config.GetInt("database.postgre.port"),
		config.GetString("database.postgre.database"),
		appName,
		timeLocation,
	)

	connDB, _ := db.Connect()

	return connDB
}

// func database(config viperHelper.Interface, timeLocation *time.Location) postgre.Database {
// 	var db postgre.Database
// 	var appName = config.GetString("app.name")
// 	logger.Default().Println("App name:", appName)
// 	logger.Default().Println("Time location:", timeLocation.String())

// 	logger.Default().Println("Let's try connect to PGPOOL")
// 	dbOms := postgre.NewPostgre(config.GetString("database.postgre.db_oms.username"),
// 		config.GetString("database.postgre.db_oms.password"),
// 		config.GetString("database.postgre.db_oms.host"),
// 		config.GetInt("database.postgre.db_oms.port"),
// 		config.GetString("database.postgre.db_oms.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbMaster := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_master.username"),
// 		config.GetString("database.postgre.db_master.password"),
// 		config.GetString("database.postgre.db_master.host"),
// 		config.GetInt("database.postgre.db_master.port"),
// 		config.GetString("database.postgre.db_master.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbWms := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_wms.username"),
// 		config.GetString("database.postgre.db_wms.password"),
// 		config.GetString("database.postgre.db_wms.host"),
// 		config.GetInt("database.postgre.db_wms.port"),
// 		config.GetString("database.postgre.db_wms.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbFin := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_fin.username"),
// 		config.GetString("database.postgre.db_fin.password"),
// 		config.GetString("database.postgre.db_fin.host"),
// 		config.GetInt("database.postgre.db_fin.port"),
// 		config.GetString("database.postgre.db_fin.database"),
// 		appName,
// 		timeLocation,
// 	)
// 	dbLog := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_log.username"),
// 		config.GetString("database.postgre.db_log.password"),
// 		config.GetString("database.postgre.db_log.host"),
// 		config.GetInt("database.postgre.db_log.port"),
// 		config.GetString("database.postgre.db_log.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbSfa := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_sfa.username"),
// 		config.GetString("database.postgre.db_sfa.password"),
// 		config.GetString("database.postgre.db_sfa.host"),
// 		config.GetInt("database.postgre.db_sfa.port"),
// 		config.GetString("database.postgre.db_sfa.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbManage := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_manage.username"),
// 		config.GetString("database.postgre.db_manage.password"),
// 		config.GetString("database.postgre.db_manage.host"),
// 		config.GetInt("database.postgre.db_manage.port"),
// 		config.GetString("database.postgre.db_manage.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	connDbOms, _ := dbOms.Connect()
// 	connDbMaster, _ := dbMaster.Connect()
// 	connDbWms, _ := dbWms.Connect()
// 	connDbFin, _ := dbFin.Connect()
// 	connDbLog, _ := dbLog.Connect()
// 	connDbSfa, _ := dbSfa.Connect()
// 	connDbManage, _ := dbManage.Connect()

// 	db.DBOms = connDbOms
// 	db.DBMaster = connDbMaster
// 	db.DBWms = connDbWms
// 	db.DBFin = connDbFin
// 	db.DBLog = connDbLog
// 	db.DBSfa = connDbSfa
// 	db.DBManage = connDbManage

// 	logger.Default().Println("Let's try connect to MASTER")
// 	dbMainOms := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_oms.username"),
// 		config.GetString("database.postgre.db_main_oms.password"),
// 		config.GetString("database.postgre.db_main_oms.host"),
// 		config.GetInt("database.postgre.db_main_oms.port"),
// 		config.GetString("database.postgre.db_main_oms.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbMainMaster := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_master.username"),
// 		config.GetString("database.postgre.db_main_master.password"),
// 		config.GetString("database.postgre.db_main_master.host"),
// 		config.GetInt("database.postgre.db_main_master.port"),
// 		config.GetString("database.postgre.db_main_master.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbMainWms := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_wms.username"),
// 		config.GetString("database.postgre.db_main_wms.password"),
// 		config.GetString("database.postgre.db_main_wms.host"),
// 		config.GetInt("database.postgre.db_main_wms.port"),
// 		config.GetString("database.postgre.db_main_wms.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbMainFin := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_fin.username"),
// 		config.GetString("database.postgre.db_main_fin.password"),
// 		config.GetString("database.postgre.db_main_fin.host"),
// 		config.GetInt("database.postgre.db_main_fin.port"),
// 		config.GetString("database.postgre.db_main_fin.database"),
// 		appName,
// 		timeLocation,
// 	)
// 	dbMainLog := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_log.username"),
// 		config.GetString("database.postgre.db_main_log.password"),
// 		config.GetString("database.postgre.db_main_log.host"),
// 		config.GetInt("database.postgre.db_main_log.port"),
// 		config.GetString("database.postgre.db_main_log.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbMainSfa := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_sfa.username"),
// 		config.GetString("database.postgre.db_main_sfa.password"),
// 		config.GetString("database.postgre.db_main_sfa.host"),
// 		config.GetInt("database.postgre.db_main_sfa.port"),
// 		config.GetString("database.postgre.db_main_sfa.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	dbMainManage := postgre.NewPostgre(
// 		config.GetString("database.postgre.db_main_manage.username"),
// 		config.GetString("database.postgre.db_main_manage.password"),
// 		config.GetString("database.postgre.db_main_manage.host"),
// 		config.GetInt("database.postgre.db_main_manage.port"),
// 		config.GetString("database.postgre.db_main_manage.database"),
// 		appName,
// 		timeLocation,
// 	)

// 	connDbMainOms, _ := dbMainOms.Connect()
// 	connDbMainMaster, _ := dbMainMaster.Connect()
// 	connDbMainWms, _ := dbMainWms.Connect()
// 	connDbMainFin, _ := dbMainFin.Connect()
// 	connDbMainLog, _ := dbMainLog.Connect()
// 	connDbMainSfa, _ := dbMainSfa.Connect()
// 	connDbMainManage, _ := dbMainManage.Connect()

// 	db.DBMainOms = connDbMainOms
// 	db.DBMainMaster = connDbMainMaster
// 	db.DBMainWms = connDbMainWms
// 	db.DBMainFin = connDbMainFin
// 	db.DBMainLog = connDbMainLog
// 	db.DBMainSfa = connDbMainSfa
// 	db.DBMainManage = connDbMainManage
// 	return db
// }
