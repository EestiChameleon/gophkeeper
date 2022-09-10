package service

import (
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

func VaultSyncConvert(in *pb.SyncVaultResponse) *models.VaultProto {
	return &models.VaultProto{
		Pair: responsePairArrayToMap(in.Pairs),
		Text: responseTextArrayToMap(in.Texts),
		Bin:  responseBinArrayToMap(in.BinData),
		Card: responseCardArrayToMap(in.Cards),
	}
}

func responsePairArrayToMap(pairs []*pb.Pair) map[string]*pb.Pair {
	result := make(map[string]*pb.Pair)
	for _, v := range pairs {
		result[v.Title] = v
	}

	return result
}

func responseTextArrayToMap(texts []*pb.Text) map[string]*pb.Text {
	result := make(map[string]*pb.Text)
	for _, v := range texts {
		result[v.Title] = v
	}

	return result
}

func responseBinArrayToMap(bins []*pb.Bin) map[string]*pb.Bin {
	result := make(map[string]*pb.Bin)
	for _, v := range bins {
		result[v.Title] = v
	}

	return result
}

func responseCardArrayToMap(cards []*pb.Card) map[string]*pb.Card {
	result := make(map[string]*pb.Card)
	for _, v := range cards {
		result[v.Title] = v
	}

	return result
}
