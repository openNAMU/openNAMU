from .tool.func import *

async def vote_list(list_type = 'normal', num = 1):    
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        data = ''
        if list_type == 'normal':
            data += '<a href="/vote/list/close">(' + await get_lang('close_vote_list') + ')</a>'
            sub = 0
            curs.execute(db_change('select name, id, type from vote where type = "open" or type = "n_open" limit ?, 50'), [sql_num])
        else:
            data += '<a href="/vote">(' + await get_lang('open_vote_list') + ')</a>'
            sub = '(' + await get_lang('closed') + ')'
            curs.execute(db_change('select name, id, type from vote where type = "close" or type = "n_close" limit ?, 50'), [sql_num])

        data += '<ul>'

        data_list = curs.fetchall()
        for i in data_list:
            if list_type == 'normal':
                open_select = await get_lang('open_vote') if i[2] == 'open' else await get_lang('not_open_vote')
            else:
                open_select = await get_lang('open_vote') if i[2] == 'close' else await get_lang('not_open_vote')

            data += '<li><a href="/vote/' + i[1] + '">' + i[1] + '. ' + html.escape(i[0]) + '</a> (' + open_select + ')</li>'

        data += '</ul>'
        menu = []
        if list_type == 'normal':
            menu = [["vote/add", await get_lang('add_vote')]] if await acl_check('', 'vote') != 1 else []
            data += await get_next_page_bottom('/vote/list/{}', num, data_list)
        else:
            data += await get_next_page_bottom('/vote/list/close/{}', num, data_list)

        return await render_template(
            await get_lang('vote_list'),
            data,
            sub,
            [['other', await get_lang('return')]] + menu
        )
