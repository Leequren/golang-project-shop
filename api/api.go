package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testingGorillaw/HashFuncs"
	"testingGorillaw/databaseFuncs"
	"testingGorillaw/env"
)

type personInfo struct {
	LoginUser string `json:"login"`
	Password  string `json:"password"`
}
type infoMatrixFrontend struct {
	Name       string  `json:"name"`
	Diagonal   float64 `json:"diagonal"`
	Resolution string  `json:"resolution"`
	Type       string  `json:"type"`
	UseGSync   bool    `json:"useGSync"`
}
type infoMonitorFrontend struct {
	Name         string  `json:"name"`
	Voltage      float64 `json:"voltage"`
	UseGSyncPrem bool    `json:"useGSyncPrem"`
	Curved       bool    `json:"curved"`
	DisplayId    int     `json:"displayId"`
}
type allInfoMonitorFrontend struct {
	Name       string  `json:"name"`
	Voltage    float64 `json:"voltage"`
	GSync      bool    `json:"gsync"`
	GSyncPrem  bool    `json:"gsync_prem"`
	Curved     bool    `json:"curved"`
	Diagonal   float64 `json:"diagonal"`
	Resolution string  `json:"resolution"`
	Matrix     string  `json:"matrix"`
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

func CreateUserPost(w http.ResponseWriter, r *http.Request) {
	//body, err := r.Body
	db := databaseFuncs.CreateDBWithDefaultConfig()
	data, _ := databaseFuncs.SelectAllPersons(db)
	log.Println(data)
	var pInfo personInfo
	//log.Print(r.Body)
	err := json.NewDecoder(r.Body).Decode(&pInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Println(pInfo)
	hashPassword, _ := HashFuncs.HashPassword(pInfo.Password)
	res := databaseFuncs.InsertNewPerson(db, pInfo.LoginUser, hashPassword)
	log.Print(res)
	if res == "Insert Error" {
		io.WriteString(w, "Don't created user")
		return
	}
	io.WriteString(w, "Created user")
}

func SignInPost(w http.ResponseWriter, r *http.Request) {

	db := databaseFuncs.CreateDBWithDefaultConfig()
	var pInfo personInfo
	err := json.NewDecoder(r.Body).Decode(&pInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	match, err := databaseFuncs.CheckSignInData(db, pInfo.LoginUser, pInfo.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if match {
		session, _ := env.Store.Get(r, "session.id")
		session.Values["authenticated"] = true
		session.Save(r, w)
		io.WriteString(w, "Access")
	} else {
		io.WriteString(w, "Denied")
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := env.Store.Get(r, "session.id")
	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Write([]byte("Logout Successful"))
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	session, _ := env.Store.Get(r, "session.id")
	log.Println(session)
	io.WriteString(w, "Данные в логах")
}
func InsertMonitorWithDataDisplay(w http.ResponseWriter, r *http.Request) {
	db := databaseFuncs.CreateDBWithDefaultConfig()
	var pMonitor allInfoMonitorFrontend

	err := json.NewDecoder(r.Body).Decode(&pMonitor)
	if err != nil {
		log.Println(err)
		io.WriteString(w, err.Error())
		return
	}
	log.Println("[INFO MONITOR]", pMonitor)

	idDisplay := databaseFuncs.InsertDisplay(db, "", pMonitor.Diagonal, pMonitor.Resolution, pMonitor.Matrix, pMonitor.GSync)
	log.Println("[ID DISPLAY]: ", idDisplay)
	_, err = databaseFuncs.InsertMonitor(db, pMonitor.Name, pMonitor.Voltage, pMonitor.GSyncPrem, pMonitor.Curved, idDisplay)
	if err != nil {
		io.WriteString(w, err.Error())
	}

	io.WriteString(w, "Monitor inserted")
}

func InsertMatrix(w http.ResponseWriter, r *http.Request) {
	db := databaseFuncs.CreateDBWithDefaultConfig()
	var pMatrix infoMatrixFrontend

	err := json.NewDecoder(r.Body).Decode(&pMatrix)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	log.Println("[INFO MATRIX]", pMatrix)

	idDisplay := databaseFuncs.InsertDisplay(db, pMatrix.Name, pMatrix.Diagonal, pMatrix.Resolution, pMatrix.Type, pMatrix.UseGSync)
	if idDisplay == 0 {
		io.WriteString(w, "Матрица не добавлена")
		return
	}
	io.WriteString(w, "Matrix inserted")
}

func InsertMonitor(w http.ResponseWriter, r *http.Request) {
	db := databaseFuncs.CreateDBWithDefaultConfig()
	var pMonitor infoMonitorFrontend

	err := json.NewDecoder(r.Body).Decode(&pMonitor)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	log.Println("[INFO MONITOR]", pMonitor)

	_, err = databaseFuncs.InsertMonitor(db, pMonitor.Name, pMonitor.Voltage, pMonitor.UseGSyncPrem, pMonitor.Curved, pMonitor.DisplayId)
	if err != nil {
		io.WriteString(w, err.Error())
	}

	io.WriteString(w, "Monitor inserted")
}

func GetAllMonitors(w http.ResponseWriter, r *http.Request) {
	db := databaseFuncs.CreateDBWithDefaultConfig()
	monitors, err := databaseFuncs.GetAllMonitors(db)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	jsonDataMonitors, err := json.Marshal(&monitors)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	w.Write(jsonDataMonitors)
}
