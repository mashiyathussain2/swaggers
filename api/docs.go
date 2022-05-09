//  CatalogV2 API:
//
//
// # Error
// The API uses standard HTTP status codes to indicate the success or failure of the API call.
// The body of the response will be JSON in the following format:
//
// ```json
//  {
//	"error": [
//		 {
//			"message": "string",
//			"type": "NoType"
//		}
//	],
//	"success": false,
//	"request_id": "string"
// }
//  ```
//
// # Responses
//  __Possible Response Status Codes__
// | Status Code | Description |
//  |-----------|-------|
//  | 200 | OK |
//  | 400 | Bad Request |
//  | 401 | Unauthorized |
//  | 500 | Server Error |
//  | 403 | Invalid User |
//
//
// Terms Of Service:
//         https://www.hypd.store/terms-and-conditions
//
//
//  Schemes: https
//  Host: catalogv2.getshitdone.in
//  BasePath: /api
//  Contact: tech<mashiyat.hussain@hypd.in>
//  License: MIT http://opensource.org/licenses/MIT
//
//
//   Extensions:
//     x-tagGroups:
//       - name: CatalogV2
//         description: This is the catalogv2 api.
//         tags:
//           - AppCatalog
//           - AppCollectionCatalogV2
//           - AppCategoryCatalog
//           - InfluencerCollectionKEEPER
//           - InfluencerCollectionApp
//           - AppReview
//           - AppSale
//           - Inventory
//           - LegacySearch
//           - UnicommerceAPIs
//         x-traitTag: true
//
//
//
//
//  InfoExtensions:
//     x-logo:
//        url: https://www.hypd.store/img/social-logo.png
//        altText: HYPD STORE
//        backgroundColor: "#FFFFFF"
//
//
//
//
//
//
//  Produces:
//    - application/json
//
//
//
//
//
// swagger:meta
package api
