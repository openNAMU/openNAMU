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
            return await render_template(
                await get_lang('bbs_delete'),
                await render_simple_set('''
                    <form method="post">
                        <span>''' + await get_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('delete') + '''</button>
                    </form>
                '''),
                '(' + bbs_name + ')',
                [['bbs/set/' + bbs_num_str, await get_lang('return')]]
            )
