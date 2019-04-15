package db

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

/* 数据格式 */

//TypeOfWork ...工种数据表格式
type TypeOfWork struct {
	gorm.Model
	TypeOfWorkField string `gorm:"type:varchar(30)"`
}

//StaffInfo ...员工信息表格式
type StaffInfo struct {
	gorm.Model
	Ref               string    `gorm:"type:varchar(20)"`            //员工编号
	Name              string    `gorm:"type:varchar(20)"`            //姓名
	Department        string    `gorm:"size:20;default:'矿管公司项目二部'"`  //所属部门
	Sex               string    `gorm:"type:varchar(2);default:'男'"` //性别
	Education         string    `gorm:"type:varchar(10)"`            //文化程度
	IDCard            string    `gorm:"type:varchar(20)"`            //身份证号
	TypeOfWorkID      string    `gorm:"type:varchar(20)"`            //工种
	InitialDate       time.Time `gorm:"default:null"`                //初领证日期
	TrainingTime      time.Time `gorm:"default:null"`                //培训时间
	NextReviewDate    time.Time `gorm:"default:null"`                //下次需复审日期
	DateOfReplacement time.Time `gorm:"default:null"`                //换证日期
	IDNumber          string    `gorm:"type:varchar(20)"`            //证件编号
	Details           string    `gorm:"type:varchar(100)"`           //详情
}

/* 数据库数据 */

//JsonTypeOfWork ... 工种数据
var jsonTypeOfWork = `[
	{"TypeOfWorkField":"安全检查工"},
	{"TypeOfWorkField":"采煤班组长"},
	{"TypeOfWorkField":"绞车司机"},
	{"TypeOfWorkField":"井下电钳工"},
	{"TypeOfWorkField":"起重机械操作工"},
	{"TypeOfWorkField":"金属焊接切割工"},
	{"TypeOfWorkField":"乳化液泵站工"},
	{"TypeOfWorkField":"输送机司机"},
	{"TypeOfWorkField":"信号把钩工"},
	{"TypeOfWorkField":"窄轨电机车司机"}
]`

/* 创建函数 */

//TgCreate ... 创建特殊工种数据库
func TgCreate() {
	db, err := gorm.Open("sqlite3", "../db/tgDb.db")
	if err != nil {
		panic("faild open tgdb.db")
	}
	db.AutoMigrate(&TypeOfWork{}, &StaffInfo{})
	fmt.Println("建立员工信息表、工种类型表共二个数据库完成")
	var count = 0
	db.Find(&TypeOfWork{}).Count(&count)
	if count == 0 {
		var typeOfWorkArray []TypeOfWork
		jerr := json.Unmarshal([]byte(jsonTypeOfWork), &typeOfWorkArray)
		if jerr != nil {
			panic("工种json数据转换为struct数组错误")
		}
		for _, v := range typeOfWorkArray {
			db.Create(&v)
		}
		fmt.Println("工种数据创建成功")
	} else {
		println("工种记录已存在，工种数据不再创建")
	}
	db.DropTableIfExists(&StaffInfo{})
	db.CreateTable(&StaffInfo{})
	file, csvErr := os.Open("../db/tgdb.csv")
	if csvErr != nil {
		fmt.Println("打开tgDb.csv文件错误", csvErr)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var staffInfoRecord StaffInfo
	var sirErr error
	var str string
	for i := 0; true; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取tgDb.csv文件错误", err)
			break
		}
		if i == 0 {
			continue
		}
		local, _ := time.LoadLocation("Local")
		var r8, _ = time.ParseInLocation("2006/1/2", record[8], local)
		var r9, _ = time.ParseInLocation("2006/1/2", record[9], local)
		var r10, _ = time.ParseInLocation("2006/1/2", record[10], local)
		var r11, _ = time.ParseInLocation("2006/1/2", record[11], local)
		var sr8 = r8.Format(time.RFC3339)
		var sr9 = r9.Format(time.RFC3339)
		var sr10 = r10.Format(time.RFC3339)
		var sr11 = r11.Format(time.RFC3339)
		// if sr8 == "0001-01-01T00:00:00Z" {
		// 	fmt.Println("未编码空值sr8：", sr8)
		// }
		str = `{"ref":` + `"` + record[2] + `"` + "," + `"name":` + `"` + record[3] + `"` + "," + `"Education":` + `"` + record[5] + `"` + "," + `"IDCard":` + `"` + record[6] + `"` + "," + `"TypeOfWorkID":` + `"` + record[7] + `"` + "," + `"InitialDate":` + `"` + sr8 + `"` + "," + `"TrainingTime":` + `"` + sr9 + `"` + "," + `"NextReviewDate":` + `"` + sr10 + `"` + "," + `"DateOfReplacement":` + `"` + sr11 + `"` + "," + `"IDNumber":` + `"` + record[12] + `"` + "," + `"Details":` + `"` + record[13] + `"` + "}"

		sirErr = json.Unmarshal([]byte(str), &staffInfoRecord)
		staffInfoRecord.ID = 0
		if sirErr != nil {
			fmt.Println("员工信息表编码错误:", sirErr)
			return
		}
		//fmt.Printf("%+v", staffInfoRecord)
		db.Create(&staffInfoRecord)
	}
	type na struct {
		InitialDate time.Time
		Name        string
		SEX         string
	}
	var names []na
	db.Model(StaffInfo{}).Select("initial_date, sex, name").Scan(&names)
	fmt.Printf("%+v", names)
}
