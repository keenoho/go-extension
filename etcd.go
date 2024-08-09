package extension

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	cli            *clientv3.Client
	contextTimeout time.Duration
}

func (m *Etcd) Init() error {
	if m.cli != nil {
		return nil
	}

	endpoints := strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
	dialTimeoutInt, _ := strconv.Atoi(os.Getenv("ETCD_DIAL_TIMEOUT"))
	if dialTimeoutInt < 1 {
		dialTimeoutInt = 1000
	}
	dialTimeout := time.Duration(dialTimeoutInt) * time.Millisecond
	contextTimeoutInt, _ := strconv.Atoi(os.Getenv("ETCD_CONTEXT_TIMEOUT"))
	if contextTimeoutInt < 1 {
		contextTimeoutInt = 500
	}
	contextTimeout := time.Duration(contextTimeoutInt) * time.Millisecond

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	m.cli = cli
	m.contextTimeout = contextTimeout
	return m.testConnect()
}

func (m *Etcd) Etcd() *clientv3.Client {
	return m.cli
}

func (m *Etcd) EtcdGet(key string, options ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.contextTimeout)
	res, err := m.cli.Get(ctx, key, options...)
	cancel()
	return res, err
}

func (m *Etcd) EtcdPut(key string, value string, options ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.contextTimeout)
	res, err := m.cli.Put(ctx, key, value, options...)
	cancel()
	return res, err
}

func (m *Etcd) EtcdDelete(key string, options ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.contextTimeout)
	res, err := m.cli.Delete(ctx, key, options...)
	cancel()
	return res, err
}

func (m *Etcd) EtcdWatch(key string, callback func(*clientv3.Event), options ...clientv3.OpOption) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	res := m.cli.Watch(ctx, key, options...)
	for re := range res {
		for _, ev := range re.Events {
			callback(ev)
		}
	}
	return cancel
}

func (m *Etcd) testConnect() error {
	var hasError error
	endPoints := m.cli.Endpoints()

	for _, ep := range endPoints {
		ctx, cancel := context.WithTimeout(context.Background(), m.contextTimeout)
		res, err := m.cli.Status(ctx, ep)
		cancel()
		if err != nil {
			hasError = err
			break
		}
		log.Println("etcd endPoint status:", ep, res.Header.ClusterId, res.Header.MemberId)
	}
	return hasError
}
