package lib

type _region uint8

type _options struct {
	China         _region
	International _region
}

const (
	international _region = iota
	china
)

var Options _options = _options{
	China:         china,
	International: international,
}

type wlocArgs struct {
	region _region
}

func newWlocArgs() wlocArgs {
	return wlocArgs{
		region: international,
	}
}

type Modifier func(*wlocArgs)

func (o _options) WithRegion(region _region) Modifier {
	return func(wa *wlocArgs) {
		wa.region = region
	}
}
