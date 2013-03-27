/**
 * init.go
 * クライアントからリクエストが来たら処理を振り分ける
 */

package okareader

import(
	"net/http"
)

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/atom", atom_test)
	http.HandleFunc("/addfolder", addFolder)
}
