from .tool.func import *

def main_func_setting_head(num, skin_name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        if flask.request.method == 'POST':
            if num == 4:
                info_d = 'body'
                end_r = 'body/top'
                coverage = ''
            elif num == 7:
                info_d = 'bottom_body'
                end_r = 'body/bottom'
                coverage = ''
            else:
                info_d = 'head'
                end_r = 'head'
                if skin_name == '':
                    coverage = ''
                else:
                    coverage = skin_name

            curs.execute(db_change("select name from other where name = ? and coverage = ?"), [info_d, coverage])
            if curs.fetchall():
                curs.execute(db_change("update other set data = ? where name = ? and coverage = ?"), [
                    flask.request.form.get('content', ''),
                    info_d,
                    coverage
                ])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values (?, ?, ?)"), [info_d, flask.request.form.get('content', ''), coverage])

            conn.commit()

            admin_check(None, 'edit_set (' + info_d + ')')

            if skin_name == '':
                return redirect('/setting/' + end_r)
            else:
                return redirect('/setting/' + end_r + '/' + skin_name)
        else:
            if num == 4:
                curs.execute(db_change("select data from other where name = 'body'"))
                title = '_body'
                start = ''
                plus = '''
                    <button id="preview" type="button" onclick="load_raw_preview(\'content\', \'see_preview\')">''' + load_lang('preview') + '''</button>
                    <hr class="main_hr">
                    <div id="see_preview"></div>
                '''
            elif num == 7:
                curs.execute(db_change("select data from other where name = 'bottom_body'"))
                title = '_bottom_body'
                start = ''
                plus = '''
                    <button id="preview" type="button" onclick="load_raw_preview(\'content\', \'see_preview\')">''' + load_lang('preview') + '''</button>
                    <hr class="main_hr">
                    <div id="see_preview"></div>
                '''
            else:
                curs.execute(db_change("select data from other where name = 'head' and coverage = ?"), [skin_name])
                title = '_head'
                start = '' + \
                    '<a href="?">(' + load_lang('all') + ')</a> ' + \
                    ' '.join(['<a href="/setting/head/' + i + '">(' + i + ')</a>' for i in load_skin('', 1)]) + '''
                    <hr class="main_hr">
                    <span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span>
                    <hr class="main_hr">
                '''
                plus = ''

            head = curs.fetchall()
            if head:
                data = head[0][0]
            else:
                data = ''

            if skin_name != '':
                sub_plus = ' (' + skin_name + ')'
            else:
                sub_plus = ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang(data = 'main' + title, safe = 1), wiki_set(), wiki_custom(), wiki_css(['(HTML)' + sub_plus, 0])],
                data = '''
                    <form method="post">
                        ''' + start + '''
                        <textarea rows="25" placeholder="''' + load_lang('enter_html') + '''" name="content" id="content">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        ''' + plus + '''
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))