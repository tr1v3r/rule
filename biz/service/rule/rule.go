package rule

// Query ...
func Query(treeName, nodeName string) string {
	return "/api/v1/rules/" + treeName + "/" + nodeName
}
