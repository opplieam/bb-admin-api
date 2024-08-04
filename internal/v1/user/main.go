package user

type Storer interface {
	AuthedI
	ManageI
}

type Handler struct {
	Store Storer
}

func NewHandler(store Storer) *Handler {
	return &Handler{
		Store: store,
	}
}
