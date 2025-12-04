package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"mesh_topology/models"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/gin-gonic/gin"
)

var NUM_NODES int = 6
var BUFFER []models.Node

func PostMeshData(c *gin.Context) {
	var node models.Node
	if err := c.BindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	BUFFER = append(BUFFER, node)

	if len(BUFFER) < NUM_NODES {
		status := fmt.Sprintf("Nodo agregado a buffer, esperando por %d paquetes más.", NUM_NODES-len(BUFFER))
		c.JSON(http.StatusOK, gin.H{"status": status, "nodes": len(BUFFER)})
		return
	}

	// ---------------------------------------------------------------
	// REORDENO BUFFER para que el root esté primero
	// ---------------------------------------------------------------
	var reordered []models.Node
	var root models.Node
	foundRoot := false

	for _, n := range BUFFER {
		if n.IsRoot {
			root = n
			foundRoot = true
			break
		}
	}

	if foundRoot {
		reordered = append(reordered, root)
		for _, n := range BUFFER {
			if !n.IsRoot {
				reordered = append(reordered, n)
			}
		}
	} else {
		reordered = BUFFER
	}

	BUFFER = reordered
	// ---------------------------------------------------------------

	g := graph.New(graph.StringHash, graph.Tree())

	// agrega nodos
	for _, n := range BUFFER {
		if n.IsRoot {
			_ = g.AddVertex(n.SelfMAC,
				graph.VertexAttribute("label", fmt.Sprintf("%s\n: ROOT", n.SelfMAC)),
				graph.VertexAttribute("fillcolor", "lightcoral"),
				graph.VertexAttribute("color", "red"),
				graph.VertexAttribute("penwidth", "2"),
				graph.VertexAttribute("style", "filled"),
				graph.VertexAttribute("shape", "doublecircle"),
			)
		} else {
			_ = g.AddVertex(n.SelfMAC,
				graph.VertexAttribute("label", fmt.Sprintf("%s\nT: %.1f°C\nH: %.1f%%", n.SelfMAC, n.Temp, n.Humidity)),
				graph.VertexAttribute("shape", "circle"),
			)
		}
	}

	// agrega aristas
	for _, n := range BUFFER {
		if !n.IsRoot {
			_ = g.AddEdge(n.ParentMAC, n.SelfMAC)
		}
	}

	// ------------- Crear archivo DOT --------------
	dotFile := "static/topology.dot"
	f, err := os.Create(dotFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = draw.DOT(g, f)
	f.Close()

	// ---------------------------------------------------------------
	// inserto rank=min para asegurar que el root esté arriba (marcha mao meno)
	// ---------------------------------------------------------------
	rootMAC := BUFFER[0].SelfMAC // porque lo pusimos primero

	contentBytes, _ := os.ReadFile(dotFile)
	content := string(contentBytes)

	// luego de la primera línea "digraph G {", insertamos el rank
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "{") {
			lines = append(lines[:i+1],
				append([]string{fmt.Sprintf("    { rank=min; \"%s\"; }", rootMAC)}, lines[i+1:]...)...)
			break
		}
	}

	_ = os.WriteFile(dotFile, []byte(strings.Join(lines, "\n")), 0644)

	// ---------------------------------------------------------------

	_ = exec.Command("dot", "-Tpng", dotFile, "-o", "static/topology.png").Run()

	c.JSON(http.StatusOK, gin.H{"status": "graph updated", "nodes": len(BUFFER)})
	BUFFER = nil
}
