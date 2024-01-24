package service

import (
	"context"
	"io"

	"github.com/raulsilva-tech/gRPC/internal/database"
	"github.com/raulsilva-tech/gRPC/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {

	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	pbCategory := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return pbCategory, nil

}

func (c *CategoryService) ListCategories(ctx context.Context, blank *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	pbCategories := make([]*pb.Category, len(categories))
	for i, category := range categories {
		pbCategories[i] = &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
	}

	return &pb.CategoryList{
		Categories: pbCategories,
	}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCategoryRequest) (*pb.Category, error) {

	category, err := c.CategoryDB.FindById(in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil

}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {

	categoryList := &pb.CategoryList{}

	for {

		categoryReceived, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categoryList)
		}
		if err != nil {
			return err
		}

		category, err := c.CategoryDB.Create(categoryReceived.Name, categoryReceived.Description)
		if err != nil {
			return err
		}

		categoryList.Categories = append(categoryList.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})

	}

}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {

	for {

		categoryReceived, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		category, err := c.CategoryDB.Create(categoryReceived.Name, categoryReceived.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
		if err != nil {
			return err
		}

	}

}
