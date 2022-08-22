module killedthis/content-generator

go 1.19

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/ryanbradynd05/go-tmdb v0.0.0-20220721194547-2ab6191c6273
	killedthis/shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/kylelemons/go-gypsy v1.0.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace killedthis/shared => ../shared
