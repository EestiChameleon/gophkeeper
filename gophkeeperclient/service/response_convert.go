package service

import (
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

// VaultSyncConvert convert gRPC response proto data (slices) to local data format (map).
func VaultSyncConvert(in *pb.SyncVaultResponse) *models.Vault {
	return &models.Vault{
		Pair: responsePairArrayToMap(in.Pairs),
		Text: responseTextArrayToMap(in.Texts),
		Bin:  responseBinArrayToMap(in.BinData),
		Card: responseCardArrayToMap(in.Cards),
	}
}

// responsePairArrayToMap converts pair slice to pair map.
func responsePairArrayToMap(pairs []*pb.Pair) map[string]*models.Pair {
	result := make(map[string]*models.Pair)
	if len(pairs) < 1 {
		return result
	}
	for _, v := range pairs {
		result[v.Title] = models.ProtoToModelsPair(v)
	}

	return result
}

// responseTextArrayToMap converts text slice to pair map.
func responseTextArrayToMap(texts []*pb.Text) map[string]*models.Text {
	result := make(map[string]*models.Text)
	if len(texts) < 1 {
		return result
	}
	for _, v := range texts {
		result[v.Title] = models.ProtoToModelsText(v)
	}

	return result
}

// responseBinArrayToMap converts bin slice to bin map.
func responseBinArrayToMap(bins []*pb.Bin) map[string]*models.Bin {
	result := make(map[string]*models.Bin)
	if len(bins) < 1 {
		return result
	}
	for _, v := range bins {
		result[v.Title] = models.ProtoToModelsBin(v)
	}

	return result
}

// responseCardArrayToMap converts card slice to card map.
func responseCardArrayToMap(cards []*pb.Card) map[string]*models.Card {
	result := make(map[string]*models.Card)
	if len(cards) < 1 {
		return result
	}
	for _, v := range cards {
		result[v.Title] = models.ProtoToModelsCard(v)
	}

	return result
}
