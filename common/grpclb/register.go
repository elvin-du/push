package grpclb

import (
	"fmt"
	"strings"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
)

// Prefix should start and end with no slash
var Prefix = "push"
var client etcd3.Client
var serviceKey string

var stopSignal = make(chan bool, 1)

// Register
func Register(name string, host string, port int, target string, interval time.Duration, ttl int) error {
	serviceValue := fmt.Sprintf("%s:%d", host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceValue)

	// get endpoints for register dial address
	var err error
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		grpclog.Errorf("grpclb: create etcd3 client failed: %v", err)
		return fmt.Errorf("grpclb: create etcd3 client failed: %v", err)
	}

	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(interval)
		for {
			// minimum lease TTL is ttl-second
			resp, _ := client.Grant(context.TODO(), int64(ttl))
			// should get first, if not exist, set it
			_, err := client.Get(context.Background(), serviceKey)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := client.Put(context.TODO(), serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
						grpclog.Errorf("grpclb: set service '%s' with ttl to etcd3 failed: %s", name, err.Error())
					}
				} else {
					grpclog.Printf("grpclb: service '%s' connect to etcd3 failed: %s", name, err.Error())
				}
			} else {
				// refresh set to true for not notifying the watcher
				if _, err := client.Put(context.Background(), serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
					grpclog.Errorf("grpclb: refresh service '%s' with ttl to etcd3 failed: %s", name, err.Error())
				}
			}
			select {
			case <-stopSignal:
				return
			case <-ticker.C:
			}
		}
	}()

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister() error {
	stopSignal <- true
	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
	var err error
	if _, err := client.Delete(context.Background(), serviceKey); err != nil {
		grpclog.Errorf("grpclb: deregister '%s' failed: %s", serviceKey, err.Error())
	} else {
		grpclog.Printf("grpclb: deregister '%s' ok.", serviceKey)
	}
	return err
}
