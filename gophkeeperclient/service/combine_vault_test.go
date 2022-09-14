package service

import (
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	halfEmptyDataOne = models.Vault{
		Pair: map[string]*models.Pair{
			"ptitle1": {
				Title:   "ptitle1",
				Login:   "log1",
				Pass:    "pass1",
				Comment: "comm1",
				Version: 1,
			},
		},
		Text: map[string]*models.Text{
			"ttitle1": {
				Title:   "ttitle1",
				Body:    "text1",
				Comment: "comm1",
				Version: 3,
			},
		},
		Bin:  nil,
		Card: nil,
	}

	halfEmptyDataTwo = models.Vault{
		Pair: map[string]*models.Pair{
			"ptitle1": {
				Title:   "ptitle1",
				Login:   "log2",
				Pass:    "pass2",
				Comment: "comm2",
				Version: 4,
			},
		},
		Text: map[string]*models.Text{
			"ttitle1": {
				Title:   "ttitle1",
				Body:    "text2",
				Comment: "comm2",
				Version: 2,
			},
		},
		Bin:  nil,
		Card: nil,
	}

	fullDataOne = models.Vault{
		Pair: map[string]*models.Pair{
			"ptitle1": {
				Title:   "ptitle1",
				Login:   "log1",
				Pass:    "pass1",
				Comment: "comm1",
				Version: 1,
			},
		},
		Text: map[string]*models.Text{
			"ttitle1": {
				Title:   "ttitle1",
				Body:    "text1",
				Comment: "comm1",
				Version: 3,
			},
		},
		Bin: map[string]*models.Bin{
			"btitle1": {
				Title:   "btitle1",
				Body:    []byte("binary data 1"),
				Comment: "comm1",
				Version: 4,
			},
		},
		Card: map[string]*models.Card{
			"ctitle1": {
				Title:          "ctitle1",
				Number:         "1111 1111 1111 1111",
				ExpirationDate: "1111/11",
				Comment:        "comm1",
				Version:        2,
			},
		},
	}

	fullDataTwo = models.Vault{
		Pair: map[string]*models.Pair{
			"ptitle1": {
				Title:   "ptitle1",
				Login:   "log2",
				Pass:    "pass2",
				Comment: "comm2",
				Version: 4,
			},
		},
		Text: map[string]*models.Text{
			"ttitle1": {
				Title:   "ttitle1",
				Body:    "text2",
				Comment: "comm2",
				Version: 2,
			},
		},
		Bin: map[string]*models.Bin{
			"btitle2": {
				Title:   "btitle2",
				Body:    []byte("binary data 2"),
				Comment: "comm2",
				Version: 3,
			},
		},
		Card: map[string]*models.Card{
			"ctitle2": {
				Title:          "ctitle2",
				Number:         "2222 2222 2222 2222",
				ExpirationDate: "2222/22",
				Comment:        "comm2",
				Version:        1,
			},
		},
	}
)

func TestCombineVault(t *testing.T) {
	type want struct {
		dataFinal *models.Vault
	}

	tests := []struct {
		name    string
		dataOne *models.Vault
		dataTwo *models.Vault
		want    want
	}{
		{
			name:    "Test #1: empty incoming data",
			dataOne: nil,
			dataTwo: nil,
			want: want{
				dataFinal: clstor.MakeVault(), //just empty struct
			},
		},
		{
			name:    "Test #2: one data empty, second data half-empty",
			dataOne: nil,
			dataTwo: &halfEmptyDataTwo,
			want: want{
				dataFinal: &halfEmptyDataTwo,
			},
		},
		{
			name:    "Test #3: one data empty, second data full",
			dataOne: nil,
			dataTwo: &fullDataTwo,
			want: want{
				dataFinal: &fullDataTwo,
			},
		},
		{
			name:    "Test #4: one data half-empty, second data half-empty",
			dataOne: &halfEmptyDataOne,
			dataTwo: &halfEmptyDataTwo,
			want: want{
				dataFinal: &models.Vault{
					Pair: map[string]*models.Pair{
						"ptitle1": {
							Title:   "ptitle1",
							Login:   "log2",
							Pass:    "pass2",
							Comment: "comm2",
							Version: 4,
						},
					},
					Text: map[string]*models.Text{
						"ttitle1": {
							Title:   "ttitle1",
							Body:    "text1",
							Comment: "comm1",
							Version: 3,
						},
					},
					Bin:  nil,
					Card: nil,
				},
			},
		},
		{
			name:    "Test #5: one data half-empty, second data full",
			dataOne: &halfEmptyDataOne,
			dataTwo: &fullDataTwo,
			want: want{
				dataFinal: &models.Vault{
					Pair: map[string]*models.Pair{
						"ptitle1": {
							Title:   "ptitle1",
							Login:   "log2",
							Pass:    "pass2",
							Comment: "comm2",
							Version: 4,
						},
					},
					Text: map[string]*models.Text{
						"ttitle1": {
							Title:   "ttitle1",
							Body:    "text1",
							Comment: "comm1",
							Version: 3,
						},
					},
					Bin: map[string]*models.Bin{
						"btitle2": {
							Title:   "btitle2",
							Body:    []byte("binary data 2"),
							Comment: "comm2",
							Version: 3,
						},
					},
					Card: map[string]*models.Card{
						"ctitle2": {
							Title:          "ctitle2",
							Number:         "2222 2222 2222 2222",
							ExpirationDate: "2222/22",
							Comment:        "comm2",
							Version:        1,
						},
					},
				},
			},
		},
		{
			name:    "Test #6: one data full, second data full",
			dataOne: &fullDataOne,
			dataTwo: &fullDataTwo,
			want: want{
				dataFinal: &models.Vault{
					Pair: map[string]*models.Pair{
						"ptitle1": {
							Title:   "ptitle1",
							Login:   "log2",
							Pass:    "pass2",
							Comment: "comm2",
							Version: 4,
						},
					},
					Text: map[string]*models.Text{
						"ttitle1": {
							Title:   "ttitle1",
							Body:    "text1",
							Comment: "comm1",
							Version: 3,
						},
					},
					Bin: map[string]*models.Bin{
						"btitle1": {
							Title:   "btitle1",
							Body:    []byte("binary data 1"),
							Comment: "comm1",
							Version: 4,
						},
						"btitle2": {
							Title:   "btitle2",
							Body:    []byte("binary data 2"),
							Comment: "comm2",
							Version: 3,
						},
					},
					Card: map[string]*models.Card{
						"ctitle1": {
							Title:          "ctitle1",
							Number:         "1111 1111 1111 1111",
							ExpirationDate: "1111/11",
							Comment:        "comm1",
							Version:        2,
						},
						"ctitle2": {
							Title:          "ctitle2",
							Number:         "2222 2222 2222 2222",
							ExpirationDate: "2222/22",
							Comment:        "comm2",
							Version:        1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CombineVault(tt.dataOne, tt.dataTwo)
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
