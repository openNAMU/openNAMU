from .tool.func import *

async def bbs_w_set(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')
        else:
            bbs_name = db_data[0][0]

        bbs_num_str = str(bbs_num)

        return await render_template(
            await get_lang('bbs_set'),
            '' + \
                '<div id="opennamu_bbs_w_set"></div>' + \
                '<script defer src="/views/main_css/js/route/bbs_w_set.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_bbs_w_set(); });</script>' + \
            '',
            '(' + bbs_name + ')',
            [['bbs/in/' + bbs_num_str, await get_lang('return')]]
        )
