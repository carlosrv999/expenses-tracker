package tag

type createRequest struct {
	TagName string  `json:"tag_name"`
	Color   *string `json:"color"`
	Icon    *string `json:"icon"`
}

type updateRequest struct {
	TagName string  `json:"tag_name"`
	Color   *string `json:"color"`
	Icon    *string `json:"icon"`
}
