package cache

import "strconv"

const TasksKeyVersion = "v1"

// KeyUserTasks 生成 user 的 tasks 快取 key：user:<uid>:tasks:v1
func KeyUserTasks(userID uint) string {
	return "user:" + strconv.FormatUint(uint64(userID), 10) + ":tasks:" + TasksKeyVersion
}
