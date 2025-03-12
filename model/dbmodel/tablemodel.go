package dbmodel

import (
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
)

type IDer interface {
	GetID() int
}
type TableModel[T IDer] struct {
	Rows   []T
	RowMap map[int]*T
}

func (tableModel *TableModel[T]) Init(rows []T) {
	tableModel.Rows = rows
	tableModel.RowMap = make(map[int]*T, len(rows))
	for i := 0; i < len(rows); i++ {
		tableModel.RowMap[rows[i].GetID()] = &rows[i]
	}
}

func (tableModel *TableModel[T]) Append(row T) {
	tableModel.Rows = append(tableModel.Rows, row)
	tableModel.RowMap[row.GetID()] = &tableModel.Rows[len(tableModel.Rows)-1]
}

func (tableModel *TableModel[T]) Get(id int) *T {
	if v, ok := tableModel.RowMap[id]; ok {
		return v
	}
	return nil
}

type DBTabelModel struct {
	table.ABEntityTableModel //AB 表

	FilterResourceTableModel    TableModel[table.FilterResourceMeta]
	FilterEntityTableModel      TableModel[table.FilterEntityMeta]
	FilterGroupEntityTableModel TableModel[table.FilterGroupEntityMeta]

	RecallEntityTableModel      TableModel[table.RecallEntityMeta]      //单路召回表
	RecallGroupEntityTableModel TableModel[table.RecallGroupEntityMeta] //召回组表

	WeightedEntityTableModel    TableModel[table.WeightedEntityMeta]
	RankEntityTableModel        TableModel[table.RankEntityMeta]
	InsertEntityTableModel      TableModel[table.InsertEntityMeta]
	InsertGroupEntityTableModel TableModel[table.InsertGroupEntityMeta]
	ScatterEntityTableModel     TableModel[table.ScatterEntityMeta]

	StrategyEntityTableModel TableModel[table.StrategyEntityMeta]
}
