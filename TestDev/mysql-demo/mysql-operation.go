package mysql_demo

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	//"gorm.io/driver/mysql"
	mysql "database/sql"
	"sync"
	"time"
)

var (
	DBHandle := &gorm.DB{}
	MysqlLock := sync.Mutex{}

)

func InitGlobalOrm()error{
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	//打开数据库连接
	DBHandle,err = gorm.Open("mysql",dsn )
	if err != nil {
		fmt.Println("Open DB by gorm failed, err is %s", err.Error())
		return err
	}

	DBHandle.DB().SetConnMaxLifetime(60*time.Second)
	DBHandle.DB().SetConnMaxIdleTime(60*time.Second)



	return nil
}

func GormDatabaseCRUD(){
	//自动迁移，根据模型创建或更新表结构
	db := DBHandle.AutoMigrate(&User{}, &TestTable{})

	//插入数据
	user := User{Name: "zhangsan", Age: 25}
	db.Create(&user)

	//查询数据
	var users []User
	db.Find(&users)
	fmt.Println("all users is ",users)

	//更新数据
	db.Model(user).Update("Name", "Alice Updated")
	fmt.Println("Updated user: %s+v", user)

	//删除数据
	db.Delete(&user)

	//查询未被删除的数据
	var activeUsers []User
	db.Where("deleted_at IS NULL").Find(&activeUsers)
}

func MysqlInit(){
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

	//连接数据库
	db, err := mysql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Open database failed, err is ",err)
		return
	}
	defer db.Close()

	//验证连接
	if err := db.Ping(); err != nil {
		fmt.Println("Ping database failed, err is ",err)
		return
	}

	//创建表
	createTable(db)

	//插入数据
	insertID := insertData(db, "Alice", 25)
	fmt.Println("Inserted data with id: %d\n", insertID)

	//查询数据
	queryData(db)

	//更新数据
	updateData(db, insertID, "Alice Updated", 30)

	//查询更新后的数据
	queryData(db)

	//删除数据
	deleteData(db, insertID)
}

func createTable(db *mysql.DB){
	query := `
			CREATE TABLE IF NOT EXISTS users (
			    id INT AUTO_INCREMENT PRIMARY KEY,
			    name VARCHAR(128) NOT NULL，
			    age INT NOT NULL,
			    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			`

	if _, err := db.Exec(query);err != nil {
		fmt.Println("create table failed, err is ", err.Error())
	}
	fmt.Println("Create table successful!")
}

func insertData(db *mysql.DB, name string, age int)int64{
	query := `
			INSERT INTO users (name, age) VALUES (?, ?);
			`

	result, err := db.Exec(query, name, age)
	if err != nil {
		fmt.Println("insert data failed,err is ", err.Error())
	}

	insertID, _ := result.LastInsertId()
	return insertID
}

func queryData(db *mysql.DB){
	query := `SELECT id, name, age, created_at FROM users`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("query failed, err is ", err.Error())
	}
	defer rows.Close()

	fmt.Println("Data in users table.")
	for rows.Next(){
		var id int64
		var name string
		var age int
		var createdAt string

		if err := rows.Scan(&id, &name, &age, &createdAt); err != nil {
			fmt.Println("")
		}

		fmt.Println("id: %d, name: %s, age: %d, CreateAt: %s", id, name, age, createdAt)
	}
}

func updateData(db *mysql.DB, insertID int64, name string, age int){
	query := `
			UPDATE users SET name=?, age=? WHERE id=?;
			`

	_, err := db.Exec(query, name, age, insertID)
	if err != nil {
		fmt.Println("update db failed, err is ", err.Error())
	}

	fmt.Println("update database success")
}

func deleteData(db *mysql.DB, insertID int64){
	query := `DELETE FROM users WHERE id=?`

	_,err := db.Exec(query)
	if err != nil {
		fmt.Println("delete from users failed, err is ",err.Error())
	}
	fmt.Println("delete id:%d from users success", insertID)
}