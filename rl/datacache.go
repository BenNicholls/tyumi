package rl

import "github.com/bennicholls/tyumi/log"

type dataCache[D any, DT ~uint32] struct {
	cache []D
}

func (dc dataCache[D, DT]) validType(data_type DT) bool {
	return int(data_type) < len(dc.cache)
}

func (dc dataCache[D, DT]) GetData(data_type DT) (data D) {
	if !dc.validType(data_type) {
		log.Error("DataType not registered.")
		return
	}

	return dc.cache[data_type]
}

func (dc *dataCache[D, DT]) RegisterDataType(data D) DT {
	if dc.cache == nil {
		dc.cache = make([]D, 0)
	}

	dc.cache = append(dc.cache, data)
	return DT(len(dc.cache) - 1)
}
