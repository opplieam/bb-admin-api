package user

type Storer interface {
	AuthedI
}

type Handler struct {
	Store Storer
}

func NewHandler(store Storer) *Handler {
	return &Handler{
		Store: store,
	}
}
