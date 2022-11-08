from .tool.func import *

def recent_discuss_delete():
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        data_html = '' + \
            '<a href="/recent_discuss">(' + load_lang('normal') + ')</a>' + \
            '<ul id="opennamu_ul">' + \
        ''
        count = 0
        
        curs.execute(db_change("select title, sub, date, code, stop from rd order by date desc limit 50"))
        for data in curs.fetchall():
            if data[4] == '':
                stop_code = ''
            elif data[4] == 'O':
                stop_code = ' (' + load_lang('close') + ')'
            else:
                stop_code = ' (' + load_lang('stop') + ')'

            data_html += '' + \
                '<li>' + \
                    '<input type="checkbox" name="checkbox_' + str(count) + '"> ' + \
                    html.escape(data[1]) + ' (' + html.escape(data[0]) + ')' + stop_code + \
                '</li>' + \
            ''
            
            count += 1
            
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('recent_discussion'), wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('delete') + ') (' + load_lang('not_working') + ')', 0])],
            data = data_html,
            menu = 0
        ))