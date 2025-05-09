from .tool.func import *

async def bbs_w_hide(bbs_num = '', post_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')

        bbs_name = db_data[0][0]
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)

        if await acl_check('', 'bbs_auth', '', '') == 1:
            return redirect(conn, '/bbs/in/' + bbs_num_str)
        
        if flask.request.method == 'POST':
            pass
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'bbs_post_hide'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + bbs_name + ')' + ' (' + post_num_str + ')', 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <button type="submit">''' + get_lang(conn, 'hide') + '''</button>
                    </form>
                '''),
                menu = [['bbs/w/' + bbs_num_str + '/' + post_num_str, get_lang(conn, 'return')]]
            ))