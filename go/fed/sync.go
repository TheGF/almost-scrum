package fed

//type Status struct {
//	Exchanges map[string]bool      `json:"exchanges"`
//	Locs      *sync.Map            `json:"locs"`
//	Throughput      map[string]Throughput      `json:"throughput"`
//	Exports   map[string]time.ModTime `json:"exports"`
//}
//
//func GetStatus(project *core.Project) Status {
//	connection, err := Connect(project)
//	if err != nil {
//		return Status{}
//	}
//
//	s := Status{
//		Exchanges: map[string]bool{},
//		Locs:      &connection.locs,
//		Throughput:      map[string]Throughput{},
//		Exports:   connection.exports,
//	}
//	for exchange, connected := range connection.exchanges {
//		s.Exchanges[exchange.Name()] = connected
//	}
//	for exchange, throughput := range connection.throughput {
//		s.Throughput[exchange.Name()] = throughput
//	}
//
//	return s
//}
