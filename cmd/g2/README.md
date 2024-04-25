## Config

### resource
* size: storage size
* value: sell price, and required work to grow or make
* speed: speed factor
* output: output num for every `value` work made

### unit
* pop: population takes
* consume: food consume rate
* workspeed: speed factor
* movespeed: speed factor
* load: carry ability
    * these two are mainly used in teams

### building
* resouces: resources required to build, this should be satifisfied befrre adding building
* work: required work to finish build
* maxunits: max working units
* size: building size

### city
* resouces: resources required to build, this should be satifisfied befrre adding building
* work: required work to finish build
* population_cap: default max population, can be increased by building
* storage_cap: default max storage, can be increased by building
* building_cap: default building size, can be increased only by upgrading city
* scan_range: ability to scan for outer resources and inbound enemies
