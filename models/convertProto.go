package models

import (
	"database/sql"
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

// ProtoToModelsPair converts proto Pair data to local Pair.
func ProtoToModelsPair(p *pb.Pair) *Pair {
	return &Pair{
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

// ProtoToModelsText converts proto Text data to local Text.
func ProtoToModelsText(t *pb.Text) *Text {
	return &Text{
		ID:        0,
		UserID:    0,
		Title:     t.GetTitle(),
		Body:      t.GetBody(),
		Comment:   t.GetComment(),
		Version:   t.GetVersion(),
		DeletedAt: sql.NullTime{},
	}
}

// ProtoToModelsBin converts proto Binary data to local Binary.
func ProtoToModelsBin(b *pb.Bin) *Bin {
	return &Bin{
		ID:        0,
		UserID:    0,
		Title:     b.GetTitle(),
		Body:      b.GetBody(),
		Comment:   b.GetComment(),
		Version:   b.GetVersion(),
		DeletedAt: sql.NullTime{},
	}
}

// ProtoToModelsCard converts proto Card data to local Card.
func ProtoToModelsCard(c *pb.Card) *Card {
	return &Card{
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

// ModelsToProtoPair converts local Pair structure to proto Pair structure.
func ModelsToProtoPair(in *Pair) *pb.Pair {
	return &pb.Pair{
		Title:   in.Title,
		Login:   in.Login,
		Pass:    in.Pass,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// ModelsToProtoText converts local Text structure to proto Text structure.
func ModelsToProtoText(in *Text) *pb.Text {
	return &pb.Text{
		Title:   in.Title,
		Body:    in.Body,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// ModelsToProtoBin converts local Bin structure to proto Bin structure.
func ModelsToProtoBin(in *Bin) *pb.Bin {
	return &pb.Bin{
		Title:   in.Title,
		Body:    in.Body,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// ModelsToProtoCard converts local Card structure to proto Card structure.
func ModelsToProtoCard(in *Card) *pb.Card {
	return &pb.Card{
		Title:   in.Title,
		Number:  in.Number,
		Expdate: in.ExpirationDate,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// convertActualDataToProto converts local data structures, used for DB interactions, to gRPC proto structures.
func ActualDataToProto(in *ActualData) *ActualProtoData {
	out := new(ActualProtoData)
	for _, v := range in.Pairs {
		out.Pairs = append(out.Pairs, ModelsToProtoPair(v))
	}

	for _, v := range in.Texts {
		out.Texts = append(out.Texts, ModelsToProtoText(v))
	}

	for _, v := range in.Bins {
		out.Bins = append(out.Bins, ModelsToProtoBin(v))
	}

	for _, v := range in.Cards {
		out.Cards = append(out.Cards, ModelsToProtoCard(v))
	}

	return out
}
