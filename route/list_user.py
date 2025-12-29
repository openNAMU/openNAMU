from .tool.func import *

async def list_user(arg_num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        list_data = '<ul>'

        curs.execute(db_change("select id, data from user_set where name = 'date' order by data desc limit ?, 50"), [sql_num])
        user_list = curs.fetchall()
        for data in user_list:
            list_data += '<li>'
            list_data += await ip_pas(data[0])
            list_data += ' | ' + data[1] if data[1] != '' else ''
            list_data += '</li>'

        list_data += '</ul>' + await get_next_page_bottom('/list/user/{}', arg_num, user_list)

        return await render_template(
            await get_lang('member_list'),
            list_data,
            0,
            [['other', await get_lang('return')]]
        )
