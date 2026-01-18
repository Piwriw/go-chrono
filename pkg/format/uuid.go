package format

import "github.com/google/uuid"

// ParseJobID parses a jobs ID string into a UUID.
// ParseJobID 将任务 ID 字符串解析为 UUID。
//
// Parameters:
//
//	jobID - The jobs ID string to parse / 要解析的任务 ID 字符串
//
// Returns:
//
//	uuid.UUID - The parsed UUID / 解析后的 UUID
//	error - Error if the parsing fails / 如果解析失败则返回错误
func ParseJobID(jobID string) (uuid.UUID, error) {
	return uuid.Parse(jobID)
}
