/*
 * HR management Api
 *
 * Ambulance HR management API
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package ambulance_hr

type PersonalDocument struct {

	// Unique identifier of the document
	Id string `json:"id"`

	// Name of the document
	Name string `json:"name"`

	// Content of the document
	Content string `json:"content"`
}
