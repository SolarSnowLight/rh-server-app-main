package main

import (
	"context"
	"encoding/json"
	"fmt"
	mainserver "main-server"
	"main-server/config"
	handler "main-server/pkg/handler"
	repository "main-server/pkg/repository"
	"main-server/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/spf13/viper"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		logrus.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		logrus.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logrus.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// @title Rental Housing
// @version 1.0
// description Проект для аренды жилья

// @host localhost:5000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	/*data, err := ioutil.ReadFile("./config/client_secret.json")
	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}

	client := conf.Client(context.Background())
	srvSheets, err := sheets.New(client)
	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}

	spreadsheetID := "1m9_IhCWtlAuifx0zd1QZHxg_E6U_T_cFHlaxkdzODpo"
	readRange := "A1:AC66"
	resp, err := srvSheets.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()

	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			for _, cell := range row {
				fmt.Printf("%s\n", cell)
			}
		}
	}*/

	/*data, err := ioutil.ReadFile("./config/client_secret.json")
	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}

	client := conf.Client(context.Background())

	service2 := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service2.FetchSpreadsheet("1m9_IhCWtlAuifx0zd1QZHxg_E6U_T_cFHlaxkdzODpo")
	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}
	sheet, err := spreadsheet.SheetByIndex(0)

	if err != nil {
		logrus.Fatalf("error : %s", err.Error())
	}

	f := excelize.NewFile()

	for _, row := range sheet.Rows {
		for _, cell := range row {

			f.SetCellValue("Sheet1", cell.Pos(), cell.Value)
			// f.SetCellStyle("Sheet1", )
		}
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}*/

	/* Init config */
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable: %s", err.Error())
	}

	/* Init logger */
	logrus.SetFormatter(new(logrus.JSONFormatter))

	fileError, err := os.OpenFile(viper.GetString("paths.logs.error"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileError,
			LogLevels: []logrus.Level{
				logrus.ErrorLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	fileInfo, err := os.OpenFile(viper.GetString("paths.logs.info"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileInfo,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	fileWarn, err := os.OpenFile(viper.GetString("paths.logs.warn"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileWarn,
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	fileFatal, err := os.OpenFile(viper.GetString("paths.logs.fatal"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileFatal,
			LogLevels: []logrus.Level{
				logrus.FatalLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	/* Connection in database */
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.username"),
		"iJ#Q@LKpkawf-)$1l25,m12l5kkm<TNN@IY)*DGjnlQ#BN",
		viper.GetString("db.dbname"),
		viper.GetString("db.port"),
		viper.GetString("db.sslmode"),
	)

	/* Init PERM model */
	dbAdapter, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dns,
	}), &gorm.Config{})

	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(dbAdapter, &config.MisuRule{}, viper.GetString("rules_table_name"))

	if err != nil {
		logrus.Fatalf("failed to initialize adapter by db with custom table: %s", err.Error())
	}

	enforcer, err := casbin.NewEnforcer(viper.GetString("paths.perm_model"), adapter)

	if err != nil {
		logrus.Fatalf("failed to initialize new enforcer: %s", err.Error())
	}

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	/* Init oauth2 services */
	config.InitOAuth2Config()
	config.InitVKAuthConfig()

	/* Dependency injection */
	repos := repository.NewRepository(db, enforcer)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := new(mainserver.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Rental Housing Main Server Started")

	// Реализация Graceful Shutdown
	// Блокировка функции main с помощью канала os.Signal
	quit := make(chan os.Signal, 1)

	// Запись в канал, если процесс, в котором выполняется приложение
	// получит сигнал SIGTERM или SIGINT
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Чтение из канала, блокирующая выполнение функции main
	<-quit

	logrus.Print("Rental Housing Main Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

	/* Closed files for logger */
	if err := fileError.Close(); err != nil {
		logrus.Error(err.Error())
	}

	if err := fileWarn.Close(); err != nil {
		logrus.Error(err.Error())
	}

	if err := fileInfo.Close(); err != nil {
		logrus.Error(err.Error())
	}

	if err := fileFatal.Close(); err != nil {
		logrus.Error(err.Error())
	}
}

// Init config
func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
