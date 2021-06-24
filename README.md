# gosm
gosm is a golang library which implements writing [OSM pbf](https://wiki.openstreetmap.org/wiki/PBF_Format) files.

# How to use?
Please check example_test.go 

# Notes
1. When you use `AppendNodes`, `AppendWays` or `AppendRelations`, at most 8000 items are allowed in one append operation
2. Use `encoder.Flush(memberType MemberType)` when you finished writing of one member type.
