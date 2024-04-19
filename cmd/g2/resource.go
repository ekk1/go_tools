package main

type Resource string

const (
	ResourceIron   Resource = "iron"
	ResourceCopper Resource = "copper"
	ResourceSilver Resource = "silver"
	ResourceGold   Resource = "gold"
	ResourceWood   Resource = "wood"
	ResourceCorn   Resource = "corn"
)

type ResourceConfig struct {
	ResourceSize  map[Resource]int64
	ResourceValue map[Resource]int64
}

var GlobalResourceConfig = &ResourceConfig{
	ResourceSize: map[Resource]int64{
		ResourceCorn: 1,
	},
}
