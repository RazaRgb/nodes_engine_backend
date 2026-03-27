package db

import (
	"backend/src/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetTreeFromDB(treeID string, tx pgx.Tx) (models.Tree, error) {
	var tree models.Tree
	tree.ID = treeID
	var err error
	tree.Nodes, err = getNodesFromDB(treeID, tx)
	if err != nil {
		return tree, err
	}
	tree.Edges, err = getEdgesFromDB(treeID, tx)
	if err != nil {
		return tree, err
	}

	return tree, nil
}
func getNodesFromDB(treeID string, tx pgx.Tx) ([]models.Node, error) {
	var nodes []models.Node

	selectQuery := `
	SELECT id, type, pos_x, pos_y, label, value FROM nodes WHERE belongs_to = $1
	`

	rows, err := tx.Query(context.Background(), selectQuery, treeID)
	if err != nil {
		fmt.Printf("unable to query tree database for nodes: %+v\n", err)
		return nodes, err
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
			&node.Data.Value,
		)
		if err != nil {
			fmt.Printf("error scanning nodes in get nodes query: %+v\n", err)
			return nodes, err
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func getEdgesFromDB(treeID string, tx pgx.Tx) ([]models.Edge, error) {
	var edges []models.Edge

	selectQuery := `
	SELECT edges.id, edges.source, edges.source_handle, edges.target, edges.target_handle 
	FROM edges
	INNER JOIN nodes ON edges.source = nodes.id
	WHERE nodes.belongs_to = $1
	`

	rows, err := tx.Query(context.Background(), selectQuery, treeID)
	if err != nil {
		fmt.Printf("unable to query tree database for edges: %+v\n", err)
		return edges, err
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
			return edges, err
		}

		edges = append(edges, edge)
	}
	return edges, nil
}

// ------------------------------------------

func ClearTreeContent(treeID string, tx pgx.Tx) error {
	clearQuery := `DELETE FROM nodes WHERE belongs_to = $1`
	_, err := tx.Exec(
		context.Background(),
		clearQuery,
		treeID,
	)
	if err != nil {
		fmt.Printf("unable to clear tree: %+v\n", err)
		return err
	}
	return nil
}

func InsertTreeContentInDB(updatedTree models.Tree, tx pgx.Tx) error {
	edgeArray := updatedTree.Edges
	nodeArray := updatedTree.Nodes

	for _, node := range nodeArray {
		err := saveNodeToDB(node, updatedTree.ID, tx)
		if err != nil {
			return err
		}
	}
	for _, edge := range edgeArray {
		err := saveEdgesToDB(edge, updatedTree.ID, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveNodeToDB(node models.Node, treeID string, tx pgx.Tx) error {
	insertNodeQuery := `
	INSERT INTO nodes (belongs_to, id, type, pos_x, pos_y, label, value)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(
		context.Background(),
		insertNodeQuery,
		treeID,
		node.ID,
		node.Type,
		node.Pos.X,
		node.Pos.Y,
		node.Data.Label,
		node.Data.Value,
	)
	if err != nil {
		fmt.Printf("Failed to execute node INSERT \n %+v\n", err)
		return err
	}
	return nil
}

func saveEdgesToDB(edge models.Edge, treeID string, tx pgx.Tx) error {
	insertQuery := `
	INSERT INTO edges (id, belongs_to, source, source_handle, target, target_handle)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(
		context.Background(),
		insertQuery,
		edge.ID,
		treeID,
		edge.Source,
		edge.SourceHandle,
		edge.Target,
		edge.TargetHandle,
	)
	if err != nil {
		fmt.Printf("Failed to execute edge INSERT \n %+v\n", err)
		return err
	}
	return nil
}

func CreateNewTree(newTree models.Tree, projID string, tx pgx.Tx) error {
	createQuery := `
	INSERT INTO trees (id, project_id)
	VALUES ($1, $2)
	`
	_, err := tx.Exec(
		context.Background(),
		createQuery,
		newTree.ID,
		projID,
	)
	if err != nil {
		fmt.Printf("unable to create a new tree: %+v\n", err)
		return err
	}

	edgeArray := newTree.Edges
	nodeArray := newTree.Nodes

	for _, node := range nodeArray {
		err := saveNodeToDB(node, newTree.ID, tx)
		if err != nil {
			return err
		}
	}
	for _, edge := range edgeArray {
		err := saveEdgesToDB(edge, newTree.ID, tx)
		if err != nil {
			return err
		}
	}
	return nil
}
