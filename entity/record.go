package entity

type Record struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`
}

type RecordRow struct {
    ID      int            `json:"id"`
    Key     string         `json:"key"`
    Value   string         `json:"value"`
}

func (d *Record) Copy() Record {
	values := d.Data

	newMap := map[string]string{}
	for key, value := range values {
		newMap[key] = value
	}

	return Record{
		ID:   d.ID,
		Data: newMap,
	}
}
