package main

import (
	"fmt"
	"os/exec"
	"time"
)

// backupMongoDB tạo bản backup cho các collection cụ thể trong cơ sở dữ liệu MongoDB.
func backupMongoDB(uri, dbName string, collections []string, backupDir string) error {
	// Tạo timestamp để sử dụng chung cho tất cả các collection trong lần backup này
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	baseBackupPath := fmt.Sprintf("%s/%s-%s", backupDir, dbName, timestamp)

	for _, collectionName := range collections {
		// Tạo thư mục cho từng collection
		backupPath := fmt.Sprintf("%s/%s", baseBackupPath, collectionName)

		// Tạo lệnh mongodump cho mỗi collection
		cmd := exec.Command("mongodump", "--uri="+uri, "--db="+dbName, "--collection="+collectionName, "--out="+backupPath)

		// Chạy lệnh và kiểm tra lỗi
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("backup failed for collection %s: %w", collectionName, err)
		}

		fmt.Printf("Backup completed for collection: %s at %s\n", collectionName, backupPath)
	}
	return nil
}

func main() {
	// Thông tin kết nối MongoDB, tên database, danh sách collection và thư mục backup
	uri := "mongodb://localhost:27017"
	dbName := "my_database"
	collections := []string{"collection1", "collection2", "collection3"}
	backupDir := "/path/to/backup"

	// Gọi hàm backup
	if err := backupMongoDB(uri, dbName, collections, backupDir); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Backup for all collections successful!")
	}
}
////

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// listBackups liệt kê các bản backup trong thư mục backupDir.
func listBackups(backupDir string) ([]string, error) {
	var backups []string

	// Duyệt qua tất cả các file và thư mục con trong backupDir
	err := filepath.WalkDir(backupDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Nếu tìm thấy thư mục chứa bản backup (dựa trên cấu trúc thư mục backup)
		if d.IsDir() && path != backupDir {
			backups = append(backups, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	return backups, nil
}

func main() {
	// Thư mục chứa các bản backup
	backupDir := "/path/to/backup"

	// Gọi hàm listBackups
	backups, err := listBackups(backupDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("List of backups:")
		for _, backup := range backups {
			fmt.Println(backup)
		}
	}
}


////
package main

import (
	"fmt"
	"os/exec"
)

// restoreMongoDB phục hồi dữ liệu từ một bản backup.
func restoreMongoDB(uri, dbName, backupPath, collectionName string) error {
	// Tạo lệnh mongorestore
	var cmd *exec.Cmd
	if collectionName == "" {
		// Phục hồi toàn bộ database nếu không truyền collection
		cmd = exec.Command("mongorestore", "--uri="+uri, "--db="+dbName, backupPath)
	} else {
		// Phục hồi một collection cụ thể
		cmd = exec.Command("mongorestore", "--uri="+uri, "--db="+dbName, "--collection="+collectionName, backupPath+"/"+collectionName+".bson")
	}

	// Chạy lệnh và kiểm tra lỗi
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	fmt.Printf("Restore completed from %s\n", backupPath)
	return nil
}

func main() {
	// Thông tin kết nối MongoDB, tên database, đường dẫn backup, và tên collection
	uri := "mongodb://localhost:27017"
	dbName := "my_database"
	backupPath := "/path/to/backup/my_database-2024-10-27"
	collectionName := "" // Để trống nếu muốn phục hồi toàn bộ database

	// Gọi hàm restore
	if err := restoreMongoDB(uri, dbName, backupPath, collectionName); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Restore successful!")
	}
}




////

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// listBackups liệt kê các bản backup của nhiều database trong thư mục backupDir.
func listBackups(backupDir string) (map[string][]string, error) {
	backups := make(map[string][]string)

	// Duyệt qua các thư mục con trong backupDir
	err := filepath.WalkDir(backupDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Kiểm tra nếu path là thư mục chứa bản backup của một database
		if d.IsDir() && path != backupDir {
			// Lấy tên database từ thư mục con
			dbName := filepath.Base(path)
			
			// Lấy danh sách bản backup trong thư mục của database đó
			var dbBackups []string
			filepath.WalkDir(path, func(subPath string, info fs.DirEntry, err error) error {
				if info.IsDir() && subPath != path {
					dbBackups = append(dbBackups, subPath)
				}
				return nil
			})
			backups[dbName] = dbBackups
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	return backups, nil
}

func main() {
	// Thư mục chứa các bản backup
	backupDir := "/path/to/backup"

	// Gọi hàm listBackups
	backups, err := listBackups(backupDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("List of backups:")
		for dbName, dbBackups := range backups {
			fmt.Printf("Database: %s\n", dbName)
			for _, backup := range dbBackups {
				fmt.Printf("  - %s\n", backup)
			}
		}
	}
}

///
package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// restoreMongoDBs phục hồi các database từ các bản backup trong backupDir.
func restoreMongoDBs(uri string, backupDir string) error {
	// Duyệt qua các thư mục trong backupDir
	err := filepath.WalkDir(backupDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Nếu thư mục con là một database
		if d.IsDir() && path != backupDir {
			dbName := filepath.Base(path)

			// Phục hồi từng database từ thư mục backup
			cmd := exec.Command("mongorestore", "--uri="+uri, "--db="+dbName, path)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("restore failed for database %s: %w", dbName, err)
			}
			fmt.Printf("Restore completed for database: %s from %s\n", dbName, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to restore backups: %w", err)
	}

	return nil
}

func main() {
	// Thông tin kết nối MongoDB và thư mục chứa các bản backup
	uri := "mongodb://localhost:27017"
	backupDir := "/path/to/backup"

	// Gọi hàm restore
	if err := restoreMongoDBs(uri, backupDir); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Restore for all databases successful!")
	}
}


