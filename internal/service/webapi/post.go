package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"imman/parser_service/internal/entity"
)

type PostWebAPI struct {
	client *http.Client
}

func NewPostWebAPI(client *http.Client) *PostWebAPI {
	return &PostWebAPI{client: client}
}

func (p *PostWebAPI) GetData(ctx context.Context, page int) ([]entity.Post, error) {
	url := fmt.Sprintf("https://gorest.co.in/public/v1/posts?page=%d", page)
	response, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	type responseData struct {
		Data []entity.Post `json:"data"`
	}

	var postData responseData
	err = json.NewDecoder(response.Body).Decode(&postData)
	if err != nil {
		return nil, err
	}

	return postData.Data, nil
}
