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
            return await render_template(
                await get_lang('bbs_post_hide'),
                await render_simple_set('''
                    <form method="post">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('hide') + '''</button>
                    </form>
                '''),
                '(' + bbs_name + ')' + ' (' + post_num_str + ')',
                [['bbs/w/' + bbs_num_str + '/' + post_num_str, await get_lang('return')]]
            )
