package driver

import "encoding/json"

var _ Driver = (*TileDriver)(nil)

func NewTileDriver() *TileDriver {
	return &TileDriver{
		PathParser: SlashPathParser,
		Realizer:   new(StdRealizer),
		Modem: &GeneralModem[*RawProcessor]{
			Marshaler:   json.Marshal,
			Unmarshaler: json.Unmarshal,
		},
	}
}

// TileDriver is a driver for raw type rule tree
type TileDriver struct {
	PathParser
	Realizer
	Modem
}

func (TileDriver) Name() string { return "tile" }
