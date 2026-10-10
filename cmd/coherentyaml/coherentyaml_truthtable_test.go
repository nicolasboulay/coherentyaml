package main

import "testing"

// TestCoherenceTruthTable specifies coherence for pairs of YAML documents.
// Only A.IsCoherentWith(B) is checked here: symmetry and the other algebraic
// laws belong in separate tests. A true expectation means a common value exists.
//
// Neutral values describe domains rather than exact values: 1 denotes positive
// integers, -1 all integers, 1.0 floating-point values, and "" strings.
// Pending rules are skipped explicitly; their wantCoherent value is unused.
func TestCoherenceTruthTable(t *testing.T) {
	tests := []struct {
		name         string
		a, b         string
		wantCoherent bool
		pending      string
	}{
		// Integers: define whether the positive domain includes zero separately.
		// 1 define the N domain
		// -1 define the Z domain
		{"positive integers", "1", "42", true, ""},
		{"positive integer and zero", "1", "0", true, ""},
		{"positive and negative integers", "1", "-4", false, ""},
		{"all integers and negative integer", "-1", "-4", true, ""},
		{"all integers and positive integer", "-1", "42", true, ""},
		{"all integers and zero", "-1", "0", true, ""},

		// Integer and floating-point domains remain distinct.
		{"floating-point values", "1.0", "3.14", true, ""},
		{"integer and floating-point value", "1", "1.0", false, ""},

		// The empty string is neutral only within the string domain.
		{"strings", `""`, `"abc"`, true, ""},
		{"string and integer", `""`, "12", false, ""},

		// These domains still need an explicit decision.
		{"boolean values", "true", "false", false, "règle à définir : booléens exacts ou valeur neutre pour le type booléen ?"},
		{"null values", "null", "null", false, "règle à définir : représenter null sans le confondre avec un nœud absent"},

		// Objects require matching field names and coherent field values.
		{"matching object fields", "{a: 1}", "{a: 5}", true, ""},
		{"different object fields", "{a: 1}", "{b: 5}", false, ""},

		// Array semantics (order, length, and multiplicity) are still open.
		{"integer arrays", "[1]", "[2,3]", false, "règle à définir : taille et multiplicité des tableaux"},
		{"mixed arrays", `[1,""]`, `[4,"x"]`, true, "règle à définir : correspondance des éléments des tableaux non ordonnés"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pending != "" {
				t.Skip(tt.pending)
			}
			a, b := makeNode(tt.a), makeNode(tt.b)
			if a == nil || b == nil {
				t.Fatalf("conversion YAML non prise en charge : A=%s, B=%s", tt.a, tt.b)
			}
			err := a.IsCoherentWith(b)
			if got := err == nil; got != tt.wantCoherent {
				t.Errorf("coherent(%s, %s) = %t, attendu %t (erreur : %v)", tt.a, tt.b, got, tt.wantCoherent, err)
			}
		})
	}
}
