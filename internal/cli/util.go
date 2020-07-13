package cli

import "strconv"

// convert int slice to string slice
func stringify(slice []int) []string {
    items := make([]string, len(slice))
    for index, e := range slice {
        items[index] = strconv.FormatInt(int64(e), 10)
    }
    return items
}

// surround every element in slice with prefix and suffix
func surround(slice []string, prefix, suffix string) []string {
    items := make([]string, len(slice))
    for index, e := range slice {
        items[index] = prefix + e + suffix
    }
    return items
}
