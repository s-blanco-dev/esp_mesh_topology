package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"mesh_topology/models"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/gin-gonic/gin"
)

var ROOT_NODE models.Node
var NUM_NODES int = 5
var BUFFER []models.Node

func PostMeshData(c *gin.Context) {
	var node models.Node
	if err := c.BindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* MIENTRAS NO LLEGUEN TODOS LOS PAQUETES:
	- No genero el grafo
	- Guardo nodo en buffer y espero por los restantes
	*/
	if node.IsRoot {
		ROOT_NODE = node
	}

	BUFFER = append(BUFFER, ROOT_NODE)
	BUFFER = append(BUFFER, node)

	if len(BUFFER) < NUM_NODES {
		status := fmt.Sprintf("Nodo agregado a buffer, esperando por %d paquetes más.", NUM_NODES-len(BUFFER))
		c.JSON(http.StatusOK, gin.H{"status": status, "nodes": len(BUFFER)})
		return
	}

	// Creamos el grafo
	g := graph.New(graph.StringHash, graph.Rooted())

	// Agrega nodos
	for _, n := range BUFFER {
		_ = g.AddVertex(n.SelfMAC, graph.VertexAttribute("label",
			fmt.Sprintf("%s\nT: %.1f°C\nH: %.1f%%", n.SelfMAC, n.Temp, n.Humidity)))
	}

	// Agrega aristas
	for _, n := range BUFFER {
		if !n.IsRoot {
			_ = g.AddEdge(n.ParentMAC, n.SelfMAC)
		}
	}

	// Guarda el grafo en formato DOT
	f, err := os.Create("static/topology.dot")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	_ = draw.DOT(g, f)

	// Genera el PNG usando Graphviz
	_ = exec.Command("dot", "-Tpng", "static/topology.dot", "-o", "static/topology.png").Run()

	c.JSON(http.StatusOK, gin.H{"status": "graph updated", "nodes": len(BUFFER)})
	BUFFER = nil
}
