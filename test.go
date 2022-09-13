package main

import (
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m1 := storage.MakeVaultProto()
	m2 := storage.MakeVaultProto()
	m3 := storage.MakeVaultProto()

	m1.Pair["ptitle1"] = &pb.Pair{Title: "ptitle1", Version: 2} // this
	m2.Pair["ptitle1"] = &pb.Pair{Title: "ptitle1", Version: 1}

	m1.Pair["ptitle2"] = &pb.Pair{Title: "ptitle2", Version: 3}
	m2.Pair["ptitle2"] = &pb.Pair{Title: "ptitle2", Version: 4} // this

	m1.Card["ctitle1"] = &pb.Card{Title: "ctitle1", Version: 1}
	m2.Card["ctitle1"] = &pb.Card{Title: "ctitle1", Version: 3} // this

	m1.Text["ttitle1"] = &pb.Text{Title: "ttitle1", Version: 5} // this
	m2.Text["ttitle1"] = &pb.Text{Title: "ttitle1", Version: 3}

	for k, v := range m1.Pair {
		if v.Version > m2.Pair[k].Version {
			m3.Pair[k] = v
		} else {
			m3.Pair[k] = m2.Pair[k]
		}
	}

	for k, v := range m1.Text {
		if v.Version > m2.Text[k].Version {
			m3.Text[k] = v
		} else {
			m3.Text[k] = m2.Text[k]
		}
	}

	for k, v := range m1.Card {
		if v.Version > m2.Card[k].Version {
			m3.Card[k] = v
		} else {
			m3.Card[k] = m2.Card[k]
		}
	}

	fmt.Println(m3)

}
