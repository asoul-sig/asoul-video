// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

type MemberSecUID string

const (
	MemberSecUIDAva    = MemberSecUID("MS4wLjABAAAAxOXMMwlShWjp4DONMwfEEfloRYiC1rXwQ64eydoZ0ORPFVGysZEd4zMt8AjsTbyt")
	MemberSecUIDBella  = MemberSecUID("MS4wLjABAAAAlpnJ0bXVDV6BNgbHUYVWnnIagRqeeZyNyXB84JXTqAS5tgGjAtw0ZZkv0KSHYyhP")
	MemberSecUIDCarol  = MemberSecUID("MS4wLjABAAAAuZHC7vwqRhPzdeTb24HS7So91u9ucl9c8JjpOS2CPK-9Kg2D32Sj7-mZYvUCJCya")
	MemberSecUIDDiana  = MemberSecUID("MS4wLjABAAAA5ZrIrbgva_HMeHuNn64goOD2XYnk4ItSypgRHlbSh1c")
	MemberSecUIDEileen = MemberSecUID("MS4wLjABAAAAxCiIYlaaKaMz_J1QaIAmHGgc3bTerIpgTzZjm0na8w5t2KTPrCz4bm_5M5EMPy92")
)

type ReportType string

const (
	ReportTypeUpdateMember ReportType = "update_member"
	ReportTypeCreateVideo  ReportType = "create_video"
)
