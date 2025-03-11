module github.com/katheland/pokedexcli

go 1.24.1

require internal/pokeapi v0.0.0
replace internal/pokeapi => ./internal/pokeapi

require internal/pokecache v0.0.0 // indirect
replace internal/pokecache => ./internal/pokecache