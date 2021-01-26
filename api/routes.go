package api

type void interface{}

type ImageRoutes struct {
	HealthCheck void `method:"GET" path:"health"`
	PostImage   void `method:"POST" path:"api/v1/image"`
}
