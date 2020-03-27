package main

import "github.com/altid/libs/fs"

var Commands = []*fs.Command{
	{
		Name: "block",
		Args: []string{"<phone number>"},
		Heading: fs.DefaultGroup,
		Description: "Block a phone number",
	},
}