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

// checks wheter a string slice contains a string
func contains(slice []string, str string) bool {
    for _, element := range slice {
        if element == str {
            return true
        }
    }
    return false
}

// remove duplications from string slice
func unique(slice []string) []string {
    var uniqueSlice []string
    for _, element := range slice {
        if !contains(uniqueSlice, element) {
            uniqueSlice = append(uniqueSlice, element)
        }
    }
    return uniqueSlice
}
