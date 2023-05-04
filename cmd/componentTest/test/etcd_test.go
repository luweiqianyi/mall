package test

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

const (
	gK1 = "k1"
	gV1 = "v1"
	gV2 = "666"
)

// TestEtcdPutV1 测试etcd的某个key放入值
func TestEtcdPutV1(t *testing.T) {
	cli := newEtcdClient()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	resp, err := cli.Put(ctx, gK1, gV1)
	cancel()
	if err != nil {
		log.Printf("error: %v", err)
	}
	// use the response:resp:{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:1041 raft_term:29  <nil> {} [] 0}
	log.Printf("resp:%v", *resp)
}

// TestEtcdPutV1 测试etcd往同一个key放入新的值
func TestEtcdPutV2(t *testing.T) {
	cli := newEtcdClient()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	resp, err := cli.Put(ctx, gK1, gV2)
	cancel()
	if err != nil {
		log.Printf("error: %v", err)
	}
	// use the response:resp:{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:1041 raft_term:29  <nil> {} [] 0}
	log.Printf("resp:%v", *resp)
}

// TestEtcdDeleteK1 删除etcd中某个key
func TestEtcdDeleteK1(t *testing.T) {
	cli := newEtcdClient()
	defer cli.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	resp, err := cli.Delete(ctx, gK1)
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("resp:%v", *resp)
}

// TestEtcdGet 测试etcd获取某个key的值
func TestEtcdGet(t *testing.T) {
	cli := newEtcdClient()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	resp, err := cli.Get(ctx, gK1)
	cancel()
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("resp:%v", *resp)
}

// TestEtcdWatch 监听etcd中某个key对应的值的变化
func TestEtcdWatch(t *testing.T) {
	cli := newEtcdClient()
	defer cli.Close()

	ctx, _ := context.WithCancel(context.Background())
	watcher := clientv3.NewWatcher(cli)
	ch := watcher.Watch(ctx, gK1)
	defer watcher.Close()

	for {
		resp, ok := <-ch

		if ok {
			log.Printf("receive data: %v %v\n", resp, resp.Events)
		} else {
			log.Fatalf("cannot receive data,err:%v %v %v\n", resp.Canceled, resp.Err(), resp.Header)
			break
		}

		if resp.Canceled {
			log.Fatalf("watch cancel,resp=%v", resp)
			break
		}
		for _, event := range resp.Events {
			log.Printf("%v", event)
		}
	}
}

// newEtcdClient 创建与etcd服务的客户端连接对象，通过该对象操作远程etcd服务
func newEtcdClient() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, //"localhost:2379", "localhost:22379", "localhost:32379"
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("new etcd client error: %v", err)
	}
	return cli
}
