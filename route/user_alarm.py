from .tool.func import *

def user_alarm():
    with get_db_connect() as conn:
        curs = conn.cursor()
    
        num = int(number_check(flask.request.args.get('num', '1')))
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0
    
        data = '<ul class="opennamu_ul">'
    
        curs.execute(db_change("select data, date from alarm where name = ? order by date desc limit ?, 50"), [ip_check(), sql_num])
        data_list = curs.fetchall()
        if data_list:
            data = '' + \
                '<a href="/alarm/delete">(' + load_lang('delete') + ')</a>' + \
                '<hr class="main_hr">' + \
                data + \
            ''
    
            for data_one in data_list:
                data_split = data_one[0].split(' | ')
                
                data += '<li>' + ip_pas(data_split[0]) + (' | ' + ' | '.join(data_split[1:]) if len(data_split) > 1 else '') + ' (' + data_one[1] + ')</li>'
    
        data += '' + \
            '</ul>' + \
            next_fix('/alarm?num=', num, data_list) + \
        ''
    
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('notice'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = [['user', load_lang('return')]]
        ))