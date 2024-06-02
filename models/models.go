package models

import (
	"time"
)

type Ant struct {
	Id              int       `xorm:"not null pk autoincr index(sz) INT"`
	Number          string    `xorm:"index(sz) VARCHAR(36)"`
	EndoscopeNumber string    `xorm:"index(sz) VARCHAR(100)"`
	EndoscopeType   string    `xorm:"VARCHAR(100)"`
	Operator        string    `xorm:"VARCHAR(100)"`
	PatientName     string    `xorm:"index(sz) VARCHAR(100)"`
	DocName         string    `xorm:"VARCHAR(100)"`
	Diseases        int       `xorm:"default 0 INT"`
	BeginTime       time.Time `xorm:"index(sz) DATETIME"`
	EndTime         time.Time `xorm:"index(sz) DATETIME"`
	TotalCostTime   int       `xorm:"index(sz) INT"`
	EndoscopeInfo   string    `xorm:"index(sz) VARCHAR(100)"`
}

type AntStep struct {
	Id             int    `xorm:"not null pk autoincr INT"`
	Number         string `xorm:"index VARCHAR(36)"`
	Step           string `xorm:"default '' VARCHAR(100)"`
	Stepip         string `xorm:"VARCHAR(50)"`
	CostTime       int    `xorm:"INT"`
	WashingMachine string `xorm:"default '' VARCHAR(100)"`
	Operator       string `xorm:"VARCHAR(50)"`
}

type DeviceInfo struct {
	Id              int    `xorm:"not null pk autoincr INT"`
	EndoscopeNumber string `xorm:"VARCHAR(100)"`
	EndoscopeType   string `xorm:"VARCHAR(100)"`
	EndoscopeInfo   string `xorm:"VARCHAR(100)"`
	Status          int    `xorm:"not null default 1 comment('镜子状态（1 正常  2 维修中）') INT"`
}

type DeviceType struct {
	Id   int    `xorm:"not null pk autoincr INT"`
	Name string `xorm:"VARCHAR(255)"`
	Ids  string `xorm:"VARCHAR(255)"`
}

type DoctorInfo struct {
	Id   int    `xorm:"not null pk autoincr INT"`
	Name string `xorm:"VARCHAR(100)"`
}

type DongguanWriteback struct {
	Id             int       `xorm:"not null pk autoincr INT"`
	Number         string    `xorm:"VARCHAR(36)"`
	Namepatient    string    `xorm:"VARCHAR(100)"`
	Patientvisitid string    `xorm:"VARCHAR(100)"`
	Checkreportid  string    `xorm:"VARCHAR(100)"`
	Rfidno         string    `xorm:"VARCHAR(100)"`
	Bindmirrostate string    `xorm:"VARCHAR(100)"`
	BeginTime      time.Time `xorm:"default CURRENT_TIMESTAMP TIMESTAMP"`
}

type EndoscopeSign struct {
	Id              int       `xorm:"not null pk autoincr INT"`
	EndoscopeNumber string    `xorm:"default '' VARCHAR(100)"`
	Time            time.Time `xorm:"DATETIME"`
	AntNumber       string    `xorm:"default '' VARCHAR(36)"`
}

type EndoscopyWriteback struct {
	Id             int       `xorm:"not null pk autoincr INT"`
	Number         string    `xorm:"VARCHAR(255)"`
	Uniqueid       string    `xorm:"not null unique VARCHAR(255)"`
	Namepatient    string    `xorm:"VARCHAR(100)"`
	Patientvisitid string    `xorm:"VARCHAR(100)"`
	Checkreportid  string    `xorm:"VARCHAR(100)"`
	Rfidno         string    `xorm:"index VARCHAR(100)"`
	Bindmirrostate string    `xorm:"VARCHAR(100)"`
	BeginTime      time.Time `xorm:"default CURRENT_TIMESTAMP TIMESTAMP"`
}

type LeakDetectionRecord struct {
	Id         int    `xorm:"not null pk autoincr INT"`
	Number     string `xorm:"VARCHAR(36)"`
	Conclusion string `xorm:"default '完好' VARCHAR(10)"`
	Leakparts  string `xorm:"default '' VARCHAR(512)"`
}

type LiquidDetectionRecord struct {
	Id             int       `xorm:"not null pk autoincr UNSIGNED INT"`
	Deviceid       int       `xorm:"not null UNSIGNED INT"`
	LiquidType     string    `xorm:"default '' VARCHAR(100)"`
	ValidityPeriod int       `xorm:"default 0 INT"`
	LimitNumber    int       `xorm:"default 0 INT"`
	Date           time.Time `xorm:"DATETIME"`
	Operator       string    `xorm:"default '' VARCHAR(100)"`
	Conclusion     string    `xorm:"default '正常' VARCHAR(10)"`
	OperatorType   string    `xorm:"default '' VARCHAR(50)"`
	OperatorInfo   string    `xorm:"default '' VARCHAR(100)"`
	UsageCount     int       `xorm:"default 0 INT"`
}

type LiquidDevice struct {
	Id             int    `xorm:"not null pk autoincr UNSIGNED INT"`
	DeviceName     string `xorm:"not null default '' VARCHAR(100)"`
	DeviceType     string `xorm:"default '' VARCHAR(100)"`
	DeviceInfo     string `xorm:"default '' VARCHAR(512)"`
	LiquidType     string `xorm:"default '' VARCHAR(100)"`
	ValidityPeriod int    `xorm:"default 0 INT"`
	LimitNumber    int    `xorm:"default 0 INT"`
}

type OperatorInfo struct {
	Id     int    `xorm:"not null pk autoincr INT"`
	Number string `xorm:"VARCHAR(100)"`
	Name   string `xorm:"VARCHAR(100)"`
}

type PreprocessLog struct {
	Id              int       `xorm:"not null pk autoincr index(序列号) INT"`
	EndoscopeNumber string    `xorm:"comment('内窥镜号') index(序列号) VARCHAR(100)"`
	Operator        string    `xorm:"VARCHAR(100)"`
	Time            time.Time `xorm:"DATETIME"`
	Ip              string    `xorm:"VARCHAR(100)"`
	Number          string    `xorm:"comment('回写洗消编号') VARCHAR(100)"`
}

type Program struct {
	Id            int    `xorm:"not null pk autoincr INT"`
	Name          string `xorm:"VARCHAR(100)"`
	TotalCostTime int    `xorm:"default 0 INT"`
}

type ProgramList struct {
	Id       int    `xorm:"not null pk autoincr INT"`
	Name     string `xorm:"VARCHAR(100)"`
	Step     string `xorm:"VARCHAR(100)"`
	CostTime int    `xorm:"INT"`
}
type StoreLog struct {
	Id              int       `xorm:"not null pk autoincr INT"`
	EndoscopeNumber string    `xorm:"default '' VARCHAR(100)"`
	Operator        string    `xorm:"default '' VARCHAR(255)"`
	Time            time.Time `xorm:"DATETIME"`
	Access          string    `xorm:"default '' VARCHAR(255)"`
}
type Step struct {
	Id   int    `xorm:"not null pk autoincr INT"`
	Step string `xorm:"VARCHAR(50)"`
}
type TimePlan struct {
	Id                int    `xorm:"not null pk autoincr INT"`
	Name              string `xorm:"VARCHAR(100)"`
	RinseTime         int    `xorm:"INT"`
	RinseSolution     string `xorm:"VARCHAR(100)"`
	DisinfectTime     int    `xorm:"INT"`
	SteriliseTime     int    `xorm:"INT"`
	DisinfectSolution string `xorm:"VARCHAR(100)"`
}
type RepairRecords struct {
	Id                 int        `gorm:"not null pk autoincr INT"`
	EndoscopeNumber    string     `gorm:"not null"`
	ResponsiblePerson  string     `gorm:"not null;size:100"`
	ProblemDescription string     `gorm:"not null"`
	ProblemLocation    string     `gorm:"size:255"`
	SendOutDate        time.Time  `gorm:"not null"`
	ReturnDate         *time.Time `gorm:"DATETIME"`
	RepairCost         float64    `gorm:"decimal"`
	RepairNotes        string     `gorm:"VARCHAR(255)"`
}

type User struct {
	Id       int    `xorm:"not null pk autoincr INT"`
	Username string `xorm:"VARCHAR(100)"`
	Password string `xorm:"CHAR(36)"`
}

type Users struct {
	Id        int       `xorm:"not null pk autoincr INT"`
	Name      string    `xorm:"VARCHAR(255)"`
	Password  string    `xorm:"VARCHAR(255)"`
	Email     string    `xorm:"VARCHAR(255)"`
	CreatedAt time.Time `xorm:"DATETIME"`
	UpdatedAt time.Time `xorm:"DATETIME"`
	DeletedAt time.Time `xorm:"DATETIME"`
	ID        int       `xorm:"INT"`
}

type WritebackLog struct {
	Id             int       `xorm:"not null pk autoincr INT"`
	Number         string    `xorm:"not null unique VARCHAR(36)"`
	Rfidno         string    `xorm:"VARCHAR(36)"`
	Patient        string    `xorm:"VARCHAR(36)"`
	Doctor         string    `xorm:"VARCHAR(36)"`
	Diseases       string    `xorm:"VARCHAR(36)"`
	Operator       string    `xorm:"VARCHAR(36)"`
	Bindmirrostate string    `xorm:"VARCHAR(10)"`
	Datetime       time.Time `xorm:"default CURRENT_TIMESTAMP TIMESTAMP"`
}

type XianWriteback struct {
	Id             int    `xorm:"not null pk autoincr INT"`
	Number         string `xorm:"VARCHAR(36)"`
	Patientid      string `xorm:"VARCHAR(100)"`
	Patientno      string `xorm:"VARCHAR(100)"`
	Namepatient    string `xorm:"VARCHAR(100)"`
	Sex            string `xorm:"VARCHAR(100)"`
	Age            string `xorm:"VARCHAR(100)"`
	Phone          string `xorm:"VARCHAR(100)"`
	Patientvisitid string `xorm:"VARCHAR(100)"`
	Hbsag          string `xorm:"VARCHAR(100)"`
	Hcv            string `xorm:"VARCHAR(100)"`
	Hiv            string `xorm:"VARCHAR(100)"`
	Tp             string `xorm:"VARCHAR(100)"`
	Hbcab          string `xorm:"VARCHAR(100)"`
	Hbeab          string `xorm:"VARCHAR(100)"`
	Hbsab          string `xorm:"VARCHAR(100)"`
	Checkreportid  string `xorm:"VARCHAR(100)"`
	Scopeno        string `xorm:"VARCHAR(100)"`
	Scopetype      string `xorm:"VARCHAR(100)"`
	Scopeinnerno   string `xorm:"VARCHAR(100)"`
	Rfidno         string `xorm:"VARCHAR(100)"`
	Bindmirrostate string `xorm:"VARCHAR(100)"`
}
