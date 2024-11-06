package appconfig

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func GetConfig() *viper.Viper {
	if Config == nil {
		log.Fatal("Config is not initialized")
	}
	return Config
}

// system config
var (
	SystemAreaUid int64
	Enable_TLS    bool = false
	Listen_addr   string
	Base_URL      string
	Service_path  string

	IWK_set_uri string
)

//Config of SQL database

var (
	Enable_sql   bool = false
	SQL_Host     string
	SQL_Port     string
	SQL_sslmode  string
	SQL_dbname   string
	SQL_Username string
	SQL_Password string
	SQL_Schema   string
)

// MQTT broke config
var (
	MQTT_Server_Public_URL     string
	MQTT_Server_Host           string
	MQTT_Server_Port           int64
	MQTT_Server_Usessl         bool = false
	MQTT_Server_ClientID       string
	MQTT_Server_Username       string
	MQTT_Server_Password       string
	MQTT_Server_Subcribe_Topic string
)

// Common API Clients
var (
	Client_Cabin_Syns_URI    string
	Client_IAM_Login_URI     string
	Client_IAM_Client_ID     string
	Client_IAM_Client_Secret string
)

// kafka broke config
var (
	Kafka_Enable                                          bool = false
	Kafka_Bootstrap_Ivi                                   string
	Kafka_Topic_Producer_Device_Command_Ivi               string
	Kafka_Topic_Producer_Device_Command_Ivi_Face_Register string
	Kafka_Topic_Consumer_Device_Feedback_Ivi              string
	Kafka_Topic_Consumer_Device_Info_Ivi                  string
	Kafka_Topic_Consumer_AI_VMS_Event                     string
	Kafka_Topic_Consumer_Group                            string
)

//Storage

var (
	Storage_Public_Download     string
	Storage_Private_Download    string
	Storage_Public_GetLibraries string
	Vector_Register_Blacklist   string
	Vector_Search_Blacklist     string
	Vector_Delete_Blacklist     string
)

// MinIO config
var (
	Minio_Endpoint          string
	Minio_Accesskey_ID      string
	Minio_Secret_Access_Key string
	Minio_UseSSL            bool = false
)

// func AppConfig() {
// 	// Setup Configs
// 	viper := viper.New()

// 	// Set Default Configs
// 	// Important: Viper configuration keys are case insensitive.
// 	viper.SetDefault("system_area_uid", 222001)
// 	viper.SetDefault("enable_tls", false)
// 	viper.SetDefault("listen_addr", "0.0.0.0:8080")
// 	viper.SetDefault("base_url", "127.0.0.1:8080")
// 	viper.SetDefault("service_path", "/ivis/vms/api/v0")

// 	viper.SetDefault("jwk_set_uri", "https://sbs.basesystem.one/ivis/iam/api/v0/.well-known/jwks.json")

// 	//viper.SetDefault("cert_file", "/certs/cert.crt")
// 	//viper.SetDefault("key_file", "/certs/key.key")
// 	//viper.SetDefault("config_watch_interval", 15)
// 	//viper.SetDefault("grpc_socket_path", "/grpc.sock")
// 	//viper.SetDefault("sql_type", "postgres")
// 	viper.SetDefault("enable_sql", true)
// 	viper.SetDefault("sql_host", "4.194.17.112")
// 	viper.SetDefault("sql_port", 5432)
// 	viper.SetDefault("sql_sslmode", "disable")
// 	viper.SetDefault("sql_dbname", "vms")
// 	viper.SetDefault("sql_user", "ivissbs")
// 	viper.SetDefault("sql_password", "hugQRT8k5HLt4b4a2kHq7fRql9i30E51ggbM")
// 	viper.SetDefault("sql_schema", "ivis_vms")

// 	// MQTT broke config
// 	viper.SetDefault("mqtt_server_public_url", "tcp://sbs.basesystem.one:1883")
// 	viper.SetDefault("mqtt_server_host", "sbs.basesystem.one")
// 	viper.SetDefault("mqtt_server_port", 1883)
// 	viper.SetDefault("mqtt_server_usessl", false)
// 	viper.SetDefault("mqtt_inflight", 65535)
// 	viper.SetDefault("mqtt_completetimeout", 60000)
// 	viper.SetDefault("mqtt_qos", 1)
// 	viper.SetDefault("mqtt_service_client_id", "8b835b46-c3a2-485a-a1af-ad42aa1b9a4d")
// 	viper.SetDefault("mqtt_service_username", "8b835b46-c3a2-485a-a1af-ad42aa1b9a4d")
// 	viper.SetDefault("mqtt_service_password", "5eb55ca8-e315-4945-9977-461359885e7a")
// 	viper.SetDefault("mqtt_service_subscribe_topic", "IVIS/VMS/8b835b46-c3a2-485a-a1af-ad42aa1b9a4d/signaling")

// 	// Common api clients
// 	viper.SetDefault("client_cabin_sync_uri", "https://sbs.basesystem.one/ivis/api/v0/cabin")
// 	viper.SetDefault("client_iam_login_uri", "https://sbs.basesystem.one/ivis/iam/api/v0/client/login")
// 	viper.SetDefault("client_iam_clientid", "e7b56179-ac9d-419b-80ed-a18c41285aab")
// 	viper.SetDefault("client_iam_clientsecret", "3b9fed7e-525a-472d-bf56-1aca6a4519ab")

// 	// Kafka broke config
// 	viper.SetDefault("kafka_enable", true)
// 	// viper.SetDefault("kafka_bootstrap_ivis", "4.194.17.112:9092")
// 	// viper.SetDefault("kafka_topic_producer_device_command_ivis", "DEV_SMARTNVR_COMMAND_D1")
// 	// viper.SetDefault("kafka_topic_consumer_device_feedback_ivis", "DEV_SMARTNVR_FEEDBACK_D1")
// 	viper.SetDefault("kafka_bootstrap_ivi", "4.193.155.125:9092")
// 	viper.SetDefault("kafka_topic_producer_device_command_ivi", "dmp_command")
// 	viper.SetDefault("kafka_topic_producer_device_command_ivi_face_register", "DEV_FACE_REGISTER")
// 	viper.SetDefault("kafka_topic_consumer_device_feedback_ivi", "dmp_device_NVR")
// 	viper.SetDefault("kafka_topic_consumer_device_info_ivi", "dmp_device_info")
// 	viper.SetDefault("kafka_topic_consumer_ai_vms_event", "EVENT_AI_CORE_RESPONSE")
// 	// viper.SetDefault("kafka_producer_id", "dev.camnetID")
// 	// viper.SetDefault("kafka_consumer_group", "dev.camnet.vms.test.wsnode")
// 	viper.SetDefault("kafka_consumer_group", "Khoatest_3")

// 	// Storage
// 	viper.SetDefault("storage_public_dowload", "https://sbs.basesystem.one/ivis/storage/api/v0/libraries/public/download/")
// 	viper.SetDefault("storage_private_dowload", "https://sbs.basesystem.one/ivis/storage/api/v0/libraries/download/")
// 	viper.SetDefault("storage_public_getlibraries", "https://sbs.basesystem.one/ivis/storage/api/v0/libraries/")
// 	viper.SetDefault("vector_register_blacklist", "https://facereg.basesystem.one/facereg/insert")
// 	viper.SetDefault("vector_delete_blacklist", "https://facereg.basesystem.one/facereg/remove")
// 	viper.SetDefault("vector_search_blacklist", "https://facereg.basesystem.one/member/search")

// 	viper.SetDefault("minio_endpoint", "dev-minio-api.basesystem.one")
// 	viper.SetDefault("minio_accesskey_id", "")      // export APP_MINIO_ACCESSKEY_ID=
// 	viper.SetDefault("minio_secret_access_key", "") // export APP_MINIO_SECRET_ACCESS_KEY=
// 	viper.SetDefault("minio_usessl", true)

// 	// Assign values to global variables
// 	SystemAreaUid = viper.GetInt64("system_area_uid")
// 	Enable_TLS = viper.GetBool("enable_tls")
// 	Listen_addr = viper.GetString("listen_addr")
// 	Base_URL = viper.GetString("base_url")
// 	Service_path = viper.GetString("service_path")
// 	IWK_set_uri = viper.GetString("jwk_set_uri")

// 	Enable_sql = viper.GetBool("enable_sql")
// 	SQL_Host = viper.GetString("sql_host")
// 	SQL_Port = viper.GetInt64("sql_port")
// 	SQL_sslmode = viper.GetString("sql_sslmode")
// 	SQL_dbname = viper.GetString("sql_dbname")
// 	SQL_Username = viper.GetString("sql_user")
// 	SQL_Password = viper.GetString("sql_password")
// 	SQL_Schema = viper.GetString("sql_schema")

// 	MQTT_Server_Public_URL = viper.GetString("mqtt_server_public_url")
// 	MQTT_Server_Host = viper.GetString("mqtt_server_host")
// 	MQTT_Server_Port = viper.GetInt64("mqtt_server_port")
// 	MQTT_Server_Usessl = viper.GetBool("mqtt_server_usessl")
// 	MQTT_Server_ClientID = viper.GetString("mqtt_service_client_id")
// 	MQTT_Server_Username = viper.GetString("mqtt_service_username")
// 	MQTT_Server_Password = viper.GetString("mqtt_service_password")
// 	MQTT_Server_Subcribe_Topic = viper.GetString("mqtt_service_subscribe_topic")

// 	Client_Cabin_Syns_URI = viper.GetString("client_cabin_sync_uri")
// 	Client_IAM_Login_URI = viper.GetString("client_iam_login_uri")
// 	Client_IAM_Client_ID = viper.GetString("client_iam_clientid")
// 	Client_IAM_Client_Secret = viper.GetString("client_iam_clientsecret")

// 	Kafka_Enable = viper.GetBool("kafka_enable")
// 	Kafka_Bootstrap_Ivi = viper.GetString("")
// 	Kafka_Topic_Producer_Device_Command_Ivi = viper.GetString("kafka_bootstrap_ivi")
// 	Kafka_Topic_Producer_Device_Command_Ivi_Face_Register = viper.GetString("kafka_topic_producer_device_command_ivi_face_register")
// 	Kafka_Topic_Consumer_Device_Feedback_Ivi = viper.GetString("kafka_topic_consumer_device_feedback_ivi")
// 	Kafka_Topic_Consumer_Device_Info_Ivi = viper.GetString("kafka_topic_consumer_device_info_ivi")
// 	Kafka_Topic_Consumer_AI_VMS_Event = viper.GetString("kafka_topic_consumer_ai_vms_event")
// 	Kafka_Topic_Consumer_Group = viper.GetString("kafka_consumer_group")

// 	Storage_Public_Download = viper.GetString("storage_public_download")
// 	Storage_Private_Download = viper.GetString("storage_private_download")
// 	Storage_Public_GetLibraries = viper.GetString("storage_public_getlibraries")
// 	Vector_Register_Blacklist = viper.GetString("vector_register_blacklist")
// 	Vector_Delete_Blacklist = viper.GetString("vector_delete_blacklist")
// 	Vector_Search_Blacklist = viper.GetString("vector_search_blacklist")

// 	viper.SetDefault("minio_endpoint", "dev-minio-api.basesystem.one")
// 	viper.SetDefault("minio_accesskey_id", "PNFIK0TCWXFZQKU0")                      // export APP_MINIO_ACCESSKEY_ID=
// 	viper.SetDefault("minio_secret_access_key", "SU54IHJRCR3SLH4C1GXMPAZWVJFPJOPP") // export

// }

func GetConfigEnv() {
	// Load environment variables from .env file
	err := godotenv.Load("/home/camnet/Desktop/vms/conf/pkg.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set configuration defaults with Viper and load from environment variables
	viper.SetDefault("system_area_uid", os.Getenv("SYSTEM_AREA_UID"))
	viper.SetDefault("enable_tls", os.Getenv("ENABLE_TLS") == "true")
	viper.SetDefault("listen_addr", os.Getenv("LISTEN_ADDR"))
	viper.SetDefault("base_url", os.Getenv("BASE_URL"))
	viper.SetDefault("service_path", os.Getenv("SERVICE_PATH"))

	viper.SetDefault("jwk_set_uri", os.Getenv("JWK_SET_URI"))

	viper.SetDefault("enable_sql", os.Getenv("ENABLE_SQL") == "true")
	viper.SetDefault("sql_host", os.Getenv("SQL_HOST"))
	viper.SetDefault("sql_port", os.Getenv("SQL_PORT"))
	viper.SetDefault("sql_sslmode", os.Getenv("SQL_SSLMODE"))
	viper.SetDefault("sql_dbname", os.Getenv("SQL_DBNAME"))
	viper.SetDefault("sql_user", os.Getenv("SQL_USER"))
	viper.SetDefault("sql_password", os.Getenv("SQL_PASSWORD"))
	viper.SetDefault("sql_schema", os.Getenv("SQL_SCHEMA"))

	viper.SetDefault("mqtt_server_public_url", os.Getenv("MQTT_SERVER_PUBLIC_URL"))
	viper.SetDefault("mqtt_server_host", os.Getenv("MQTT_SERVER_HOST"))
	viper.SetDefault("mqtt_server_port", os.Getenv("MQTT_SERVER_PORT"))
	viper.SetDefault("mqtt_server_usessl", os.Getenv("MQTT_SERVER_USESSL") == "true")
	viper.SetDefault("mqtt_inflight", os.Getenv("MQTT_INFLIGHT"))
	viper.SetDefault("mqtt_completetimeout", os.Getenv("MQTT_COMPLETETIMEOUT"))
	viper.SetDefault("mqtt_qos", os.Getenv("MQTT_QOS"))
	viper.SetDefault("mqtt_service_client_id", os.Getenv("MQTT_SERVICE_CLIENT_ID"))
	viper.SetDefault("mqtt_service_username", os.Getenv("MQTT_SERVICE_USERNAME"))
	viper.SetDefault("mqtt_service_password", os.Getenv("MQTT_SERVICE_PASSWORD"))
	viper.SetDefault("mqtt_service_subscribe_topic", os.Getenv("MQTT_SERVICE_SUBSCRIBE_TOPIC"))

	viper.SetDefault("client_cabin_sync_uri", os.Getenv("CLIENT_CABIN_SYNC_URI"))
	viper.SetDefault("client_iam_login_uri", os.Getenv("CLIENT_IAM_LOGIN_URI"))
	viper.SetDefault("client_iam_clientid", os.Getenv("CLIENT_IAM_CLIENTID"))
	viper.SetDefault("client_iam_clientsecret", os.Getenv("CLIENT_IAM_CLIENTSECRET"))

	viper.SetDefault("kafka_enable", os.Getenv("KAFKA_ENABLE") == "true")
	viper.SetDefault("kafka_bootstrap_ivi", os.Getenv("KAFKA_BOOTSTRAP_IVI"))
	viper.SetDefault("kafka_topic_producer_device_command_ivi", os.Getenv("KAFKA_TOPIC_PRODUCER_DEVICE_COMMAND_IVI"))
	viper.SetDefault("kafka_topic_producer_device_command_ivi_face_register", os.Getenv("KAFKA_TOPIC_PRODUCER_DEVICE_COMMAND_IVI_FACE_REGISTER"))
	viper.SetDefault("kafka_topic_consumer_device_feedback_ivi", os.Getenv("KAFKA_TOPIC_CONSUMER_DEVICE_FEEDBACK_IVI"))
	viper.SetDefault("kafka_topic_consumer_device_info_ivi", os.Getenv("KAFKA_TOPIC_CONSUMER_DEVICE_INFO_IVI"))
	viper.SetDefault("kafka_topic_consumer_ai_vms_event", os.Getenv("KAFKA_TOPIC_CONSUMER_AI_VMS_EVENT"))
	viper.SetDefault("kafka_consumer_group", os.Getenv("KAFKA_CONSUMER_GROUP"))

	// Storage
	viper.SetDefault("storage_public_dowload", os.Getenv("STORAGE_PUBLIC_DOWNLOAD"))
	viper.SetDefault("storage_private_dowload", os.Getenv("STORAGE_PRIVATE_DOWNLOAD"))
	viper.SetDefault("storage_public_getlibraries", os.Getenv("STORAGE_PUBLIC_GETLIBRARIES"))
	viper.SetDefault("vector_register_blacklist", os.Getenv("VECTOR_REGISTER_BLACKLIST"))
	viper.SetDefault("vector_delete_blacklist", os.Getenv("VECTOR_DELETE_BLACKLIST"))
	viper.SetDefault("vector_search_blacklist", os.Getenv("VECTOR_SEARCH_BLACKLIST"))

	viper.SetDefault("minio_endpoint", os.Getenv("MINIO_ENDPOINT"))
	viper.SetDefault("minio_accesskey_id", os.Getenv("MINIO_ACCESSKEY_ID"))           // export APP_MINIO_ACCESSKEY_ID=
	viper.SetDefault("minio_secret_access_key", os.Getenv("MINIO_SECRET_ACCESS_KEY")) // export APP_MINIO_SECRET_ACCESS_KEY=
	viper.SetDefault("minio_usessl", os.Getenv("MINIO_USESSL"))

	// Assign values to global variables
	SystemAreaUid = viper.GetInt64("system_area_uid")
	Enable_TLS = viper.GetBool("enable_tls")
	Listen_addr = viper.GetString("listen_addr")
	Base_URL = viper.GetString("base_url")
	Service_path = viper.GetString("service_path")
	IWK_set_uri = viper.GetString("jwk_set_uri")

	Enable_sql = viper.GetBool("enable_sql")
	SQL_Host = viper.GetString("sql_host")
	SQL_Port = viper.GetString("sql_port")
	SQL_sslmode = viper.GetString("sql_sslmode")
	SQL_dbname = viper.GetString("sql_dbname")
	SQL_Username = viper.GetString("sql_user")
	SQL_Password = viper.GetString("sql_password")
	SQL_Schema = viper.GetString("sql_schema")

	MQTT_Server_Public_URL = viper.GetString("mqtt_server_public_url")
	MQTT_Server_Host = viper.GetString("mqtt_server_host")
	MQTT_Server_Port = viper.GetInt64("mqtt_server_port")
	MQTT_Server_Usessl = viper.GetBool("mqtt_server_usessl")
	MQTT_Server_ClientID = viper.GetString("mqtt_service_client_id")
	MQTT_Server_Username = viper.GetString("mqtt_service_username")
	MQTT_Server_Password = viper.GetString("mqtt_service_password")
	MQTT_Server_Subcribe_Topic = viper.GetString("mqtt_service_subscribe_topic")

	Client_Cabin_Syns_URI = viper.GetString("client_cabin_sync_uri")
	Client_IAM_Login_URI = viper.GetString("client_iam_login_uri")
	Client_IAM_Client_ID = viper.GetString("client_iam_clientid")
	Client_IAM_Client_Secret = viper.GetString("client_iam_clientsecret")

	Kafka_Enable = viper.GetBool("kafka_enable")
	Kafka_Bootstrap_Ivi = viper.GetString("")
	Kafka_Topic_Producer_Device_Command_Ivi = viper.GetString("kafka_bootstrap_ivi")
	Kafka_Topic_Producer_Device_Command_Ivi_Face_Register = viper.GetString("kafka_topic_producer_device_command_ivi_face_register")
	Kafka_Topic_Consumer_Device_Feedback_Ivi = viper.GetString("kafka_topic_consumer_device_feedback_ivi")
	Kafka_Topic_Consumer_Device_Info_Ivi = viper.GetString("kafka_topic_consumer_device_info_ivi")
	Kafka_Topic_Consumer_AI_VMS_Event = viper.GetString("kafka_topic_consumer_ai_vms_event")
	Kafka_Topic_Consumer_Group = viper.GetString("kafka_consumer_group")

	Storage_Public_Download = viper.GetString("storage_public_download")
	Storage_Private_Download = viper.GetString("storage_private_download")
	Storage_Public_GetLibraries = viper.GetString("storage_public_getlibraries")
	Vector_Register_Blacklist = viper.GetString("vector_register_blacklist")
	Vector_Delete_Blacklist = viper.GetString("vector_delete_blacklist")
	Vector_Search_Blacklist = viper.GetString("vector_search_blacklist")
	// Example usage

	Minio_Endpoint = viper.GetString("minio_endpoint")
	Minio_Accesskey_ID = viper.GetString("minio_accesskey_id")
	Minio_Secret_Access_Key = viper.GetString("minio_secret_access_key")
	Minio_UseSSL = viper.GetBool("minio_usessl")

	log.Println("System Area UID:", viper.GetString("system_area_uid"))
	log.Println("Listen Address:", viper.GetString("listen_addr"))

}
