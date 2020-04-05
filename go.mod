module github.com/nxzz/HEM-GW16A-exporter

go 1.14

require (
	github.com/buger/jsonparser v0.0.0-20200322175846-f7e751efca13
	github.com/prometheus/client_golang v1.5.1
	local.packages/HEMGW16A v0.0.0-00010101000000-000000000000
)

replace local.packages/HEMGW16A => ./HEMGW16A
