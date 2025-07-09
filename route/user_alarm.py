from .tool.func import *

async def user_alarm():
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        num = int(number_check(flask.request.args.get('num', '1')))
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0
    
        data = '<ul>'

        ip = ip_check()
    
        curs.execute(db_change("select data, date, readme, id from user_notice where name = ? order by date desc limit ?, 50"), [ip, sql_num])
        data_list = curs.fetchall()
        if data_list:
            data = '' + \
                '<a href="/alarm/delete">(' + await get_lang('delete') + ')</a>' + \
                '<hr class="main_hr">' + \
                data + \
            ''
    
            for data_one in data_list:
                data_split = data_one[0].split(' | ')
                data_style = ''
                if data_one[2] == '1':
                    data_style = 'opacity: 0.75;'
                
                data += '' + \
                    '<li style="' + data_style + '">' + \
                        await ip_pas(data_split[0]) + (' | ' + ' | '.join(data_split[1:]) if len(data_split) > 1 else '') + \
                        ' | ' + data_one[1] + \
                        ' <a href="/alarm/delete/' + url_pas(data_one[3]) + '">(' + await get_lang('delete') + ')</a>' + \
                    '</li>' + \
                ''

        curs.execute(db_change("update user_notice set readme = '1' where name = ?"), [ip])
    
        data += '' + \
            '</ul>' + \
            await get_next_page_bottom('/alarm?num={}', num, data_list) + \
        ''
    
        return easy_minify(flask.render_template(await skin_check(),
            imp = [await get_lang('notice'), await wiki_set(), await wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = [['user', await get_lang('return')]]
        ))