package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/delaunay"
	"github.com/koron/clustering-examination/internal/graphmatrix"
	"github.com/koron/clustering-examination/internal/loader"
	"github.com/koron/clustering-examination/internal/wardsmethod"
	"gonum.org/v1/gonum/graph/coloring"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func writeTree(name string, tree wardsmethod.Tree) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	defer b.Flush()
	return wardsmethod.DumpTree(b, tree)
}

func statistics(label string, nodes []wardsmethod.Node, alives []int) {
	num := float64(len(alives))
	var sumW, sumD float64
	for _, v := range alives {
		n := nodes[v]
		sumW += n.Weight
		sumD += n.Delta
	}
	meanW := sumW / num
	meanD := sumD / num
	var varW, varD float64
	for _, v := range alives {
		n := nodes[v]
		varW += math.Pow(n.Weight-meanW, 2)
		varD += math.Pow(n.Delta-meanD, 2)
	}
	varW /= num
	varD /= num
	fmt.Printf("%6s: weight=%e±%e delta=%e±%e sumW=%e\n", label, meanW, varW, meanD, varD, sumW)
}

func drawHist(name, title string, data []float64, n int) (*plotter.Histogram, error) {
	h, err := plotter.NewHist(plotter.Values(data), n)
	if err != nil {
		return nil, err
	}
	p := plot.New()
	p.Title.Text = title
	p.Add(h)
	err = p.Save(1024, 1024, name)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func wards(name string) error {
	pp, err := loader.LoadTSVFile(name)
	if err != nil {
		return err
	}

	var (
		midTree   wardsmethod.Tree
		midAlives []int
	)

	start := time.Now()
	tree := wardsmethod.Clustering(pp, wardsmethod.MonitorFunc(func(nodes []wardsmethod.Node, alives []int) {
		if len(alives) == 45 {
			midTree = wardsmethod.Tree(nodes)
			midAlives = make([]int, len(alives))
			copy(midAlives, alives)
		}
	}))
	log.Printf("Clustering elapsed %s, len(nodes)=%d", time.Since(start), len(tree))

	statistics("mid", midTree, midAlives)
	tops := wardsmethod.Top(tree, 45)
	statistics("last", tree, tops)
	means := wardsmethod.Mean(tree, 45)
	statistics("mean", tree, means)

	start = time.Now()
	tri, err := delaunay.Triangulate(collectCenter(tree, means))
	if err != nil {
		return err
	}
	log.Printf("Triangulate elapsed %s, points#%d convexHull#%d triangles#%d halfedges=%d", time.Since(start), len(tri.Points), len(tri.ConvexHull), len(tri.Triangles), len(tri.Halfedges))

	start = time.Now()
	ug := graphmatrix.NewUndirect(len(means))
	for i := 0; i+2 < len(tri.Triangles); i += 3 {
		ug.Set(tri.Triangles[i+0], tri.Triangles[i+1], true)
		ug.Set(tri.Triangles[i+0], tri.Triangles[i+2], true)
		ug.Set(tri.Triangles[i+1], tri.Triangles[i+2], true)
	}
	k, colors, err := coloring.DsaturExact(nil, ug)
	if err != nil {
		return err
	}
	log.Printf("DsaturExact elapsed %s, k=%d colors=%+v\n", time.Since(start), k, colors)

	// draw histograms of weight, delta of nodes
	//drawHist("tmp/ward-mid-weigts.png", "Middle Weight", midTree.Weights(midAlives), 10)
	//drawHist("tmp/ward-mid-deltas.png", "Middle Delta", midTree.Deltas(midAlives), 10)
	//drawHist("tmp/ward-last-weigts.png", "Last Weigt", tree.Weights(tops), 10)
	//drawHist("tmp/ward-last-deltas.png", "Last Delta", tree.Deltas(tops), 10)
	//drawHist("tmp/ward-mean-weigts.png", "Mean Weigt", tree.Weights(means), 10)
	//drawHist("tmp/ward-mean-deltas.png", "Mean Delta", tree.Deltas(means), 10)

	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, midTree, midAlives)
	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, tree, tops)
	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, tree, means)

	// enumerate points of clusters
	clustersPoints := make([]plotter.XYs, len(means))
	points := make(plotter.XYs, len(pp))
	for i, nn := range means {
		clustersPoints[i] = vectors(points[:0], tree, nn)
		points = points[len(clustersPoints[i]):]
		//fmt.Printf("#%d w=%d/%d %+v\n", nn, int(tree[nn].Weight), len(clustersPoints[i]), clustersPoints[i])
	}

	// draw clusters with colors
	p := plot.New()
	for i, points := range clustersPoints {
		sc, err := plotter.NewScatter(points)
		if err != nil {
			return err
		}
		sc.Color = plotutil.DarkColors[colors[int64(i)]]
		sc.Radius = 1.5
		p.Add(sc)
	}
	return p.Save(1000, 1000, "tmp/out.png")

	//err = drawClusters("tmp/means", clustersPoints, tree, means)
	//if err != nil {
	//	return err
	//}

	//err = writeTree("tmp/tree.txt", tree)
	//if err != nil {
	//	return err
	//}

	return nil
}

func drawClusters(dirname string, clustersPoints []plotter.XYs, tree wardsmethod.Tree, indexes []int) error {
	err := os.MkdirAll(dirname, 0777)
	if err != nil {
		return err
	}
	for i, idx := range indexes {
		n := tree[idx]
		name := filepath.Join(dirname, fmt.Sprintf("N%05d-W%d-D%f.png", idx, int(n.Weight), n.Delta))
		title := fmt.Sprintf("#%d/%d node#%05d weight=%d delta=%f", i+1, len(indexes), idx, int(n.Weight), n.Delta)
		err := drawCluster(name, title, clustersPoints, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func drawCluster(name, title string, clustersPoints []plotter.XYs, target int) error {
	p := plot.New()
	p.Title.Text = title
	// add other scatters.
	for i, points := range clustersPoints {
		if i == target {
			continue
		}
		sc, err := plotter.NewScatter(points)
		if err != nil {
			return err
		}
		sc.Color = color.Gray16{0x9999}
		sc.Radius = 1.5
		p.Add(sc)
	}
	// add a target scatter.
	sc, err := plotter.NewScatter(clustersPoints[target])
	if err != nil {
		return err
	}
	sc.Color = plotutil.Color(0)
	p.Add(sc)
	// plot
	return p.Save(1000, 1000, name)
}

func collectCenter(tree wardsmethod.Tree, indexes []int) []delaunay.Point {
	points := make([]delaunay.Point, 0, len(indexes))
	for _, idx := range indexes {
		n := tree[idx]
		points = append(points, delaunay.Point(n.Center))
	}
	return points
}

func vectors(vecs plotter.XYs, tr wardsmethod.Tree, nodeNum int) plotter.XYs {
	tr.ForEach(nodeNum, func(node wardsmethod.Node) {
		if node.Left < 0 && node.Right < 0 {
			vecs = append(vecs, plotter.XY(node.Center))
		}
	})
	return vecs
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("require a TSV file")
	}
	err := wards(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
}
