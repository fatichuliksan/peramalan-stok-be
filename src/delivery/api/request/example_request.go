package request

// PostExample ...
type PostExample struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}
