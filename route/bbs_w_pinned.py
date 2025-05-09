from .tool.func import *

from .go_api_bbs_w import api_bbs_w

async def bbs_w_pinned(bbs_num = '', post_num = ''):
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
        
        temp_dict = await api_bbs_w(bbs_num_str + '-' + post_num_str)
        if not 'user_id' in temp_dict:
            return redirect(conn, '/bbs/main')
        
        if flask.request.method == 'POST':
            curs.execute(db_change('select set_data from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            if not curs.fetchall():
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', ?, ?, ?)"), [post_num_str, bbs_num_str, get_time()])
            else:
                curs.execute(db_change('delete from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            
            return redirect(conn, '/bbs/in/' + bbs_num_str)
        else:
            curs.execute(db_change('select set_data from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            pinned = get_lang(conn, 'pinned') if not curs.fetchall() else get_lang(conn, 'pinned_release')

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'bbs_post_pinned'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + bbs_name + ')' + ' (' + post_num_str + ')', 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <button type="submit">''' + pinned + '''</button>
                    </form>
                '''),
                menu = [['bbs/w/' + bbs_num_str + '/' + post_num_str, get_lang(conn, 'return')]]
            ))