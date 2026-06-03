package core

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	// Initialize GORM adapter
	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin adapter: %v", err)
	}

	// Define RBAC model
	text := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	m, err := model.NewModelFromString(text)
	if err != nil {
		log.Fatalf("Failed to load Casbin model: %v", err)
	}

	Enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	err = Enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	Logger.Info("Casbin enforcer initialized successfully")
}
