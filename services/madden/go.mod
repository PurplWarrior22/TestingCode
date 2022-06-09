module ../services/madden

go 1.17

require gorm.io/gorm v1.23.2

require (
	../services/echopprof v0.0.0-20220304154656-126cd83bba02
	../services/maddendb v0.0.0-20220406174729-01aa20dd59e2
	../services/models v0.0.0-20220529142624-60b3349973d5
	../services/utilities v0.0.0-20220529142624-60b3349973d5
	github.com/deepmap/oapi-codegen v1.9.1
	github.com/getkin/kin-openapi v0.90.0
	github.com/labstack/echo/v4 v4.6.3
)
