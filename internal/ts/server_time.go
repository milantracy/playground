package ts

/**
 * The handler for TimeServer RPC service.
 */

import (
	"time"
)

type TimeServerArgs struct{}

type TimeServer struct{}

func (s *TimeServer) ServerTime(args *TimeServerArgs, reply *int64) error {
	*reply = time.Now().Unix()
	return nil
}
