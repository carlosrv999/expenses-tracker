package category

type createRequest struct {
	ParentCategoryID *int64  `json:"parent_category_id"`
	CategoryName     string  `json:"category_name"`
	Icon             *string `json:"icon"`
	Color            *string `json:"color"`
}

type updateRequest struct {
	ParentCategoryID *int64  `json:"parent_category_id"`
	CategoryName     string  `json:"category_name"`
	Icon             *string `json:"icon"`
	Color            *string `json:"color"`
}
