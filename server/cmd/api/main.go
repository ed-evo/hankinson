package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/webui"
)

func main() {
	ctx := context.Background()

	db := orm.GetClient(ctx)

	r := webui.NewWebRouter(db)

	// 4. Start Server
	log.Println("🚀 Hankinson backend spinning up on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
