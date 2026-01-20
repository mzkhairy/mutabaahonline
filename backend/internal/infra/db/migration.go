package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migrate menjalankan semua file *.up.sql di folder target
func (d *DB) Migrate(dir string) error {
	// 1. Cek folder migrasi
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("migration directory not found: %s", dir)
	}

	// 2. Baca file di folder
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	// 3. Filter hanya file .up.sql
	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}

	// 4. Sort agar urutan eksekusi benar (0001 dulu, baru 0002)
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		return fmt.Errorf("no migration files found in %s", dir)
	}

	// 5. Eksekusi SQL satu per satu
	for _, file := range sqlFiles {
		path := filepath.Join(dir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Eksekusi query
		_, err = d.DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
		fmt.Printf("[INFO] Migration executed: %s\n", file)
	}

	return nil
}
