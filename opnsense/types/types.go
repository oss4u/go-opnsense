package types

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == `"1"` {
		*b = Bool(true)
	} else {
		*b = false
	}
	return nil
}

func (h Bool) MarshalJSON() ([]byte, error) {
	if h {
		return []byte(`"1"`), nil
	} else {
		return []byte(`"0"`), nil
	}
}
