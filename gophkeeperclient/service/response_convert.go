package service

import (
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

// VaultSyncConvert convert gRPC response data (slices) to local data format (map).
func VaultSyncConvert(in *pb.SyncVaultResponse) *models.VaultProto {
	return &models.VaultProto{
		Pair: responsePairArrayToMap(in.Pairs),
		Text: responseTextArrayToMap(in.Texts),
		Bin:  responseBinArrayToMap(in.BinData),
		Card: responseCardArrayToMap(in.Cards),
	}
}

// responsePairArrayToMap converts pair slice to pair map.
func responsePairArrayToMap(pairs []*pb.Pair) map[string]*pb.Pair {
	result := make(map[string]*pb.Pair)
	for _, v := range pairs {
		result[v.Title] = v
	}

	return result
}

// responseTextArrayToMap converts text slice to pair map.
func responseTextArrayToMap(texts []*pb.Text) map[string]*pb.Text {
	result := make(map[string]*pb.Text)
	for _, v := range texts {
		result[v.Title] = v
	}

	return result
}

// responseBinArrayToMap converts bin slice to bin map.
func responseBinArrayToMap(bins []*pb.Bin) map[string]*pb.Bin {
	result := make(map[string]*pb.Bin)
	for _, v := range bins {
		result[v.Title] = v
	}

	return result
}

// responseCardArrayToMap converts card slice to card map.
func responseCardArrayToMap(cards []*pb.Card) map[string]*pb.Card {
	result := make(map[string]*pb.Card)
	for _, v := range cards {
		result[v.Title] = v
	}

	return result
}
