package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"pehks1980/be2_hw9/graphql1/graph/generated"
	"pehks1980/be2_hw9/graphql1/graph/model"
	"pehks1980/be2_hw9/graphql1/internal/pkg/repository"
	"strconv"
)

func (r *mutationResolver) CreatePriceListItem(ctx context.Context, input model.NewPriceListItem) (*model.PriceListItem, error) {
	var newPriceListEntity repository.Pricelist // элемент для репозитория
	newPriceListEntity.Id, _ = strconv.Atoi(input.PricelistID)
	newPriceListEntity.Good = input.Name
	newPriceListEntity.Price, _ = strconv.Atoi(input.Price)

	id, err := repository.Reposit.AddPriceListEntity(ctx, newPriceListEntity)
	if err != nil {
		log.Printf("add repo error: %v", err)
		return nil, err
	}

	var listItem model.PriceListItem
	listItem.ItemID = strconv.Itoa(id)
	listItem.PricelistID = input.PricelistID
	listItem.Name = input.Name
	listItem.Price = input.Price
	return &listItem, nil
}

func (r *mutationResolver) UpdatePriceListItem(ctx context.Context, input model.UpdPriceListItem) (*model.PriceListItem, error) {
	var newPriceListEntity repository.Pricelist
	id, _ := strconv.Atoi(input.PricelistID)
	newPriceListEntity.Id, _ = strconv.Atoi(input.ItemID)
	newPriceListEntity.Good = input.Name
	newPriceListEntity.Price, _ = strconv.Atoi(input.Price)
	err := repository.Reposit.UpdatePriceList(ctx, id, newPriceListEntity)
	if err != nil {
		log.Printf("update repo error: %v", err)
		return nil, err
	}
	//copy result as per input
	return &model.PriceListItem{
		PricelistID: input.PricelistID,
		ItemID:      input.ItemID,
		Name:        input.Name,
		Price:       input.Price,
	}, nil
}

func (r *mutationResolver) DeletePriceList(ctx context.Context, input model.ID) (*model.Result, error) {
	id, _ := strconv.Atoi(input.PricelistID)
	err := repository.Reposit.DelPriceList(ctx, id)
	if err != nil {
		log.Printf("delete repo error: %v", err)
		return &model.Result{Result: false}, err
	}
	return &model.Result{Result: true}, err
}

func (r *queryResolver) ListRows(ctx context.Context, pricelistID string) ([]*model.PriceListItem, error) {
	var res []*model.PriceListItem
	id, _ := strconv.Atoi(pricelistID)
	result, _ := repository.Reposit.GetPriceList(ctx, id)
	//repack to output if struct
	for _, dbrow := range result {
		price := strconv.Itoa(dbrow.Price)
		itemId := strconv.Itoa(dbrow.Id)
		row := model.PriceListItem{
			PricelistID: pricelistID,
			ItemID:      itemId,
			Name:        dbrow.Good,
			Price:       price,
		}
		res = append(res, &row)
	}
	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
