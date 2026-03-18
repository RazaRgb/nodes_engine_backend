package db

import (
	"backend/src/models"
	"context"
	"fmt"
)

func GetTreeFromDB(treeID string) models.Tree {
	var tree models.Tree
	tree.Nodes = getNodesFromDB(treeID)
	tree.Edges = getEdgesFromDB(treeID)
	return tree
}

func getNodesFromDB(treeID string) []models.Node {
	var nodes []models.Node

	selectQuery := `
	SELECT id, type, pos_x, pos_y, label FROM nodes WHERE belongs_to = $1
	`

	rows, err := DB.Query(context.Background(), selectQuery, treeID)
	if err != nil {
		fmt.Printf("unable to query tree database for nodes: %+v\n", err)
		fmt.Printf("query : %+v", selectQuery)
	}

	defer rows.Close()

	for rows.Next() {
		var node models.Node

		err := rows.Scan(
			&node.ID,
			&node.Type,
			&node.Pos.X,
			&node.Pos.Y,
			&node.Data.Label,
		)
		if err != nil {
			fmt.Printf("error scanning nodes in get nodes query: %+v\n", err)
		}

		nodes = append(nodes, node)
	}
	return nodes
}

func getEdgesFromDB(treeID string) []models.Edge {
	var edges []models.Edge

	selectQuery := `
	SELECT edges.id, edges.source, edges.source_handle, edges.target, edges.target_handle 
	FROM edges
	INNER JOIN nodes ON edges.source = nodes.id
	WHERE nodes.belongs_to = $1
	`

	rows, err := DB.Query(context.Background(), selectQuery, treeID)
	if err != nil {
		fmt.Printf("unable to query tree database for edges: %+v\n", err)
		fmt.Printf("query : %+v", selectQuery)
	}

	defer rows.Close()

	for rows.Next() {
		var edge models.Edge

		err := rows.Scan(
			&edge.ID,
			&edge.Source,
			&edge.SourceHandle,
			&edge.Target,
			&edge.TargetHandle,
		)
		if err != nil {
			fmt.Printf("error scanning edges in get edges query: %+v\n", err)
		}

		edges = append(edges, edge)
	}
	return edges
}

func ClearTreeContent(treeID string) {
	clearQuery := `DELETE FROM nodes WHERE belongs_to = $1`
	_, err := DB.Exec(
		context.Background(),
		clearQuery,
		treeID,
	)
	if err != nil {
		fmt.Printf("unable to clear tree for updation: %+v\n", err)
	}
}

func InsertTreeInDB(updatedTree models.Tree) {
	edgeArray := updatedTree.Edges
	nodeArray := updatedTree.Nodes

	for _, node := range nodeArray {
		saveNodeToDB(node, updatedTree.ID)
	}
	for _, edge := range edgeArray {
		saveEdgesToDB(edge)
	}
}

func saveNodeToDB(node models.Node, treeID string) {
	insertNodeQuery := `
	INSERT INTO nodes (belongs_to, id, type, pos_x, pos_y, label)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := DB.Exec(
		context.Background(),
		insertNodeQuery,
		treeID,
		node.ID,
		node.Type,
		node.Pos.X,
		node.Pos.Y,
		node.Data.Label,
	)
	if err != nil {
		fmt.Printf("Failed to execute node INSERT \n %+v\n", err)
	}
}

func saveEdgesToDB(edge models.Edge) {
	insertQuery := `
	INSERT INTO edges (id, source, source_handle, target, target_handle)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := DB.Exec(
		context.Background(),
		insertQuery,
		edge.ID,
		edge.Source,
		edge.SourceHandle,
		edge.Target,
		edge.TargetHandle,
	)
	if err != nil {
		fmt.Printf("Failed to execute edge INSERT \n %+v\n", err)
	}
}
