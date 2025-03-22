package viper

import (
	"fmt"
	"peramalan-stok-be/src/helper/logger"
	"strings"

	"github.com/spf13/viper"
)

// Interface ...
type Interface interface {
	GetDirectory() string
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetStringSlice(key string) []string
	GetStringMapString(key string) map[string]string

	DbLinkMaster() string
	DbLinkOms() string
	DbLinkWms() string
	DbLinkFin() string
	DbLinkLog() string
	DbLinkSfa() string
	DbLinkManage() string
}

type viperHelper struct {
	dir      string
	fileType string
}

// NewViper ...
func NewViper(dir string, fileType string) Interface {
	v := &viperHelper{
		dir:      dir,
		fileType: fileType,
	}
	v.Init()
	return v
}

// Init ...
func (t *viperHelper) Init() {
	viper.SetEnvPrefix(`test`)
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(`.`, `_`)
	viper.SetEnvKeyReplacer(replacer)
	viper.AddConfigPath(".")
	viper.SetConfigType(t.fileType)
	viper.SetConfigFile(t.dir)
	logger.Default().Println("Config: " + viper.ConfigFileUsed())
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// GetString ...
func (t *viperHelper) GetDirectory() string {
	return t.dir
}

// GetString ...
func (t *viperHelper) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt ...
func (t *viperHelper) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool ...
func (t *viperHelper) GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetBool ...
func (t *viperHelper) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (v *viperHelper) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (v *viperHelper) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// // prod/demo
func (v *viperHelper) DbLinkMaster() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_master.host"), v.GetString("database.postgre.db_master.database"), v.GetString("database.postgre.db_oms.username"), v.GetString("database.postgre.db_oms.password"))
}

func (v *viperHelper) DbLinkOms() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_oms.host"), v.GetString("database.postgre.db_oms.database"), v.GetString("database.postgre.db_oms.username"), v.GetString("database.postgre.db_oms.password"))
}

func (v *viperHelper) DbLinkWms() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_wms.host"), v.GetString("database.postgre.db_wms.database"), v.GetString("database.postgre.db_wms.username"), v.GetString("database.postgre.db_wms.password"))
}

func (v *viperHelper) DbLinkFin() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_fin.host"), v.GetString("database.postgre.db_fin.database"), v.GetString("database.postgre.db_fin.username"), v.GetString("database.postgre.db_fin.password"))
}

func (v *viperHelper) DbLinkLog() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_log.host"), v.GetString("database.postgre.db_log.database"), v.GetString("database.postgre.db_log.username"), v.GetString("database.postgre.db_log.password"))
}

func (v *viperHelper) DbLinkSfa() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_sfa.host"), v.GetString("database.postgre.db_sfa.database"), v.GetString("database.postgre.db_sfa.username"), v.GetString("database.postgre.db_sfa.password"))
}

func (v *viperHelper) DbLinkManage() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s", v.GetString("database.postgre.db_manage.host"), v.GetString("database.postgre.db_manage.database"), v.GetString("database.postgre.db_manage.username"), v.GetString("database.postgre.db_manage.password"))
}
