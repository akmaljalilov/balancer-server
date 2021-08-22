package service

import (
	_ "github.com/go-sql-driver/mysql"
	"velox/server/src/service/replic"
)

var MapServer = map[string]*replic.Server{
	"localhost:3310": replic.NewServer("localhost:3310", &replic.User{
		Name:     "root",
		Password: "123456",
	}, &replic.User{
		Name:     "abc",
		Password: "123456",
	}),
	"localhost:3307": replic.NewServer("localhost:3307", &replic.User{
		Name:     "root",
		Password: "123456",
	}, nil),
	"localhost:3308": replic.NewServer("localhost:3308", &replic.User{
		Name:     "root",
		Password: "123456",
	}, &replic.User{
		Name:     "abc",
		Password: "123456",
	}),
	"localhost:3309": replic.NewServer("localhost:3309", &replic.User{
		Name:     "root",
		Password: "123456",
	}, nil),
}

func ConnectAndGetToServer() (*replic.Server, *replic.Server, *replic.Server, *replic.Server, error) {
	f := replic.MysqlGTIDHandler{}

	master := MapServer["localhost:3310"]
	slave := MapServer["localhost:3307"]
	err := f.ChangeMasterTo(slave, master, "mysql-master")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	master2 := MapServer["localhost:3308"]
	slave2 := MapServer["localhost:3309"]

	err = f.ChangeMasterTo(slave2, master2, "mysql-master2")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return master, slave, master2, slave2, err

}
