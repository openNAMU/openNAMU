package route

import "opennamu/route/tool"

func View_login_logout(config tool.Config) string {
    if !tool.IP_or_user(config.IP) {
        return tool.Get_redirect("/user")
    }

    Api_login_logout(config)

    return tool.Get_redirect("/user")
}