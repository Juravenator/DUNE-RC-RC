package internal

type ServiceAddress struct {
	Address string
	Port int
}

func GetFirstServiceAddress(serviceName string) ServiceAddress {
	all := GetServiceAddresses(serviceName)
	if len(all) == 0 {
		panic("no such service")
	}
	return all[0]
}