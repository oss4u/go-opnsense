package types

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == `"1"` {
		*b = Bool(true)
		return nil
	}

	*b = false
	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`"1"`), nil
	}

	return []byte(`"0"`), nil
}

func (b *Bool) Bool() bool {
	return *(*bool)(b)
}
