package xid

import "github.com/rs/xid"

func GetXid() string {
	id := xid.New()
	xidStr := id.String()
	return xidStr[10:] //默认返回太长，截取掉一部分
}
