package ipdb

func DefaultDBList(pathPrefix string) List {
	filePath := ""
	if pathPrefix != "" {
		filePath = pathPrefix
	}
	return List{
		&DB{
			Name: "qqwry",
			NameAlias: []string{
				"chunzhen",
			},
			Format:    FormatQQWry,
			FilePath:  filePath,
			File:      "qqwry.dat",
			Languages: LanguagesZH,
			Types:     TypesIPv4,
		},
		&DB{
			Name: "geoip",
			NameAlias: []string{
				"geoip2",
				"geolite",
				"geolite2",
			},
			Format:    FormatMMDB,
			FilePath:  filePath,
			File:      "GeoLite2-City-1.mmdb",
			Languages: LanguagesAll,
			Types:     TypesIP,
		},
		&DB{
			Name: "ip2region_v1",
			NameAlias: []string{
				"i2r_v1",
			},
			Format:    FormatIP2RegionV1,
			FilePath:  filePath,
			File:      "ip2region.db",
			Languages: LanguagesZH,
			Types:     TypesIPv4,
		},
		&DB{
			Name: "ip2region_v2",
			NameAlias: []string{
				"i2r_v2",
			},
			Format:    FormatIP2RegionV2,
			FilePath:  filePath,
			File:      "ip2region.xdb",
			Languages: LanguagesZH,
			Types:     TypesIPv4,
		},
		&DB{
			Name:      "ip2location",
			Format:    FormatIP2Location,
			FilePath:  filePath,
			File:      "IP2LOCATION-LITE-DB11.IPV6.BIN",
			Languages: LanguagesEN,
			Types:     TypesIP,
		},
		//&DB{
		//	Name: "dbip",
		//	NameAlias: []string{
		//		"db-ip",
		//	},
		//	Format:    FormatMMDB,
		//	File:      "dbip.mmdb",
		//	Languages: LanguagesAll,
		//	Types:     TypesIP,
		//},
		//&DB{
		//	Name:      "ipip",
		//	Format:    FormatIPIP,
		//	File:      "ipipfree.ipdb",
		//	Languages: LanguagesZH,
		//	Types:     TypesIP,
		//},
	}
}
