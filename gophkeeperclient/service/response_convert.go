package service

import (
	"database/sql"
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
		result[v.Title] = ProtoToModelsPair(v)
	}

	return result
}

// ProtoToModelsPair converts proto Pair data to local Pair.
func ProtoToModelsPair(p *pb.Pair) *models.Pair {
	return &models.Pair{
		ID:        0,
		UserID:    0,
		Title:     p.GetTitle(),
		Login:     p.GetLogin(),
		Pass:      p.GetPass(),
		Comment:   p.GetComment(),
		Version:   p.GetVersion(),
		DeletedAt: sql.NullTime{},
	}
}

// responseTextArrayToMap converts text slice to pair map.
func responseTextArrayToMap(texts []*pb.Text) map[string]*models.Text {
	result := make(map[string]*models.Text)
	if len(texts) < 1 {
		return result
	}
	for _, v := range texts {
		result[v.Title] = ProtoToModelsText(v)
	}

	return result
}

// ProtoToModelsText converts proto Text data to local Text.
func ProtoToModelsText(t *pb.Text) *models.Text {
	return &models.Text{
		ID:        0,
		UserID:    0,
		Title:     t.GetTitle(),
		Body:      t.GetBody(),
		Comment:   t.GetComment(),
		Version:   t.GetVersion(),
		DeletedAt: sql.NullTime{},
	}
}

// responseBinArrayToMap converts bin slice to bin map.
func responseBinArrayToMap(bins []*pb.Bin) map[string]*models.Bin {
	result := make(map[string]*models.Bin)
	if len(bins) < 1 {
		return result
	}
	for _, v := range bins {
		result[v.Title] = ProtoToModelsBin(v)
	}

	return result
}

// ProtoToModelsBin converts proto Binary data to local Binary.
func ProtoToModelsBin(b *pb.Bin) *models.Bin {
	return &models.Bin{
		ID:        0,
		UserID:    0,
		Title:     b.GetTitle(),
		Body:      b.GetBody(),
		Comment:   b.GetComment(),
		Version:   b.GetVersion(),
		DeletedAt: sql.NullTime{},
	}
}

// responseCardArrayToMap converts card slice to card map.
func responseCardArrayToMap(cards []*pb.Card) map[string]*models.Card {
	result := make(map[string]*models.Card)
	if len(cards) < 1 {
		return result
	}
	for _, v := range cards {
		result[v.Title] = ProtoToModelsCard(v)
	}

	return result
}

// ProtoToModelsCard converts proto Card data to local Card.
func ProtoToModelsCard(c *pb.Card) *models.Card {
	return &models.Card{
		ID:             0,
		UserID:         0,
		Title:          c.GetTitle(),
		Number:         c.GetNumber(),
		ExpirationDate: c.GetExpdate(),
		Comment:        c.GetComment(),
		Version:        c.GetVersion(),
		DeletedAt:      sql.NullTime{},
	}
}
