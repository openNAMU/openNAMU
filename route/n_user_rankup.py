from .tool.func import *

def user_rankup():
    with get_db_connect() as conn:
        ip = ip_check()

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'rankup'), wiki_set(conn), wiki_custom(conn), wiki_css(['(' + ip + ')', 0])],
            data = '' + \
                '<div id="opennamu_user_rankup"></div>' + \
                '<script defer src="/views/main_css/js/route/user_rankup.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_user_rankup(); });</script>' + \
            '',
            menu = [['user/' + url_pas(ip)]]
        ))