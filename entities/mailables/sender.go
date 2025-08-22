package mailables

type Sender interface {
	~string | ~*address
}
