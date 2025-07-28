package util

import "github.com/bennicholls/tyumi/log"

// DataCache is a generic cache for flyweight-style data. It takes 2 type arguments: a type describing the data you want
// to cache (D), and a numeric type for indexing the data (DT). Add data to the cache with RegisterDataType; it will
// return a DT, which you can then use to retrieve data later with GetData().
//
// Example: for storing tile data! Create a struct for holding data common to types of tiles, and a TileType type for
//
//	       the indexes.
//
//	type TileType uint32
//
//	type TileData struct {
//	   Colour
//	}
//
// Then we make a DataCache, and register some types of tiles:
//
//	var tileCache = DataCache[TileData, TileType]
//
//	var GrassTile = tileCache.RegisterDataType( TileData{ col.GREEN } )
//	var DirtTile = tileCache.RegisterDataType( TileData{ col.BROWN } )
//	var WallTile = tileCache.RegisterDataType( TileData{ col.GREY } )
//
// Now, when making your tilemap, you can just store the DT values you got back from the datacache instead of storing
// the entire TileData struct. Big memory savings! Then later when you need to get the data:
//
//	tile := tilemap[0] // get the type of the first tile in the tilemap
//	colourToDraw := tileCache.GetData(tile).Colour // get its colour from the cache!
//
// Easy peasy! Of course this example is kind of silly, but you can imagine the memory savings possible if the tiledata
// struct had dozens of fields in it.
type DataCache[D any, DT ~uint32] struct {
	cache []D
}

func (dc DataCache[D, DT]) validType(data_type DT) bool {
	return int(data_type) < len(dc.cache)
}

func (dc DataCache[D, DT]) GetData(data_type DT) (data D) {
	if !dc.validType(data_type) {
		log.Error("DataType not registered.")
		return
	}

	return dc.cache[data_type]
}

func (dc *DataCache[D, DT]) ReplaceData(data_type DT, replacement D) {
	if !dc.validType(data_type) {
		log.Error("DataType not registered.")
		return
	}

	dc.cache[data_type] = replacement
}

func (dc *DataCache[D, DT]) RegisterDataType(data D) DT {
	if dc.cache == nil {
		dc.cache = make([]D, 0)
	}

	dc.cache = append(dc.cache, data)
	return DT(len(dc.cache) - 1)
}
