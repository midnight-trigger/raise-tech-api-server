package api

type void interface{}

type ImageRoutes struct {
	GetImages void `method:"GET" path:"api/v1/image"`
	PostImage void `method:"POST" path:"api/v1/image"`
}
