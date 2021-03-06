package api

import (
	"context"
	_ "fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	_ "io"
	"log"
	pricelists "pehks1980/be2_hw9/grpc"
	"pehks1980/be2_hw9/internal/pkg/repository"
	"strconv"
	_ "time"
)

var _ pricelists.PricelistsServer = &Server{}

//Server - app structs
type Server struct {
	Ctx  context.Context
	Repo *repository.Repo
	pricelists.UnimplementedPricelistsServer
}

//CreatePriceListItem - create pricelist entity
func (srv *Server) CreatePriceListItem(ctx context.Context, req *pricelists.PriceListItem) (*pricelists.PriceListItem, error) {

	var newPriceListEntity repository.Pricelist // элемент для репозитория
	newPriceListEntity.Id, _ = strconv.Atoi(req.PricelistId)
	newPriceListEntity.Good = req.Name
	newPriceListEntity.Price, _ = strconv.Atoi(req.Price)

	id, err := srv.Repo.AddPriceListEntity(ctx, newPriceListEntity)
	if err != nil {
		log.Printf("add repo error: %v", err)
		return nil, err
	}
	req.ItemId = strconv.Itoa(id)

	return req, nil
}

//DeletePriceList -delete pricelist
func (srv *Server) DeletePriceList(ctx context.Context, req *pricelists.ID) (*emptypb.Empty, error) {
	id, _ := strconv.Atoi(req.Id)
	err := srv.Repo.DelPriceList(ctx, id)
	if err != nil {
		log.Printf("delete repo error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

//UpdatePriceListItem - update pricelist entity(row)
func (srv *Server) UpdatePriceListItem(ctx context.Context, req *pricelists.PriceListItem) (*emptypb.Empty, error) {
	var newPriceListEntity repository.Pricelist
	id, _ := strconv.Atoi(req.PricelistId)
	newPriceListEntity.Id, _ = strconv.Atoi(req.ItemId)
	newPriceListEntity.Good = req.Name
	newPriceListEntity.Price, _ = strconv.Atoi(req.Price)
	err := srv.Repo.UpdatePriceList(ctx, id, newPriceListEntity)
	if err != nil {
		log.Printf("Update repo error: %v", err)
		return nil, err
	}
	//no data updated here -)
	return &empty.Empty{}, nil
}

//GetPriceList - get pricelist as array of rows
func (srv *Server) GetPriceList(ctx context.Context, req *pricelists.ID) (*pricelists.PriceListRowsResponse, error) {
	res := []*pricelists.PriceListItem{}
	id, _ := strconv.Atoi(req.Id)
	result, _ := srv.Repo.GetPriceList(ctx, id)
	//repack to output if struct
	for _, dbrow := range result {
		price := strconv.Itoa(dbrow.Price)
		itemId := strconv.Itoa(dbrow.Id)
		row := pricelists.PriceListItem{
			PricelistId: req.Id,
			ItemId:      itemId,
			Name:        dbrow.Good,
			Price:       price,
		}
		res = append(res, &row)
	}
	//finish - insert to array struct
	ans := pricelists.PriceListRowsResponse{ListRows: res}
	return &ans, nil
}
