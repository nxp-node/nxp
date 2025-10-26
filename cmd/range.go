package cmd

type Range struct {
	Minimum uint16
	Maximum *uint16
}

func MakeRange(Minimum uint16, Maximum int16) Range {
	var max *uint16

	if Maximum == -1 {
		max = nil
	} else if Maximum < 0 {
		panic("the 'Maximum' parameter cannot be smaller than -1")
	} else {
		umax := uint16(Maximum)
		max = &umax
	}

	return Range{
		Minimum: Minimum,
		Maximum: max,
	}
}
