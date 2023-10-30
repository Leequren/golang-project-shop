package databaseFuncs

import (
	"database/sql"
	"fmt"
	"log"
	"testingGorillaw/HashFuncs"
)

type Person struct {
	idUser             int
	loginUser          string
	hashedPasswordUser string
}
type allInfoMonitor struct {
	IdMonitor             int     `json:"idMonitor"`
	NameMonitor           string  `json:"nameMonitor"`
	VoltageMonitor        float64 `json:"voltageMonitor"`
	UseGSyncPremMonitor   bool    `json:"useGSyncPremMonitor"`
	CurvedMonitor         bool    `json:"curvedMonitor"`
	DisplayMonitorId      int     `json:"displayMonitorId"`
	NameDisplay           string  `json:"nameDisplay"`
	LengthDiagonalDisplay float64 `json:"lengthDiagonalDisplay"`
	ResolutionDisplay     string  `json:"resolutionDisplay"`
	MatrixDisplay         string  `json:"matrixDisplay"`
	UseGSync              bool    `json:"useGSync"`
}
type ArrayMonitors struct {
	Monitors []allInfoMonitor `json:"monitors"`
}

func CreatePostgresConn(username, password, host, port, dbname, sslMode string) *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", username, password, dbname, sslMode, host, port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(db)
	return db
}

func InsertNewPerson(db *sql.DB, login, hashedPassword string) string {
	sqlExpr := fmt.Sprintf("insert into person(loginUser, hashedPasswordUser) values('%s', '%s')", login, hashedPassword)
	_, err := db.Exec(sqlExpr)
	log.Println(err)
	if err != nil {
		return "Insert Error"
	}
	return "Inserted"
}

func SelectAllPersons(db *sql.DB) (data []Person, err error) {
	rows, err := db.Query("select idUser, loginUser, hashedPasswordUser from person")
	if err != nil {
		return []Person{}, err
	}
	var persons []Person
	defer rows.Close()

	for rows.Next() {
		var person Person
		err := rows.Scan(&person.idUser, &person.loginUser, &person.hashedPasswordUser)
		if err != nil {
			return []Person{}, err
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func CheckUserByLoginInDB(db *sql.DB, login string) (bool, error) {
	sqlExrp := fmt.Sprintf("select idUser, loginUser, hashedPasswordUser from person where loginUser = '%s'", login)
	//log.Println(sqlExrp)
	rows, err := db.Query(sqlExrp)
	if err != nil {
		log.Println("Error")
	}
	var person Person
	errScan := rows.Scan(&person.idUser, &person.loginUser, &person.hashedPasswordUser)
	if err != nil {
		return false, errScan
	}
	log.Println(person)

	//log.Println(rows.Next())
	for rows.Next() {
		log.Println(123)
		err := rows.Scan(&person.idUser, &person.loginUser, &person.hashedPasswordUser)
		if err != nil {
			return false, err
		}
		//log.Println(person)
	}
	//log.Print(person)
	return person.idUser != 0, nil
}

func CheckSignInData(db *sql.DB, login string, password string) (bool, error) {
	sqlExrp := fmt.Sprintf("select idUser, loginUser, hashedPasswordUser from person where loginUser = '%s'", login)
	//log.Println(sqlExrp)
	rows, err := db.Query(sqlExrp)
	if err != nil {
		log.Println("Error")
	}
	var person Person
	errScan := rows.Scan(&person.idUser, &person.loginUser, &person.hashedPasswordUser)
	if err != nil {
		return false, errScan
	}
	log.Println(person)

	//log.Println(rows.Next())
	for rows.Next() {
		log.Println(123)
		err := rows.Scan(&person.idUser, &person.loginUser, &person.hashedPasswordUser)
		if err != nil {
			return false, err
		}
		//log.Println(person)
	}
	//log.Print(person)
	if person.idUser != 0 {
		match := HashFuncs.CheckPasswordHash(password, person.hashedPasswordUser)
		return match, nil
	}
	return person.idUser != 0, nil
}

func InsertDisplay(db *sql.DB, NameDisplay string, lengthDiagonal float64, resolution, matrix string, useGSync bool) int {
	log.Println(lengthDiagonal, resolution, matrix, useGSync)
	sqlExpr := fmt.Sprintf("insert into display(NameDisplay, LengthDiagonalDisplay, ResolutionDisplay, MatrixDisplay, UseGSync) values ('%v', '%v', '%v', '%v', '%v') returning IdDisplay;", NameDisplay, lengthDiagonal, resolution, matrix, useGSync)
	var idDisplay int
	log.Println(sqlExpr)
	err := db.QueryRow(sqlExpr).Scan(&idDisplay)
	if err != nil {
		log.Println(err)
		return 0
	}
	return idDisplay
}

func InsertMonitor(db *sql.DB, NameMonitor string, voltageMonitor float64, UseGSyncPremMonitor, CurvedMonitor bool, DisplayMonitorId int) (string, error) {
	sqlExpr := fmt.Sprintf("insert into monitor(nameMonitor, VoltageMonitor, UseGSyncPremMonitor, CurvedMonitor, DisplayMonitorId) values('%v','%v', '%v', '%v', '%v')", NameMonitor, voltageMonitor, UseGSyncPremMonitor, CurvedMonitor, DisplayMonitorId)
	_, err := db.Query(sqlExpr)
	if err != nil {
		log.Println(err)
		return "not inserted", err
	}
	return "inserted", nil
}

func GetAllMonitors(db *sql.DB) (*ArrayMonitors, error) {
	sqlExrp := "select idmonitor, namemonitor, voltagemonitor, usegsyncpremmonitor, curvedmonitor, iddisplay, namedisplay, lengthdiagonaldisplay, resolutiondisplay, matrixdisplay, usegsync from monitor join display on display.iddisplay = monitor.displaymonitorid;"

	rows, err := db.Query(sqlExrp)
	if err != nil {
		return nil, err
	}
	var monitors ArrayMonitors
	for rows.Next() {
		var monitor allInfoMonitor
		err := rows.Scan(&monitor.IdMonitor, &monitor.NameMonitor, &monitor.VoltageMonitor, &monitor.UseGSyncPremMonitor, &monitor.CurvedMonitor, &monitor.DisplayMonitorId, &monitor.NameDisplay, &monitor.LengthDiagonalDisplay, &monitor.ResolutionDisplay, &monitor.MatrixDisplay, &monitor.UseGSync)
		if err != nil {
			return nil, err
		}

		//log.Println("[ALL MONITOR INFO DATA]", monitor)
		monitors.Monitors = append(monitors.Monitors, monitor)
	}

	return &monitors, nil
}

func CreateDBWithDefaultConfig() *sql.DB {
	db := CreatePostgresConn("postgres", "1234", "localhost", "5432", "TechShop", "disable")
	return db
}
