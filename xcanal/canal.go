package xcanal

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"time"
)

//
func DumpAndReplication(opt DumpConfig, rowhandler RowDataHandler) error {
	if opt.Stop {
		return nil
	}
	cfg := canal.NewDefaultConfig()
	//root:vncweA8XYR60@tcp(rm-wz9i7l7xhp5u02w21.mysql.rds.aliyuncs.com:3306)/asset_v2?charset=utf8
	cfg.Addr = opt.Host
	cfg.User = opt.User
	cfg.Password = opt.Password
	cfg.ServerID = uint32(RandInt(1000000, 99999999))
	if opt.ReadTimeout > 0 {
		cfg.ReadTimeout = time.Duration(opt.ReadTimeout) * time.Second
	}

	cfg.Dump.ExecutionPath = opt.ExecutionPath
	if cfg.Dump.ExecutionPath == "" {
		cfg.Dump.ExecutionPath = "mysqldump"
	}
	if opt.MaxAllowedPacket > 0 {
		cfg.Dump.MaxAllowedPacketMB = opt.MaxAllowedPacket
	}
	if len(opt.Database) > 0 {
		if len(opt.Database) == 1 && len(opt.Table) > 0 {
			//cfg.Dump.Databases = opt.Database
			cfg.Dump.TableDB = opt.Database[0]
			cfg.Dump.Tables = opt.Table
		} else {
			cfg.Dump.Databases = opt.Database
		}
	}
	if opt.SkipMasterData {
		cfg.Dump.SkipMasterData = true
	}
	if opt.Charset != "" {

	}
	if len(opt.ExtraOptions) > 0 {
		//d.SetExtraOptions(opt.ExtraOptions)
		cfg.Dump.ExtraOptions = opt.ExtraOptions
	}
	if opt.Where != "" {
		cfg.Dump.Where = opt.Where
	}
	if opt.Charset != "" {
		cfg.Charset = opt.Charset
	}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		return err
	}

	defer c.Close()
	// Register a handler to handle RowsEvent
	c.SetEventHandler(&DefaultEventHandler{fn: rowhandler})

	// Start canal
	return c.Run()
}
