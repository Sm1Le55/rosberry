/*
 * Rosberry Mobile APP API (Test work)
 *
 * Test work for rosberry
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DisplaySettings struct {

	UserId int32 `json:"userId,omitempty"`

	ShowMeAges string `json:"showMeAges,omitempty"`

	HideMeFromAges string `json:"hideMeFromAges,omitempty"`

	ShowThemes []string `json:"showThemes,omitempty"`

	HideThemes []string `json:"hideThemes,omitempty"`

	Location string `json:"location,omitempty"`
}
