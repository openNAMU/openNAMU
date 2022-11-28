from .tool.func import *

def vote_list(list_type = 'normal', num = 1):    
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        data = ''
        if list_type == 'normal':
            data += '<a href="/vote/list/close">(' + load_lang('close_vote_list') + ')</a>'
            sub = 0
            curs.execute(db_change('select name, id, type from vote where type = "open" or type = "n_open" limit ?, 50'), [sql_num])
        else:
            data += '<a href="/vote">(' + load_lang('open_vote_list') + ')</a>'
            sub = '(' + load_lang('closed') + ')'
            curs.execute(db_change('select name, id, type from vote where type = "close" or type = "n_close" limit ?, 50'), [sql_num])

        data += '<ul class="opennamu_ul">'

        data_list = curs.fetchall()
        for i in data_list:
            if list_type == 'normal':
                open_select = load_lang('open_vote') if i[2] == 'open' else load_lang('not_open_vote')
            else:
                open_select = load_lang('open_vote') if i[2] == 'close' else load_lang('not_open_vote')

            data += '<li><a href="/vote/' + i[1] + '">' + i[1] + '. ' + html.escape(i[0]) + '</a> (' + open_select + ')</li>'

        data += '</ul>'
        if list_type == 'normal':
            data += ('<a href="/vote/add">(' + load_lang('add_vote') + ')</a>') if admin_check() == 1 else ''
            data += next_fix('/vote/list/', num, data_list)
        else:
            data += next_fix('/vote/list/close/', num, data_list)

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('vote_list'), wiki_set(), wiki_custom(), wiki_css([sub, 0])],
            data = data,
            menu = [['other', load_lang('return')]]
        ))