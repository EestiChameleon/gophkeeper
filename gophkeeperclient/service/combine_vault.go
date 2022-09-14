package service

import (
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	"github.com/EestiChameleon/gophkeeper/models"
)

// CombineVault check for latest version from both passed storage. Returns vault with the latest versions from both.
func CombineVault(localVault, dbVault *models.VaultProto) *models.VaultProto {
	out := clstor.MakeVaultProto()
	//anti nil incoming data.
	if localVault == nil {
		localVault = clstor.MakeVaultProto()
	}
	if dbVault == nil {
		dbVault = clstor.MakeVaultProto()
	}

	for k, v := range localVault.Pair {
		if v.Version > dbVault.Pair[k].Version {
			out.Pair[k] = v
		} else {
			out.Pair[k] = dbVault.Pair[k]
		}
	}

	for k, v := range localVault.Text {
		if v.Version > dbVault.Text[k].Version {
			out.Text[k] = v
		} else {
			out.Text[k] = dbVault.Text[k]
		}
	}

	for k, v := range localVault.Bin {
		if v.Version > dbVault.Bin[k].Version {
			out.Bin[k] = v
		} else {
			out.Bin[k] = dbVault.Bin[k]
		}
	}

	for k, v := range localVault.Card {
		if v.Version > dbVault.Card[k].Version {
			out.Card[k] = v
		} else {
			out.Card[k] = dbVault.Card[k]
		}
	}

	return out
}
