package main

import . "github.com/enetx/g"

func main() {
	NewDir(".").
		Walk().
		Exclude(func(f *File) bool { return f.IsDir() && f.Dir().Ok().IsLink() }).
		Exclude((*File).IsLink).
		ForEach(func(v Result[*File]) {
			if v.IsOk() {
				v.Ok().Path().Ok().Println()
			}
		})
}
