# gosm
gosm is a golang library which implements writing [OSM pbf](https://wiki.openstreetmap.org/wiki/PBF_Format) files. The initial idea is that when we use osmosis tool to dump data from database to pbf file, it causes high load of the database during that time. We need to maintain a large size of database just for the dump. For osmosis, we do not have control on how to query the database, so we implemented this library in golang. With this library, we can do the database query by ourselves during the data dump phase and we managed to reduce half size of the original database to finish the data dump.

So if you are looking for a full control during the data dump in your golang code, gosm is your choice. 
This is running in our production usage and under a stable version right now.

# Usage
## Quick start
```
	// initialize an encoder
	wc := &myWriter{}
	encoder := NewEncoder(&NewEncoderRequiredInput{
		RequiredFeatures: []string{"OsmSchema-V0.6"},
		Writer:           wc,
	},
		WithWritingProgram("example"),
		WithZlipEnabled(false),
	)
	errChan, err := encoder.Start()

	// write some nodes
	nodes := []*Node{
		{
			ID: 7278995748,
			Latitude:  -7.2380901,
			Longitude: 112.6773289,
			Tags: map[string]string{
				"node": "node1",
				"erp":  "no",
			},
		},
		{
			ID: 6978510772,
			Latitude:  -7.2381273,
			Longitude: 112.6775354,
		},
	}
	encoder.AppendNodes(nodes)
	encoder.Close()
```



## End to end demo
This is not a fully runnable code, would like to demo how to use the library.

You can replace oriNodes, oriWays, oriRelations with your own data.
```
	encoder := gosm.NewEncoder(&gosm.NewEncoderRequiredInput{
		RequiredFeatures: []string{"OsmSchema-V0.6", "DenseNodes"},
		Writer:           f,
	},
		gosm.WithWritingProgram("wp1"),
		gosm.WithZlipEnabled(true),
	)
	defer func() {
		_ = encoder.Close()
	}()

	errChan, err := encoder.Start()
	if err != nil {
		return err
	}

	var errs []error
	go func() {
		for e := range errChan {
			errs = append(errs, e)
		}
	}()

    sorts := make([]*Node, len(oriNodes))
	i := 0
	for _, node := range oriNodes {
		sorts[i] = node
		i++
	}
	
	// write nodes order by id 
	sort.Slice(sorts, func(i, j int) bool {
		return sorts[i].ID < sorts[j].ID
	})

    // at most write 8000 nodes in one batch
	nodes := make([]*gosm.Node, 0, gosmWriteElementsMax)
	count := 0
	for _, n := range sorts {
		nodes = append(nodes, n.ToGosmNode())
		count++
		if count == gosmWriteElementsMax {
			encoder.AppendNodes(nodes)
			nodes = make([]*gosm.Node, 0, gosmWriteElementsMax)
			count = 0
		}
	}
	if len(nodes) > 0 {
		encoder.AppendNodes(nodes)
	}
	// remember to flush
	encoder.Flush(gosm.NodeType)

    // write ways order by way id
	sorts := make([]*Way, len(oriWays))
	i := 0
	for _, way := range oriWays {
		sorts[i] = way
		i++
	}
	sort.Slice(sorts, func(i, j int) bool {
		return sorts[i].ID < sorts[j].ID
	})

	ways := make([]*gosm.Way, 0, gosmWriteElementsMax)
	count := 0
	for _, w := range sorts {
		ways = append(ways, w.ToGosmWay())
		count++
		if count == gosmWriteElementsMax {
			encoder.AppendWays(ways)
			ways = make([]*gosm.Way, 0, gosmWriteElementsMax)
			count = 0
		}
	}
	if len(ways) > 0 {
		encoder.AppendWays(ways)
	}
	encoder.Flush(gosm.WayType)
	
	// write relations by relation id
	sorts := make([]*Relation, len(oriRelations))
	i := 0
	for _, rel := range oriRelations {
		sorts[i] = rel
		i++
	}
	sort.Slice(sorts, func(i, j int) bool {
		return sorts[i].ID < sorts[j].ID
	})

	relations := make([]*gosm.Relation, 0, gosmWriteElementsMax)
	count := 0
	for _, r := range sorts {
		relations = append(relations, r.ToGosmRelation())
		count++
		if count == gosmWriteElementsMax {
			encoder.AppendRelations(relations)
			relations = make([]*gosm.Relation, 0, gosmWriteElementsMax)
			count = 0
		}
	}
	if len(relations) > 0 {
		encoder.AppendRelations(relations)
	}

	if len(errs) > 0 {
		...
	}
```
# Notes
1. When you use `AppendNodes`, `AppendWays` or `AppendRelations`, at most 8000 items are allowed in one append operation
2. Use `encoder.Flush(memberType MemberType)` when you finished writing of one member type.
