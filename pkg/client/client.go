package client

type Client struct {
	Adapters AdapterService
	/*
		TODO
		Tags TagService
		Jobs JobService
		Runs RunService
		Events EventService
		About AboutService
		Snippets SnippetService
		Ws WsService
	*/
}

func New(url string) *Client {
	return &Client{
		Adapters: newAdapterService(url),
	}
}
