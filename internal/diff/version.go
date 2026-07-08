package diff

import "github.com/martyn/apidiff/internal/model"

// VersionBump represents a semantic version bump recommendation.
type VersionBump string

const (
	BumpNone  VersionBump = "none"
	BumpPatch VersionBump = "patch"
	BumpMinor VersionBump = "minor"
	BumpMajor VersionBump = "major"
)

func (v VersionBump) String() string {
	return string(v)
}

// RecommendVersion analyzes the list of changes and returns the highest
// semantic version bump required.
//
// Rules:
//   - Removed endpoint        → Major
//   - Required field added    → Major
//   - Property removed        → Major
//   - Property type changed   → Major
//   - New optional endpoint   → Minor
//   - New optional property   → Minor
//   - Documentation only      → Patch
//   - No changes              → None
func RecommendVersion(changes []model.Change) VersionBump {
	if len(changes) == 0 {
		return BumpNone
	}

	bump := BumpPatch

	for _, c := range changes {
		switch {
		case c.Type == model.EndpointRemoved:
			return BumpMajor
		case c.Type == model.PropertyRemoved:
			return BumpMajor
		case c.Type == model.PropertyTypeChanged:
			return BumpMajor
		case c.Type == model.PropertyAdded && c.Severity == model.SeverityBreaking:
			// Required field added
			return BumpMajor
		case c.Type == model.EndpointAdded:
			bump = BumpMinor
		case c.Type == model.PropertyAdded && c.Severity != model.SeverityBreaking:
			bump = BumpMinor
		case c.Type == model.SchemaChanged:
			if bump != BumpMinor {
				bump = BumpMinor
			}
		}
	}

	return bump
}
