package service

import (
	"context"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/siddontang/go-log/log"
	"os"
	"strconv"
	"velox/server/src/service/replic"
)

var hostReplic = "localhost"
var portReplic = "3308"
var hostMaser = "localhost"
var portMaster = "3307"

func replicate() {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "localhost",
		Port:     3307,
		User:     "root",
		Password: "123456",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// Start sync with specified binlog file and position
	streamer, _ := syncer.StartSync(mysql.Position{Name: "", Pos: 1})

	// or you can start a gtid replication like
	// streamer, _ := syncer.StartSyncGTID(gtidSet)
	// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
	// the mariadb GTID set likes this "0-1-100"

	for {
		ev, _ := streamer.GetEvent(context.Background())
		// Dump event
		ev.Dump(os.Stdout)
	}

}

func StatusMaster() (string, int64, error) {
	pool := client.NewPool(log.Debugf, 100, 400, 5, "localhost:3307", `root`, `123456`, `test`)
	// ...
	conn, _ := pool.GetConn(context.Background())
	defer pool.PutConn(conn)
	r, err := conn.Execute("SHOW MASTER STATUS")
	if err != nil {
		return "", 0, err
	}
	binFile, _ := r.GetString(0, 0)
	binPos, _ := r.GetInt(0, 1)
	return binFile, binPos, nil
}

func createUserForSlave2() error {
	db, err := sqlx.Open("mysql", "root:123456@("+hostMaser+":"+portMaster+")/")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	_, err = db.Exec("grant replication slave on *.* to 'replic'@'%' IDENTIFIED BY '123456';")
	return err
}

func registerSlave() error {
	conn, err := sqlx.Open("mysql", "root:123456@("+hostReplic+":"+portReplic+")/")
	_, _ = conn.Exec("stop slave;")
	_, err = conn.Exec("change master to master_host='" + hostMaser + "',master_port=" + portMaster + ",master_user='replic',master_password='123456',master_auto_position=1;")
	if err != nil {
		return err
	}
	_, err = conn.Exec("start slave;")
	return err
}

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Infof("%s %v\n", e.Action, e.Rows)
	return nil

}

func (h *MyEventHandler) OnTableChanged(schema string, table string) error {
	log.Infof("%s %v\n", schema, table)
	return nil
}
func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}
func StartCanal() {

	port, _ := strconv.Atoi(portMaster)
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     hostMaser,
		Port:     uint16(port),
		User:     "root",
		Password: "123456",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// Start sync with specified binlog file and position
	streamer, _ := syncer.StartSync(mysql.Position{"", 1})

	// or you can start a gtid replication like
	// streamer, _ := syncer.StartSyncGTID(gtidSet)
	// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
	// the mariadb GTID set likes this "0-1-100"

	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {

			println()
		}
		//Dump event
		ev.Dump(os.Stdout)
	}
}

func failoverCheck() (*replic.Server, *replic.Server, error) {
	f := replic.MysqlGTIDHandler{}

	master := replic.NewServer("localhost:3307", replic.User{
		Name:     "root",
		Password: "123456",
	}, replic.User{
		Name:     "replic",
		Password: "123456",
	})
	slave := replic.NewServer("localhost:3308", replic.User{
		Name:     "root",
		Password: "123456",
	}, replic.User{
		Name:     "replic",
		Password: "123456",
	})
	err := f.ChangeMasterTo(slave, master)
	return master, slave, err

}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}

type myMember string

func (m myMember) String() string {
	return string(m)
}
func Consistent() {
	cfg := consistent.Config{
		PartitionCount:    7,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}
	c := consistent.New(nil, cfg)

	// Add some members to the consistent hash table.
	// Add function calculates average load and distributes partitions over members
	node1 := myMember("nodgergrm")
	c.Add(node1)

	node2 := myMember("noergergom")
	c.Add(node2)
	key := []byte("my-key")
	owner := c.LocateKey(key)
	println(owner)
}
