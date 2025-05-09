from .tool.func import *

async def bbs_delete(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')
        
        bbs_name = db_data[0][0]
        
        bbs_num_str = str(bbs_num)

        if await acl_check('', 'owner_auth', '', '') == 1:
            return redirect(conn, '/bbs/in/' + bbs_num_str)
        
        if bbs_num_str == 0:
            return redirect(conn, '/bbs/in/' + bbs_num_str)
        
        if flask.request.method == 'POST':
            curs.execute(db_change('delete from bbs_data where set_id = ?'), [bbs_num_str])
            curs.execute(db_change('delete from bbs_set where set_id = ?'), [bbs_num_str])
            curs.execute(db_change('delete from bbs_data where set_id like ?'), [bbs_num_str + '-%'])
            
            return redirect(conn, '/bbs/main')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'bbs_delete'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + bbs_name + ')', 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <span>''' + get_lang(conn, 'delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'delete') + '''</button>
                    </form>
                '''),
                menu = [['bbs/set/' + bbs_num_str, get_lang(conn, 'return')]]
            ))