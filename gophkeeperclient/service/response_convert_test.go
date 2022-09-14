package service

import (
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVaultSyncConvert(t *testing.T) {
	type want struct {
		dataFinal *models.Vault
	}

	tests := []struct {
		name         string
		incomingData *pb.SyncVaultResponse
		want         want
	}{
		{
			name: "Test #1: empty incoming data",
			incomingData: &pb.SyncVaultResponse{
				Pairs:   nil,
				Texts:   nil,
				BinData: nil,
				Cards:   nil,
				Status:  "not found",
			},
			want: want{
				dataFinal: clstor.MakeVault(),
			},
		},
		{
			name: "Test #2: full incoming data",
			incomingData: &pb.SyncVaultResponse{
				Pairs: []*pb.Pair{
					{
						Title:   "p1",
						Login:   "l1",
						Pass:    "p1",
						Comment: "c1",
						Version: 1,
					},
				},
				Texts: []*pb.Text{
					{
						Title:   "t1",
						Body:    "b1",
						Comment: "c1",
						Version: 2,
					},
				},
				BinData: []*pb.Bin{
					{
						Title:   "b1",
						Body:    []byte(`byte1`),
						Comment: "c1",
						Version: 3,
					},
				},
				Cards: []*pb.Card{
					{
						Title:   "c1",
						Number:  "1111 1111 1111 1111",
						Expdate: "2022/22",
						Comment: "c1",
						Version: 4,
					},
				},
				Status: "",
			},
			want: want{
				dataFinal: &models.Vault{
					Pair: map[string]*models.Pair{
						"p1": {
							Title:   "p1",
							Login:   "l1",
							Pass:    "p1",
							Comment: "c1",
							Version: 1,
						},
					},
					Text: map[string]*models.Text{
						"t1": {
							Title:   "t1",
							Body:    "b1",
							Comment: "c1",
							Version: 2,
						},
					},
					Bin: map[string]*models.Bin{
						"b1": {
							Title:   "b1",
							Body:    []byte(`byte1`),
							Comment: "c1",
							Version: 3,
						},
					},
					Card: map[string]*models.Card{
						"c1": {
							Title:          "c1",
							Number:         "1111 1111 1111 1111",
							ExpirationDate: "2022/22",
							Comment:        "c1",
							Version:        4,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VaultSyncConvert(tt.incomingData)
			for k, v := range result.Pair {
				assert.Equal(t, tt.want.dataFinal.Pair[k], v)
			}
			for k, v := range result.Text {
				assert.Equal(t, tt.want.dataFinal.Text[k], v)
			}
			for k, v := range result.Bin {
				assert.Equal(t, tt.want.dataFinal.Bin[k], v)
			}
			for k, v := range result.Card {
				assert.Equal(t, tt.want.dataFinal.Card[k], v)
			}
		})
	}
}
