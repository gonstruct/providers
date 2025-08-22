package mailables

import "net/mail"

type address = mail.Address
type addressSlice []*address

func (addresses addressSlice) String() []string {
	stringAddresses := make([]string, len(addresses))
	for i, address := range addresses {
		stringAddresses[i] = address.String()
	}

	return stringAddresses
}

func Address(email, name string) *address {
	return &mail.Address{
		Address: email,
		Name:    name,
	}
}

func Addresses[S Sender](addresses ...S) addressSlice {
	if len(addresses) == 0 {
		return nil
	}

	slice := make(addressSlice, len(addresses))
	for i, addr := range addresses {
		switch email := any(addr).(type) {
		case string:
			slice[i] = Address("", email)
		case *address:
			slice[i] = email
		}
	}

	return slice
}
