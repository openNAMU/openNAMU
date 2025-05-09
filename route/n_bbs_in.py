from .tool.func import *

async def bbs_in(bbs_num = 1, page = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        bbs_num_str = str(bbs_num)

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')
    
        bbs_name = db_data[0][0]

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [bbs_name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'bbs') + ') (' + str(page) + ')', 0])],
            data = '' + \
                '<div id="opennamu_bbs_in"></div>' + \
                '<script defer src="/views/main_css/js/route/bbs_in.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_bbs_in(); });</script>' + \
            '',
            menu = [['bbs/main', get_lang(conn, 'return')], ['bbs/edit/' + bbs_num_str, get_lang(conn, 'add')], ['bbs/set/' + bbs_num_str, get_lang(conn, 'bbs_set')]]
        ))