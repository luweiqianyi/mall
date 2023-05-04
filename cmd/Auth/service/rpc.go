package service

import (
	"context"
	"mall/cmd/Auth/config"
	"mall/cmd/Auth/pb"
	"mall/pkg/log"
	"mall/pkg/util"
	"time"
)

const (
	QueryTokenCodeSuccess = 1001
	QueryTokenCodeFail    = 1002
	SaveTokenCodeSuccess  = 1003
	SaveTokenCodeFail     = 1004
	DelTokenSuccess       = 1005
	DelTokenFail          = 1006
	RedisDisconnected     = 1007

	Success = "success"
	Fail    = "fail"
)

type AuthServiceServerImpl struct {
}

func (l AuthServiceServerImpl) QueryToken(ctx context.Context, request *pb.QueryTokenRequest) (*pb.QueryTokenResponse, error) {
	if config.GetGConfig() == nil || config.GetGConfig().GetRedisClient() == nil {
		log.PrintLog("redis connection invalid,remote address:%s", config.GetGConfig().GetRedisHost())
		return &pb.QueryTokenResponse{
			Code: RedisDisconnected,
			Msg:  Fail,
		}, nil
	}
	redisClient := config.GetGConfig().GetRedisClient()

	accountName := request.AccountName

	token, err := util.RedisGet(redisClient, accountName)
	if err != nil {
		log.PrintLog("get key[%s]'s value failed,err=%v", accountName, err)
		return &pb.QueryTokenResponse{
			Code: QueryTokenCodeFail,
			Msg:  Fail,
		}, nil
	}
	return &pb.QueryTokenResponse{
		Code:  QueryTokenCodeSuccess,
		Msg:   Success,
		Token: token,
	}, nil
}

func (l AuthServiceServerImpl) SaveToken(ctx context.Context, request *pb.SaveTokenRequest) (*pb.SaveTokenResponse, error) {
	if config.GetGConfig() == nil || config.GetGConfig().GetRedisClient() == nil {
		log.PrintLog("redis connection invalid,remote address:%s", config.GetGConfig().GetRedisHost())
		return &pb.SaveTokenResponse{
			Code: RedisDisconnected,
			Msg:  Fail,
		}, nil
	}
	redisClient := config.GetGConfig().GetRedisClient()

	accountName := request.AccountName
	token := request.Token
	tokenExpireTime := request.TokenExpireTime

	err := util.RedisSet(redisClient, accountName, token, time.Duration(tokenExpireTime))
	if err != nil {
		log.PrintLog("save key:value [%s]:[%s] failed,err=%v", accountName, token, err)
		return &pb.SaveTokenResponse{
			Code: SaveTokenCodeFail,
			Msg:  Fail,
		}, nil
	}
	return &pb.SaveTokenResponse{
		Code: SaveTokenCodeSuccess,
		Msg:  Success,
	}, nil
}

func (l AuthServiceServerImpl) DelToken(ctx context.Context, request *pb.DelTokenRequest) (*pb.DelTokenResponse, error) {
	if config.GetGConfig() == nil || config.GetGConfig().GetRedisClient() == nil {
		log.PrintLog("redis connection invalid,remote address:%s", config.GetGConfig().GetRedisHost())
		return &pb.DelTokenResponse{
			Code: RedisDisconnected,
			Msg:  Fail,
		}, nil
	}
	redisClient := config.GetGConfig().GetRedisClient()

	accountName := request.AccountName

	_, err := util.RedisGet(redisClient, accountName)
	if err != nil {
		return &pb.DelTokenResponse{
			Code: DelTokenFail,
			Msg:  Fail,
		}, nil
	}

	err = util.RedisDel(redisClient, accountName)
	if err != nil {
		log.PrintLog("del key[%s]failed,err=%v", accountName, err)
		return &pb.DelTokenResponse{
			Code: DelTokenFail,
			Msg:  Fail,
		}, nil
	}

	return &pb.DelTokenResponse{
		Code: DelTokenSuccess,
		Msg:  Success,
	}, nil
}
