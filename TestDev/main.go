package TestDev

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/waitgroup"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"net"
	"net/http"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/jinzhu/gorm"


	iris "github.com/kataras/iris/v12"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*1.参数解析、配置加载 */
//a.参数定义
type TestString string

type TestStruct struct {
	Var1    string	`json: var1, omitempty`
	Slice1	[]int	`json: Slice1, omitempty`
	Map1	map[string]int
	Chan1	chan struct{}
	Func1   func()
	Struct1 struct{}
}

//b.初始化参数
func AddFlags(fs *pflag.FlagSet){
	testVar := 1
	testStruct := TestStruct{}
	testString := "test"

	print("testvar is "+testVar)

	flagSet := pflag.NewFlagSet("testFlagSet",pflag.ExitOnError)
	flagSet.StringVar(&fs.,"test", )


	flagSet.Parse()
}

//c.解析参数
argparser.Parse()

//d.校验参数
func VerifyTestStruct(v1 string){
	print("var is "+v1)
	v1Type := reflect.TypeOf(v1)
	print(v1Type)

	return
}

/*2.优雅终止 */
signal.TERM

func gracefulExit(ctx) {

	select {
	case <-c == 1:
		gracefulExit(ctx)
	case <-c == 2:
		return
	}
}


/* 3.lead选举 */
func leaderElection(ctx context.Context, leaseName string, leaseNamespace string, c *client.Client){
	//a.定义租约配置
	identity := "test"
	lock := &resourcelock.LeaseLock{
		LeaseMeta: v1.ObjectMeta{
			Name: leaseName,
			Namespace: leaseNamespace,
		},
		Client: c,
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: identity,

		},
	}

	leaderElectionConfig := leaderelection.LeaderElectionConfig{
		Lock: lock,
		LeaseDuration: 60*time.Second,
		RenewDeadline: 15*time.Second,
		RetryPeriod: 5*time.Second,

		Name: "test",

		//选举后执行的函数，即进入函数主体
		Callbacks: leaderelection.LeaderCallbacks{
			OnNewLeader: func(identity string){
				print("NewLeader"+identity)
			},

			OnStoppedLeading: func(){
				fmt.Print("Stoppe leading!")
			},

			OnStartedLeading: func(ctx context.Context){
				fmt.Print("Start Leading!")
			},
		},

	}

	leaderelection.RunOrDie(ctx, leaderElectionConfig)
	//b.定义
}

/*4.select/switch */
func SelectTest(ctx context.Context){
	c := make(chan struct{})
	c1 := make(chan string)

	select{
	case x:=<-c:
		fmt.Print("test result %s",x)
	case y:=<-c1:
		fmt.Print("hahaha"+y)
	default:
		fmt.Print("SelectTest finish!")
	}
}

func SwitchTest(ctx context.Context){
	c := make(chan struct{})
	c1 := make(chan string)

	switch <-c {
	case struct{}{}:
		fmt.Print("switch result 1")
	default:
		fmt.Print("switch 1 finish !")
	}

	switch <-c1 {
	case "name":
		fmt.Print("switch result is name")
	case "password":
		fmt.Print("switch result is password")
	default:
		fmt.Print("switch 2 finish!")
	}
}

/* 5.context/reflect */
func ContextProcess(){
	ctx,cancel := context.WithCancel(context.Background())


	ctx1,cancel1 := context.WithDeadline(ctx, 3*time.Second)


	ctx2,cancel2 := context.WithTimeout(ctx, 5*time.Second)

	key := "key"
	value := "value"
	ctx3,cancel3 := context.WithValue(ctx, &key, &value)
}

/* 6.DB连接、读写 */
//表结构定义
type testTable struct{
	gorm.Model
	Uuid string `json:"colume:uuid,type:varchar(128)"`
	name string	`json:"colume:name,type:varchar(128)"`
	password string `json:"colume:name,type:varchar(128）"`
}

//读数据
func DatabaseRead(host net.IP, port int, name string, passwd string){
	database := "testDB"
	table := "testTable"

	mysqlDev := "root:root@tcp(localhost:3306)/niwo?charset=utf8&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=60s"

	sqlQuery := "select * from testTable where name=yes"

	dbHandle, err := gorm.Open("mysql", mysqlDev)
	if err != nil{
		fmt.Print("Open mysql failed, err is "+err.Error())
	}
	//查询使用方法，1. where函数，2.先写好查询语句
	dbHandle.Where()
}

//写数据
func DatabaseWrite(host net.IP, port int, name string, passwd string){
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Database open failed, err is ",err.Error())
	}

	query := `
		SELECT * FROM users WHERE name=?
		`

	resultRaws, err := db.Query(query, "alice")
	if err != nil {
		fmt.Println("Select Query failed, err is ",err.Error())
	}
	defer resultRaws.Close()

	var insertID int64
	var userName string
	var age int
	var createdAt string
	err1 := resultRaws.Scan(&insertID, &userName, &age, &createdAt)
	if err1 != nil {
		fmt.Println("result scan failed, err is ", err1.Error())
	}
}

/* 7.MQ通信 */

/* 8.net IO */
hostIP := "127.0.0.1"
port := 12345
//1.普通网络通信
func testSocket(hostIP string, port int){
	//用Dial函数
	address := hostIP+":"+strconv.Itoa(port)
	conn,err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		fmt.Print("test tcp socket by DialTimeout failed! err is %s", err.Error())
	}


	//用DialIP函数
	localAddress := &net.IPAddr{
		IP: net.IP{127,0,0,1},
		Zone: "local",
	}
	remoteAddress := &net.IPAddr{
		IP: net.IP{127,0,0,1},
		Zone: "remote",
	}
	conn1,err1 := net.DialIP("tcp",localAddress, remoteAddress )
	if err1 != nil {
		fmt.Print("test tcp socket by DialIP failed! err is %s", err.Error())
	}

	localTcpAddress := &net.TCPAddr{
		IP: net.IP{127,0,0,1},
		Port: 55555,
		Zone: "local",
	}
	remoteTcpAddress := &net.TCPAddr{
		IP: net.IP{127,0,0,1},
		Port: 12345,
		Zone: "remote",
	}
	conn2,err2 := net.DialTCP("tcp", localTcpAddress, remoteTcpAddress)
	if err2 != nil {
		fmt.Print("test tcp socket by DialIP failed! err is %s", err.Error())
	}
}
//2.线程池
func testThreadPool(){
	//waitgroup使用

}

//3.grpc
//先定义protobuf

/* 9.文件读写 */
//文件处理的库，有os,io,ioutil,bufio,path/filepath
func testFileProcessByOS(fileName string){
	fileHandle,err := os.Create(fileName)
	buffer := []byte{}
	readCount1,osReadErr1 := fileHandle.Read(buffer)

	data := "test file process with reader by ReadFrom!!！"
	reader := strings.NewReader(data)
	readCount2,osReadErr2 := fileHandle.ReadFrom(reader)
	writeCount1,osWriteErr1 := fileHandle.Write()
	writeCount2,osWriteErr2 := fileHandle.WriteString("testFileProcess function test")
	fileHandle1,err1 := os.Open(fileName)
}

func testFileProcessByIO(filename string){
	fileHandle,err := os.Open(filename)
	writeCount, ioWriteErr := io.WriteString(fileHandle, "testFileProcessByIO")
	if ioWriteErr != nil {
		fmt.Println("WriteString failed, err is %s", ioWriteErr.Error())
	}

	data := "testFileProcesssByIO"
	reader := strings.NewReader(data)
	dataBuffer, ioReadErr := io.ReadAll(reader)
	if ioReadErr != nil{
		fmt.Printf("ReadAll failed, err is %s", ioReadErr.Error())
	}
}


/* 10.webserver  */
func testWebServer(){
	//addr := net.IPAddr{
	//	IP : []byte{127,0,0,1},
	//	Zone: "local",
	//}

	// 方式一：带handle
	//handle := mux.NewHandle()
	//handle.HandleFunc("/", testHandleRoot)
	//handle.HandleFunc("/hello", testHandlePath)
	handle := iris.New()
	rootParty := handle.Party("/")
	rootParty.HandleFunc("Post", "/", testHandleRoot)

	helloParty := handle.Party("/hello")
	helloParty.HandleFunc("Post", "/", testHandlePath)


	addr := "127.0.0.1:8080"
	if err := http.ListenAndServe(addr, handle);err != nil {
		fmt.Println("ListerAndServe failed, err is %s", err.Error())
		return
	}

	// 方式二：不带handle
	http.HandleFunc("/", testHandleRoot)
	http.HandleFunc("/hello", testHandlePath)

	port := ":8080"
	if err := http.ListenAndServe(port, nil);err != nil {
		fmt.Println("Listen Server failed, err is %s", err.Error())
		return
	}
}

func testHandleRoot(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w, "Welcome to Go http server ")
}

func testHandlePath(w http.ResponseWriter, r *http.Request){
	name := r.URL.Query().Get("name")
	if name == ""{
		name = "World"
	}

	fmt.Fprintln(w, "Welcome to Go http server, Hello, %s! \n", name)
}

/* 11.多协程处理  */
//1.goroutine使用，waitgroup使用
func testGoRoutine(){
	wg1 := waitgroup.SafeWaitGroup{}
	wg1.Add(1)
	wg1.Wait()
	wg1.Done()
	//waitgroup.RateLimiter()

	var wg sync.WaitGroup

	numWorks :=5
	wg.Add(numWorks)
	for i:=1;i<=5;i++{
		go func(i int)int{
			fmt.Println("starting func")
			time.Sleep(15*time.Second)
			wg.Done()
			fmt.Println("stoping func")
			return i*i
		}(1)

		//wg.Add(1)
	}

	wg.Wait()
	fmt.Println("All workers Done")
}

/* 12.互斥锁、读写锁等使用  */
func testLock(){
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	fileHandle,err := os.Open("/testfile")
	if err != nil {
		fmt.Println("Open file failed,err is %s", err.Error())
	}

	_,err1 := fileHandle.WriteString("hahahaha")
	if err1 == nil{
		return
	}

	//RWMutex
	var (
		rwLock = sync.RWMutex{}
		wg = sync.WaitGroup{}
		data = 10000
	)

	wg.Add(2)
	go func(){
		rwLock.RLock()
		defer rwLock.RUnlock()
		fmt.Println("Using RLock, data is ", data)
		wg.Done()
	}()

	go func(){
		rwLock.Lock()
		defer rwLock.Lock()
		fmt.Println("Using Lock, Before process, data  is ", data)
		data = 20
		fmt.Println("Using Lock, After processed, data  is ", data)
		wg.Done()
	}()

	wg.Wait()

	//sync.Cond，配合sync.Mutex使用
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	go func(){
		mu.Lock()
		defer mu.Unlock()
		cond.Wait()
		fmt.Println("Condition met")
	}()

	time.Sleep(10 *time.Second)
	mu.Lock()
	//cond.Signal()	//唤醒一个等待的协程
	cond.Broadcast()	//唤醒所有等待的协程
	mu.Unlock()

	//sync.Once，常用于初始化或单例场景
	once := sync.Once{}
	once.Do(func(){
		fmt.Println("this will be printed only once")
	})

	//sync.Atomic，原子操作，原子增加、原子读取、原子写入、原子交换
	//var counter uint32 = 0
	counter := uint32(18)
	atomic.AddUint32(&counter, 12)
	fmt.Println(atomic.LoadUint32(&counter))


	//sync.Map并发安全字典

	//sync.Pool对象池
}

/* 13.  */
func main() {

	print()
}
