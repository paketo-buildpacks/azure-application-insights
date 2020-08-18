module github.com/paketo-buildpacks/azure-application-insights

go 1.15

require (
	github.com/buildpacks/libcnb v1.16.0
	github.com/onsi/gomega v1.10.1
	github.com/paketo-buildpacks/libpak v1.40.0
	github.com/sclevine/spec v1.4.0
	github.com/stretchr/testify v1.6.1
)

replace github.com/paketo-buildpacks/libpak => ../libpak
