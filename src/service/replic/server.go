package replic

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ResultExec struct {
}

type User struct {
	Name     string
	Password string
}

func (u User) FullName() string {
	return u.Name + ":" + u.Password
}

type Server struct {
	Addr string

	User     *User
	ReplUser *User

	conn *sqlx.DB
}

func NewServer(addr string, user *User, replUser *User) *Server {
	s := new(Server)

	s.Addr = addr

	s.User = user
	s.ReplUser = replUser
	if replUser != nil {
		_, err := s.Execute("CREATE USER " + replUser.Name + "@'%' IDENTIFIED BY '" + replUser.Password + "';")
		if err != nil {
			fmt.Println("Error in create user")
		}
		_, err = s.Execute("GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO " + replUser.Name + "@'%';")
		if err != nil {
			fmt.Println("Error in grant")
		}
	}

	return s
}

func (s *Server) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Server) Execute(cmd string, args ...interface{}) (map[string]sql.NullString, error) {
	var err error
	if s.conn == nil {
		s.conn, err = sqlx.Open("mysql", s.User.FullName()+"@("+s.Addr+")/")
		if err != nil {
			return nil, err
		}
	}

	r, err := s.conn.Query(cmd)
	if err != nil {
		return nil, err
	}
	return scanMap(r)
}

func (s *Server) StartSlave() error {
	_, err := s.Execute("START SLAVE")
	return err
}

func (s *Server) StopSlave() error {
	_, err := s.Execute("STOP SLAVE")
	return err
}

func (s *Server) StopSlaveIOThread() error {
	_, err := s.Execute("STOP SLAVE IO_THREAD")
	return err
}

func (s *Server) SlaveStatus() (map[string]sql.NullString, error) {
	res, err := s.Execute("SHOW SLAVE STATUS")
	return res, err
}
func scanMap(rows *sql.Rows) (map[string]sql.NullString, error) {
	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	}

	values := make([]interface{}, len(columns))

	for index := range values {
		values[index] = new(sql.NullString)
	}

	err = rows.Scan(values...)

	if err != nil {
		return nil, err
	}

	result := make(map[string]sql.NullString)

	for index, columnName := range columns {
		result[columnName] = *values[index].(*sql.NullString)
	}

	return result, nil
}

func (s *Server) MasterStatus() (*ResultExec, error) {
	_, err := s.Execute("SHOW MASTER STATUS")
	if err != nil {
		return nil, err
	} else {
		//return r, nil
	}
	return nil, err
}

func (s *Server) ResetSlave() error {
	_, err := s.Execute("RESET SLAVE")
	return err
}

func (s *Server) ResetSlaveALL() error {
	_, err := s.Execute("RESET SLAVE ALL")
	return err
}

func (s *Server) ResetMaster() error {
	_, err := s.Execute("RESET MASTER")
	return err
}

// TODO
//func (s *Server) MysqlGTIDMode() (string, error) {
//	r, err := s.Execute("SELECT @@gtid_mode")
//	if err != nil {
//		return GTIDModeOff, err
//	}
//	on, _ := r.GetString(0, 0)
//	if on != GTIDModeOn {
//		return GTIDModeOff, nil
//	} else {
//		return GTIDModeOn, nil
//	}
//}

func (s *Server) SetReadonly(b bool) error {
	var err error
	if b {
		_, err = s.Execute("SET GLOBAL read_only = ON")
	} else {
		_, err = s.Execute("SET GLOBAL read_only = OFF")
	}
	return err
}

func (s *Server) LockTables() error {
	_, err := s.Execute("FLUSH TABLES WITH READ LOCK")
	return err
}

func (s *Server) UnlockTables() error {
	_, err := s.Execute("UNLOCK TABLES")
	return err
}

//// FetchSlaveReadPos gets current binlog filename and position read from master
//func (s *Server) FetchSlaveReadPos() (Position, error) {
//	r, err := s.SlaveStatus()
//	if err != nil {
//		return Position{}, err
//	}
//
//	fname, _ := r.GetStringByName(0, "Master_Log_File")
//	pos, _ := r.GetIntByName(0, "Read_Master_Log_Pos")
//
//	return Position{Name: fname, Pos: uint32(pos)}, nil
//}

//// FetchSlaveExecutePos gets current executed binlog filename and position from master
//func (s *Server) FetchSlaveExecutePos() (Position, error) {
//	r, err := s.SlaveStatus()
//	if err != nil {
//		return Position{}, err
//	}
//
//	fname, _ := r.GetStringByName(0, "Relay_Master_Log_File")
//	pos, _ := r.GetIntByName(0, "Exec_Master_Log_Pos")
//
//	return Position{Name: fname, Pos: uint32(pos)}, nil
//}

//func (s *Server) MasterPosWait(pos Position, timeout int) error {
//	_, err := s.Execute(fmt.Sprintf("SELECT MASTER_POS_WAIT('%s', %d, %d)", pos.Name, pos.Pos, timeout))
//	return err
//}
