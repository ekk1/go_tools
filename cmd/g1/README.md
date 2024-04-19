# G1

## Basic Architecture

* Players: a player can have multiple cities

* City: a city consists of:
    * buildings:        grow food ingredients, factories, schools for research, producing units
    * working units:    normal workers, special workers, working vehicles, like cars, planes etc (or should planes be managed by building section?), can be assigned to buildings or working teams
    * troops:           army units, can be assigned to marching teams
    * working teams:    a group of units, can be assigned to resources nearby, or building an outpost or new city
    * marching teams:   a group of armys, can be sent to attack other cities and outposts

* Units: all units are currently only managed by City and outposts, player can not control units directly

* Loop: Server's main loop will be executed every 10s (pending review), and user inputs will be freezed during
    * a loop execution loops the city list, and loops all buildings, teams, special event etc.
    * player can assign units, build buildings, change formulas at anytime, but the timeline won't be moving

* Battle: Currently, a battle is sending a marching team to a location, and cities have a observing range, if army is within range, player can see a warning, indicating the speed and army details, defence towers and remaining troops will be able to attack the inbound army within range
    * All cities and outposts can attack teams nearby, this is achived in the ATTACKING city's main loop. it will calculate the army location, and loop all cities and outputs (maybe nearby?), and if within that city's range, it will be attacked (and will attack back if it's also within the troops range, but the team never stops until reaching the destination)

* How to use multi thread ?
