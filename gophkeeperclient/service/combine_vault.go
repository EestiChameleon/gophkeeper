package service

import (
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	"github.com/EestiChameleon/gophkeeper/models"
)

// CombineVault check for latest version from both passed storage. Returns vault with the latest versions from both.
func CombineVault(localVault, dbVault *models.Vault) *models.Vault {
	out := clstor.MakeVault()
	//anti nil incoming data.
	if localVault == nil {
		localVault = clstor.MakeVault()
	}
	if dbVault == nil {
		dbVault = clstor.MakeVault()
	}

	// save key-titles from both structures
	keysP := dataKeys("pair", localVault, dbVault)
	keysT := dataKeys("text", localVault, dbVault)
	keysB := dataKeys("bin", localVault, dbVault)
	keysC := dataKeys("card", localVault, dbVault)

	// sync Pairs
	for title := range keysP {
		out.Pair[title] = FindLatestPair(title, localVault.Pair, dbVault.Pair)
	}

	// sync Text
	for title := range keysT {
		out.Text[title] = FindLatestText(title, localVault.Text, dbVault.Text)
	}

	// sync Text
	for title := range keysB {
		out.Bin[title] = FindLatestBin(title, localVault.Bin, dbVault.Bin)
	}

	// sync Text
	for title := range keysC {
		out.Card[title] = FindLatestCard(title, localVault.Card, dbVault.Card)
	}

	return out
}

func dataKeys(dataType string, v1, v2 *models.Vault) map[string]struct{} {
	keys := map[string]struct{}{}

	switch dataType {
	case "pair":
		for k := range v1.Pair {
			keys[k] = struct{}{}
		}
		for k := range v2.Pair {
			keys[k] = struct{}{}
		}
	case "text":
		for k := range v1.Text {
			keys[k] = struct{}{}
		}
		for k := range v2.Text {
			keys[k] = struct{}{}
		}
	case "bin":
		for k := range v1.Bin {
			keys[k] = struct{}{}
		}
		for k := range v2.Bin {
			keys[k] = struct{}{}
		}
	case "card":
		for k := range v1.Card {
			keys[k] = struct{}{}
		}
		for k := range v2.Card {
			keys[k] = struct{}{}
		}
	}

	return keys
}

func FindLatestPair(title string, localMap, dbMap map[string]*models.Pair) *models.Pair {
	// search for the title in both storages
	loc, locOK := localMap[title]
	db, dbOK := dbMap[title]
	// if found in both - compare
	if locOK && dbOK {
		if loc.Version > db.Version {
			return loc
		} else {
			return db
		}
	}
	// found in local and not found in db => save local
	if locOK && !dbOK {
		return loc
	}
	if !locOK && dbOK {
		return db
	}

	return nil
}

func FindLatestText(title string, localMap, dbMap map[string]*models.Text) *models.Text {
	// search for the title in both storages
	loc, locOK := localMap[title]
	db, dbOK := dbMap[title]
	// if found in both - compare
	if locOK && dbOK {
		if loc.Version > db.Version {
			return loc
		} else {
			return db
		}
	}
	// found in local and not found in db => save local
	if locOK && !dbOK {
		return loc
	}
	if !locOK && dbOK {
		return db
	}

	return nil
}

func FindLatestBin(title string, localMap, dbMap map[string]*models.Bin) *models.Bin {
	// search for the title in both storages
	loc, locOK := localMap[title]
	db, dbOK := dbMap[title]
	// if found in both - compare
	if locOK && dbOK {
		if loc.Version > db.Version {
			return loc
		} else {
			return db
		}
	}
	// found in local and not found in db => save local
	if locOK && !dbOK {
		return loc
	}
	if !locOK && dbOK {
		return db
	}

	return nil
}

func FindLatestCard(title string, localMap, dbMap map[string]*models.Card) *models.Card {
	// search for the title in both storages
	loc, locOK := localMap[title]
	db, dbOK := dbMap[title]
	// if found in both - compare
	if locOK && dbOK {
		if loc.Version > db.Version {
			return loc
		} else {
			return db
		}
	}
	// found in local and not found in db => save local
	if locOK && !dbOK {
		return loc
	}
	if !locOK && dbOK {
		return db
	}

	return nil
}
