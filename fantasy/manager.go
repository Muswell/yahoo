// Fantasy wraps the Yahoo fantasy API into Go types and provides functions to generate proper request URLs.
package fantasy

// Manager type represents a single Yahoo fantasy team manager.
type Manager struct {
	// Guid is the unique ID of the user
	Guid string
	// ManagerID is the id of the manager within the League.
	ManagerID string `xml:"manager_id"`
	// Name is the nickname the manager is using within the League.
	Name string `xml:"nickname"`
	// Email is the email address the manager is using within the League.
	Email string
	// ImageURL is the address of the manager's avatar
	ImageURL string `xml:"image_url"`
	// IsCurrentLogin is a bool value indicating if this manager is the logged in user
	IsActiveUser IntAsBool `xml:"is_current_login"`
}
