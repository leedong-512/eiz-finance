package tools

import (
	"fmt"
)

func Del(exType int, userId int)  {
	//taskDb, _ := database.GetTaskDB()
	sql := fmt.Sprintf(`DELETE FROM xero_data where type = %d`, exType)
	_, err := taskDb.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
}