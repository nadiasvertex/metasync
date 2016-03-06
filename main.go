package main

import "fmt"
import "flag"

var databaseDataFolder = flag.String("db-data-folder", "./db/data", "path to the folder where the database data will be stored.")
var databaseEngineFolder = flag.String("db-engine-folder", "./db/engine", "path to the folder where the database data will be stored.")

func main() {
	fmt.Printf("metasync server\ndb data folder=%v\ndb engine folder=%v\n", *databaseDataFolder, *databaseEngineFolder)

}
