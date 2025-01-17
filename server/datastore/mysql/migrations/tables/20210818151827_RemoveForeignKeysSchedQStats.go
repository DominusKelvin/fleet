package tables

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20210818151827, Down_20210818151827)
}

func Up_20210818151827(tx *sql.Tx) error {
	referencedTables := map[string]struct{}{"hosts": {}, "scheduled_queries": {}}
	table := "scheduled_query_stats"

	constraints, err := constraintsForTable(tx, table, referencedTables)
	if err != nil {
		return err
	}

	if len(constraints) == 0 {
		return errors.New("Found no constraints in scheduled_query_stats")
	}

	for _, constraint := range constraints {
		_, err = tx.Exec(fmt.Sprintf(`ALTER TABLE scheduled_query_stats DROP FOREIGN KEY %s;`, constraint))
		if err != nil {
			return errors.Wrapf(err, "dropping fk %s", constraint)
		}
	}
	return nil
}

func constraintsForTable(tx *sql.Tx, table string, referencedTables map[string]struct{}) ([]string, error) {
	var constraints []string
	query := `SELECT DISTINCT CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME = ? AND CONSTRAINT_NAME <> 'PRIMARY'`
	rows, err := tx.Query(query, table) //nolint
	if err != nil {
		return nil, errors.Wrap(err, "getting fk for scheduled_query_stats")
	}
	for rows.Next() {
		var constraintName string
		var referencedTable string
		err := rows.Scan(&constraintName, &referencedTable)
		if err != nil {
			return nil, errors.Wrap(err, "scanning fk for scheduled_query_stats")
		}
		if _, ok := referencedTables[referencedTable]; ok {
			constraints = append(constraints, constraintName)
		}
	}
	return constraints, nil
}

func Down_20210818151827(tx *sql.Tx) error {
	return nil
}
