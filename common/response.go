package common

type Response[T any] struct {
	Success    bool           `json:"success"`
	Messages   []any          `json:"messages"` // Use `string` instead of `any` if they're always text
	Result     []T              `json:"result"`
	Pagination PaginationInfo `json:"pagination"`
}

type OneResponse[T any] struct {
	Success    bool           `json:"success"`
	Messages   []any          `json:"messages"` // Use `string` instead of `any` if they're always text
	Result     T              `json:"result"`
	Pagination PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	Count int    `json:"count"`
	Next  string `json:"next"`
}
