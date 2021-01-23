module github.com/paketo-buildpacks/microsoft-azure

go 1.15

require (
	github.com/buildpacks/libcnb v1.18.1
	github.com/onsi/gomega v1.10.4
	github.com/paketo-buildpacks/libpak v1.50.1
	github.com/sclevine/spec v1.4.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/paketo-buildpacks/libpak => ../../paketo-buildpacks/libpak
