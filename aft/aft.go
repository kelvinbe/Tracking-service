package aft

//Africa's Talking API logic

type AftClient struct {
	Username string
	ApiKey   string
}

func NewAftClient() (*AftClient, error) {
	// TODO: add implimentation
	return nil, nil
}

func (aft *AftClient) ActivateDevice(device_id string) error {
	// TODO: add implimentation
	return nil
}

func (aft *AftClient) DeactivateDevice(device_id string) error {
	// TODO: add implimentation
	return nil
}
