package cherryConst

import (
	"fmt"
)

// component name
const (
	HandlerComponent         = "handler_component"
	SessionComponent         = "session_component"
	ORMComponent             = "db_orm_component"
	QueueComponent           = "db_queue_component"
	ConnectorPomeloComponent = "connector_pomelo_component"
	ConnectorSimpleComponent = "connector_simple_component"
)

const (
	ProfileNameFormat = "profile-%s.json"
	version           = "1.0.0"
)

var logo = `

░█████╗░██╗░░██╗███████╗██████╗░██████╗░██╗░░░██╗
██╔══██╗██║░░██║██╔════╝██╔══██╗██╔══██╗╚██╗░██╔╝
██║░░╚═╝███████║█████╗░░██████╔╝██████╔╝░╚████╔╝░
██║░░██╗██╔══██║██╔══╝░░██╔══██╗██╔══██╗░░╚██╔╝░░
╚█████╔╝██║░░██║███████╗██║░░██║██║░░██║░░░██║░░░
░╚════╝░╚═╝░░╚═╝╚══════╝╚═╝░░╚═╝╚═╝░░╚═╝░░░╚═╝░░░

game sever framework@v%s
`

func GetLOGO() string {
	return fmt.Sprintf(logo, Version())
}

func Version() string {
	return version
}
