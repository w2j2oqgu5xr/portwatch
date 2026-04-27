// Package portclassify provides risk-tier classification for observed ports.
//
// A Classifier is constructed with a Policy that maps port numbers to one of
// three tiers: Safe, Caution, or Critical. Ports with no policy match are
// returned as Unknown.
//
// Usage:
//
//	policy := portclassify.DefaultPolicy()
//	classifier := portclassify.New(policy)
//	result := classifier.Classify(22)
//	fmt.Println(result) // port 22: caution (sensitive service)
package portclassify
