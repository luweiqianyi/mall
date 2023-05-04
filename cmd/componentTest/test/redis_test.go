package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	RedisHost     = "localhost:6379"
	Password      = ""
	gRedisK1      = "k1"
	gRedisK2      = "k2"
	gRedisV1      = "v1"
	gRedisV2      = "v2"
	gRedisHashKey = "hashKey"
)

// HSet: Time complexity:O(1) for each field/value pair added, so O(N) to add N field/value pairs when the command is called with multiple field/value pairs.
// HGet: Time complexity:O(1)
// HDel: Time complexity:O(N) where N is the number of fields to be removed.
// HExists: Time complexity: O(1)

// go get github.com/go-redis/redis/v8
// TestRedisSet 测试Redis的Set功能,Redis采用覆盖写的方式
func TestRedisSet(t *testing.T) {
	rdb := newRedisClient()
	err := rdb.Set(context.Background(), gRedisK1, gRedisV2, 0).Err()
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
}

// TestRedisGet 测试Redis的Get功能
func TestRedisGet(t *testing.T) {
	rdb := newRedisClient()
	v, err := rdb.Get(context.Background(), gRedisK1).Result()
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	log.Printf("%s:%s", gRedisK1, v)
}

func TestRedisDel(t *testing.T) {
	rdb := newRedisClient()
	err := rdb.Del(context.Background(), gRedisK1).Err()
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	log.Printf("delete key:%s success", gRedisK1)
}

// TestRedisGetKeyNotExist 测试Redis key不存在时的Get功能
func TestRedisGetKeyNotExist(t *testing.T) {
	rdb := newRedisClient()
	_, err := rdb.Get(context.Background(), gRedisK2).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("key not exist!!!\n")
		}
		log.Printf("%v", err)
		panic(err)
	}
}

// TestSetGetConcurrently 测试Redis 多用户并发Set Get同一key的Get功能
func TestSetGetConcurrently(t *testing.T) {
	rdb := newRedisClient()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		rdb.Set(context.Background(), gRedisK1, gRedisV1, 0)
		log.Printf("go routine1: set key[%s]=[%s]\n", gRedisK1, gRedisV1)
		time.Sleep(time.Second * 3) // 等3s后再获取值，此时另外一个协程会将gRedisK1更新为它所想要的值，会覆盖掉本协程赋予的值
		v, err := rdb.Get(context.Background(), gRedisK1).Result()
		if err != nil {
			log.Printf("err:%v\n", err)
			return
		}
		log.Printf("go routine1 get, %s:%s", gRedisK1, v)
	}()

	go func() {
		defer wg.Done()
		rdb.Set(context.Background(), gRedisK1, gRedisV2, 0)
		log.Printf("go routine1: set key[%s]=[%s]\n", gRedisK1, gRedisV2)
		v, err := rdb.Get(context.Background(), gRedisK1).Result()
		if err != nil {
			log.Printf("err:%v\n", err)
			return
		}
		log.Printf("go routine1 get, %s:%s", gRedisK1, v)
	}()

	wg.Wait()

	//=== RUN   TestSetGetConcurrently
	//2023/04/24 12:03:15 go routine1: set key[k1]=[v1]
	//2023/04/24 12:03:15 go routine1: set key[k1]=[v2]
	//2023/04/24 12:03:15 go routine1 get, k1:v2
	//2023/04/24 12:03:18 go routine1 get, k1:v2
	//--- PASS: TestSetGetConcurrently (3.03s)
	//PASS
	// 如果将协程1和协程2看成两个用户1和2对同一项目进行更新的话，用户1将k1设置为v1,但是得到的结果确实v2，站在用户1的角度，是不想看到的
	// 因为用户1会想我明明已经更新了为v1,为什么看到的值是v2
}

type Client struct {
	db *redis.Client
}

func (cli *Client) writeDataToRedisUsingHSet(key string, k, v string) {
	err := cli.db.HSet(context.Background(), key, k, v).Err()
	if err != nil {
		log.Printf("save %s %s:%s err=%v", key, k, v, err)
	}
}

// TestRedisHSet 测试Redis的HSet功能
func TestRedisHSet(t *testing.T) {
	rdb := newRedisClient()
	cli := &Client{
		db: rdb,
	}

	// 设计1：用户ID作为field
	appID := "app1"
	appRoomID := "room1"
	key := fmt.Sprintf("wplive:userroom:service1:dc:dc1:app:%s:appRoom:%s:users", appID, appRoomID)
	session := map[string]string{"user1": "{}", "user2": "{}", "user3": "{}", "user4": "{}"}
	for k, v := range session {
		cli.writeDataToRedisUsingHSet(key, k, v)
	}

	// 单个用户数据
	appUser1 := "u1"
	appUser2 := "u2"
	key = fmt.Sprintf("wplive:userroom:service1:dc:dc1:app:%s:appRoom:%s:appUser:%s", appID, appRoomID, appUser1)
	cli.writeDataToRedisUsingHSet(key, "name", "zhangSan")
	cli.writeDataToRedisUsingHSet(key, "privilege", "1,2,3")
	key = fmt.Sprintf("wplive:userroom:service1:dc:dc1:app:%s:appRoom:%s:appUser:%s", appID, appRoomID, appUser2)
	cli.writeDataToRedisUsingHSet(key, "name", "leeSi")
	cli.writeDataToRedisUsingHSet(key, "privilege", "1,3")

}

type RedisUserLoginTime struct {
	LoginNodeID  string `json:"loginNodeID"`
	CreateTimeMs int64  `json:"createTime"`
}

type RedisStreamInfo struct {
	StreamID string
	GroupID  string

	AppID     string
	AppRoomID string
	AppUserID string

	EdgeEmitterID string
	EdgeRoomID    string
	EdgeUserID    string

	EdgePublishID string
	Bitrate       int
	Audio         string
	Video         string
	StreamName    string
	SelectPair    string

	HandlerServiceID string
	ReceivedMsgID    string
	CreateTimeMs     int64
}

// TestCurrentUserRoomDesign 测试Redis 多用户并发Set Get同一key的Get功能
func TestCurrentUserRoomDesign(t *testing.T) {
	rdb := newRedisClient()
	cli := &Client{
		db: rdb,
	}

	RedisHeader := "wplive:userroom:service1:dc:dc1"
	// 设计1：用户ID作为field
	appID := "app1"
	appRoom1 := "room1"
	appRoom2 := "room2"
	key := fmt.Sprintf("%s:app:%s:appRooms", RedisHeader, appID)
	cli.writeDataToRedisUsingHSet(key, appRoom1, "{}")
	cli.writeDataToRedisUsingHSet(key, appRoom2, "{}")

	appUser1 := "u1"
	appUser2 := "u2"
	key = fmt.Sprintf("%s:app:%s:appRoom:%s:appUsers", RedisHeader, appID, appRoom1)

	data, _ := json.Marshal(
		&RedisUserLoginTime{
			LoginNodeID:  "service1",
			CreateTimeMs: time.Now().UnixMilli(),
		})
	cli.writeDataToRedisUsingHSet(key, appUser1, string(data))
	cli.writeDataToRedisUsingHSet(key, appUser2, string(data))

	key = fmt.Sprintf("%s:app:%s:appRoom:%s:appUser:%s", RedisHeader, appID, appRoom1, appUser1)
	cli.writeDataToRedisUsingHSet(key, "name", "ZhangSan")
	cli.writeDataToRedisUsingHSet(key, "role", "teacher")

	key = fmt.Sprintf("%s:app:%s:appRoom:%s:appUser:%s", RedisHeader, appID, appRoom1, appUser2)
	cli.writeDataToRedisUsingHSet(key, "name", "LiSi")
	cli.writeDataToRedisUsingHSet(key, "role", "student")

	group1 := "group1"
	stream1 := "stream1"
	stream2 := "stream2"
	// 哪个group中有哪些stream
	key = fmt.Sprintf("%s:app:%s:appRoom:%s:group:%s", RedisHeader, appID, appRoom1, group1)
	streamInfo1 := RedisStreamInfo{
		StreamID: stream1,
		GroupID:  group1,

		AppID:     appID,
		AppRoomID: appRoom1,
		AppUserID: appUser1,

		EdgeEmitterID: "service2",
		EdgeRoomID:    "edgeRoom1",
		EdgeUserID:    "edgeUser1",

		EdgePublishID: "edgePublishID1",
		Bitrate:       60,
		Audio:         "audio",
		Video:         "video",
		StreamName:    fmt.Sprintf("%s's stream", appUser1),
		SelectPair:    "???",

		HandlerServiceID: "service1",
		ReceivedMsgID:    "msg1",
		CreateTimeMs:     time.Now().UnixMilli(),
	}
	streamInfoData1, _ := json.Marshal(streamInfo1)
	cli.writeDataToRedisUsingHSet(key, stream1, string(streamInfoData1))
	streamInfo2 := RedisStreamInfo{
		StreamID: stream1,
		GroupID:  group1,

		AppID:     appID,
		AppRoomID: appRoom1,
		AppUserID: appUser2,

		EdgeEmitterID: "service2",
		EdgeRoomID:    "edgeRoom1",
		EdgeUserID:    "edgeUser1",

		EdgePublishID: "edgePublishID1",
		Bitrate:       60,
		Audio:         "audio",
		Video:         "video",
		StreamName:    fmt.Sprintf("%s's stream", appUser2),
		SelectPair:    "???",

		HandlerServiceID: "service1",
		ReceivedMsgID:    "msg2",
		CreateTimeMs:     time.Now().UnixMilli(),
	}
	streamInfoData2, _ := json.Marshal(streamInfo2)
	cli.writeDataToRedisUsingHSet(key, stream2, string(streamInfoData2))

	// 不区分stream在哪个group中,即room中的所有stream
	key = fmt.Sprintf("%s:app:%s:appRoom:%s:streams", RedisHeader, appID, appRoom1)
	cli.writeDataToRedisUsingHSet(key, stream1, string(streamInfoData1))
	cli.writeDataToRedisUsingHSet(key, stream2, string(streamInfoData2))

	// 发起的点名情况
	rollCall1 := "rollCall1"
	rollCall2 := "rollCall2"
	key = fmt.Sprintf("%s:app:%s:appRoom:%s:rollCalls", RedisHeader, appID, appRoom1)
	cli.writeDataToRedisUsingHSet(key, "latestRollCall", "1") // 最近一次点名
	cli.writeDataToRedisUsingHSet(key, rollCall1, "{}")
	cli.writeDataToRedisUsingHSet(key, rollCall2, "{}")

	// 记录参与点名的用户名单
	key = fmt.Sprintf("%s:app:%s:appRoom:%s:rollCall:%s", RedisHeader, appID, appRoom1, rollCall1)
	cli.writeDataToRedisUsingHSet(key, appUser1, "{information}")
	cli.writeDataToRedisUsingHSet(key, appUser2, "{information}")

}

// 创建Redis客户端对象，用于和Redis服务端进行通信
func newRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: Password,
		DB:       1,
	})
	return rdb
}
