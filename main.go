package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/txross1993/superman-api/api"
	"github.com/txross1993/superman-api/db"
	"github.com/txross1993/superman-api/geolocate"
	"github.com/txross1993/superman-api/superman"
)

func main() {
	var apiCfg api.Config
	var geoliteRepository string
	var dataPath string
	flag.StringVar(&apiCfg.Host, "host", getEnvOrDefault("HOST", "0.0.0.0"), "Provide the bind address for hosting the api")
	flag.StringVar(&apiCfg.Port, "port", getEnvOrDefault("PORT", "8080"), "Provide the bind port for hosting the api")
	flag.StringVar(&geoliteRepository, "geodb", getEnvOrDefault("GEODB", "GeoLite2-City_20200602/GeoLite2-City.mmdb"), "Provide the fully qualified path to the GeoLite2 database *.mmdb file")
	flag.StringVar(&dataPath, "dbpath", getEnvOrDefault("DBPATH", "local-db"), "Provide the fully qualified path to the sqlite database host directory")
	flag.Parse()

	geoSvc, err := geolocate.NewGeoService(geoliteRepository)
	if err != nil {
		log.Fatal(err)
	}
	defer geoSvc.Close()

	localDb := path.Join(dataPath, "local.db")

	sqlDB, err := db.InitDB(localDb)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	superman := superman.NewService(geoSvc, sqlDB)
	apiCfg.Superman = superman

	api := api.NewAPI(apiCfg)

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}

}

func getEnvOrDefault(val, defaultVal string) string {
	if env := os.Getenv(val); env != "" {
		return env
	}

	return defaultVal
}
