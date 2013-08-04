package register

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"labix.org/v2/mgo"
  "os"
)

func Register(useremail, password string) {
	//db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/msgstory")
  dburl := os.Getenv("CLEARDB_DATABASE_URL")
	db, err := sql.Open("mysql", dburl[8:])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO users (useremail, password) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(useremail, password)
	if err != nil {
		panic(err.Error())
	} else {
		createPerson(useremail)
	}
}

func createPerson(userEmail string) {
	mdb, err := mgo.Dial("localhost")
	if err != nil {
		panic(err.Error())
	}

	mdb.SetMode(mgo.Monotonic, true)

	// person := User{}

	// c := mdb.DB("msgme").C("jove")
	// jove := User.M
	// // err = c.Find(bson.M{"firstName": "Jasdeep1"}).One(&result)
	// err = c.Insert(&jove{nil, nil, userEmail})
	// if err != nil {
	// 	panic(err)
	// }
}
