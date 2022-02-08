package xcanal

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"math/rand"
	"time"
)

type (
	SyncerConfig struct {
		Stop        bool   `json:"Stop" yaml:"Stop"`
		Host        string `json:"Host" yaml:"Host"`
		User        string `json:"User" yaml:"User"`
		Password    string `json:"Password" yaml:"Password"`
		ReadTimeout int64  `json:"ReadTimeout" yaml:"ReadTimeout"` //1s

	}
)

//
func Replication(sfg SyncerConfig, spos *Position, rowhandler RowDataHandler) error {
	if sfg.Stop {
		return nil
	}
	cfg := canal.NewDefaultConfig()
	//root:vncweA8XYR60@tcp(rm-wz9i7l7xhp5u02w21.mysql.rds.aliyuncs.com:3306)/asset_v2?charset=utf8
	cfg.Addr = sfg.Host
	cfg.User = sfg.User
	cfg.Password = sfg.Password
	cfg.ServerID = uint32(RandInt(1000000, 99999999))
	if sfg.ReadTimeout > 0 {
		cfg.ReadTimeout = time.Duration(sfg.ReadTimeout) * time.Second
	}
	cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)
	if err != nil {
		return err
	}
	defer c.Close()
	//log.SetLevel(log.LevelError)
	// Register a handler to handle RowsEvent
	dh := &DefaultEventHandler{fn: rowhandler}

	c.SetEventHandler(dh)
	if spos == nil {
		p, err := c.GetMasterPos()
		if err != nil {
			return err
		}
		dh.binlogName = p.Name
		return c.Run()
	}
	err = checkPosition(c, spos)
	if err == nil {
		dh.binlogName = spos.Name
		return c.RunFrom(mysql.Position{
			Name: spos.Name,
			Pos:  uint32(spos.Pos),
		})
	}
	mp, err := c.GetMasterPos()
	if err != nil {
		return err
	}
	dh.binlogName = mp.Name
	return c.RunFrom(mysql.Position{
		Name: mp.Name,
		Pos:  mp.Pos,
	})
}

//RandInt RandInt 伪随机数 [min,max)
func RandInt(min, max int) int {
	if min >= max || (min == 0 && max == 0) {
		return max
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := r.Intn(max)
	if ret < min {
		return ret + min
	}
	return ret
}

//CheckPosition 检查日志是否存在　否则返回当前最新最新位置
func checkPosition(c *canal.Canal, pos *Position) error {
	if pos.Name == "" {
		return xtools.XErr(codes.InvalidArgument, "pos.name is nil")
	}
	rr, err := c.Execute("SHOW MASTER LOGS;")
	if err != nil {
		return err
	}
	rownum := rr.RowNumber()

	for index := 0; index < rownum; index++ {
		val, err := rr.GetStringByName(index, "Log_name")
		if err != nil {
			continue
		}
		if val == pos.Name {
			return nil
		}
	}
	return xtools.XErr(codes.InvalidArgument, "pos. is err")
}
func getFirstPostion(c *canal.Canal) (*mysql.Position, error) {
	rr, err := c.Execute("SHOW MASTER LOGS;")
	if err != nil {
		return nil, xtools.XErr(codes.Internal, err.Error(), true)
	}
	rownum := rr.RowNumber()
	if rownum < 1 {
		return nil, xtools.XErr(codes.Internal, "empty is logs", true)
	}
	pos := &mysql.Position{}
	pos.Name, _ = rr.GetStringByName(0, "Log_name")
	pos.Pos = 0
	return pos, nil
}
