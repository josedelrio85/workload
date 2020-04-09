package workload

import "net/http"

// Modelable is blablabla
type Modelable interface {
	// TableName() string
	GetElements(Input, Storer) error
	CreateElement(Input, Storer) error
	UpdateElement(Input, Storer) error
	DeleteElement(Input, Storer) error
	GetType() int64
	Response(http.ResponseWriter)
}
