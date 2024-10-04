package chart

import (
	"github.com/beto20/gofluence/model"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Chart struct{}

const (
	LAYOUT         = "orthogonal"
	ORIENTATION    = "LR"
	TREE_DEPTH     = -1
	POSITION_RIGHT = "right"
	POSITION_TOP   = "top"
	COLOR_BLACK    = "black"
)

func GenerateTreeChart(documents []model.Document, projectName string) *charts.Tree {
	graph := charts.NewTree()
	graph.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "80vh", BackgroundColor: "white"}),
		charts.WithTitleOpts(opts.Title{Title: projectName}),
	)

	tree := buildTree(documents, projectName)
	graph.AddSeries(projectName, tree).SetSeriesOptions(
		charts.WithTreeOpts(
			opts.TreeChart{
				Layout:           LAYOUT,
				Orient:           ORIENTATION,
				InitialTreeDepth: TREE_DEPTH,
				Leaves: &opts.TreeLeaves{
					Label: &opts.Label{
						Show:     opts.Bool(true),
						Position: POSITION_RIGHT,
						Color:    COLOR_BLACK,
					},
					LineStyle: &opts.LineStyle{
						Type: "dashed",
					},
				},
				Right:  "10%",
				Left:   "6%",
				Bottom: "2%",
				Top:    "5%",
			},
		),
		charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: POSITION_TOP, Color: COLOR_BLACK}),
	)

	return graph
}

func GenerateProjectChart(documents []model.Document) map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for i := 0; len(documents) > i; i++ {
		var readyRead []string
		var depArr []string

		mapDep := make(map[string][]string)
		for j := 0; len(documents[i].Dependencies.Dependency) > j; j++ {
			depArr = nil

			if !(contains(readyRead, documents[i].Dependencies.Dependency[j].GroupId)) {
				for k := 0; len(documents[i].Dependencies.Dependency) > k; k++ {
					if documents[i].Dependencies.Dependency[j].GroupId == documents[i].Dependencies.Dependency[k].GroupId {
						depArr = append(depArr, documents[i].Dependencies.Dependency[k].ArtifactId)
					}
				}
			}

			if len(depArr) != 0 {
				mapDep[documents[i].Dependencies.Dependency[j].GroupId] = depArr
			}

			readyRead = append(readyRead, documents[i].Dependencies.Dependency[j].GroupId)
		}
		result[documents[i].Name] = mapDep
	}

	return result
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func buildTree(documents []model.Document, nodeName string) []opts.TreeData {
	results := GenerateProjectChart(documents)
	var nArr []node

	for _, d := range documents {
		groupDeps := results[d.Name]
		// fmt.Println("groupDeps: ", d.Name, "-", groupDeps)
		var chArr []children
		for key, _ := range groupDeps {
			scArr := []subChildren{}
			for _, v := range groupDeps[key] {
				sub := subChildren{
					name: v,
				}
				scArr = append(scArr, sub)
			}

			ch := children{
				name:        key,
				subChildren: scArr,
			}

			chArr = append(chArr, ch)

			// fmt.Println("CHILDREN: ", ch)
		}

		n := node{
			nodeName: d.Name,
			children: chArr,
		}

		nArr = append(nArr, n)
		// fmt.Println("NODE: ", n)
	}

	sn := superNode{
		node: nArr,
	}

	tree := toTreeData4(sn, nodeName)

	return tree

}

func toTreeData4(n superNode, name string) []opts.TreeData {
	var res []opts.TreeData
	x := opts.TreeData{
		Name:     name,
		Children: toTreeData3(n.node),
	}
	res = append(res, x)
	return res
}

func toTreeData3(node []node) []*opts.TreeData {
	var res []*opts.TreeData
	for _, n := range node {
		x := &opts.TreeData{
			Name:     n.nodeName,
			Children: toTreeData2(n.children),
		}
		res = append(res, x)
	}

	return res
}

func toTreeData2(children []children) []*opts.TreeData {
	var res []*opts.TreeData
	for _, c := range children {
		x := &opts.TreeData{
			Name:     c.name,
			Children: toTreeData(c.subChildren),
		}
		res = append(res, x)
	}

	return res
}

func toTreeData(subChildren []subChildren) []*opts.TreeData {
	var res []*opts.TreeData
	for _, sc := range subChildren {
		x := &opts.TreeData{
			Name: sc.name,
		}
		res = append(res, x)
	}

	return res
}

type superNode struct {
	node []node
}

type node struct {
	nodeName string
	children []children
}

type children struct {
	name        string
	subChildren []subChildren
}

type subChildren struct {
	name string
}
