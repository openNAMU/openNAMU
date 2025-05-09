from .tool.func import *

async def main_setting_head(num, skin_name = '', set_preview = 0):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
        if flask.request.method == 'POST' and set_preview == 0:
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
                curs.execute(db_change("update other set data = ? where name = ? and coverage = ?"), [flask.request.form.get('content', ''), info_d, coverage])
            else:
                curs.execute(db_change("insert into other (name, data, coverage) values (?, ?, ?)"), [info_d, flask.request.form.get('content', ''), coverage])

            await acl_check(tool = 'owner_auth', memo = 'edit_set (' + info_d + ')')

            if skin_name == '':
                return redirect(conn, '/setting/' + end_r)
            else:
                return redirect(conn, '/setting/' + end_r + '/' + skin_name)
        else:
            title = ''
            start = ''
            form_action = ''
            data_preview = ''
            plus = ''

            if num == 4:
                curs.execute(db_change("select data from other where name = 'body'"))
                title = '_body'
                form_action = 'formaction="/setting/body/top"'
                data_preview = flask.request.form.get('content', '') if set_preview == 1 else ''
                plus = '''
                    <button id="opennamu_preview_button" type="submit" formaction="/setting_preview/body/top">''' + get_lang(conn, 'preview') + '''</button>
                    <hr class="main_hr">
                    <div id="opennamu_preview_area">''' + data_preview + '''</div>
                '''
            elif num == 7:
                curs.execute(db_change("select data from other where name = 'bottom_body'"))
                title = '_bottom_body'
                data_preview = flask.request.form.get('content', '') if set_preview == 1 else ''
                form_action = 'formaction="/setting/body/bottom"'
                plus = '''
                    <button id="opennamu_preview_button" type="submit" formaction="/setting_preview/body/bottom">''' + get_lang(conn, 'preview') + '''</button>
                    <hr class="main_hr">
                    <div id="opennamu_preview_area">''' + data_preview + '''</div>
                '''
            else:
                skin_list = ''
                for for_a in load_skin(conn, '', 1):
                    skin_list += '<a href="/setting/head/' + url_pas(for_a) + '">(' + html.escape(for_a) + ')</a> '
                    skin_list += '<a href="/setting/head/' + url_pas(for_a) + '-cssdark">(' + html.escape(for_a) + '-cssdark)</a> '

                curs.execute(db_change("select data from other where name = 'head' and coverage = ?"), [skin_name])
                title = '_head'
                start = '' + \
                    '<a href="/setting/head">(' + get_lang(conn, 'all') + ')</a> ' + \
                    skin_list + '''
                    <hr class="main_hr">
                    <span>
                        &lt;style&gt;CSS&lt;/style&gt;
                        <br>
                        &lt;script&gt;JS&lt;/script&gt;
                    </span>
                    <hr class="main_hr">
                '''

            if set_preview == 1:
                data = data_preview
            else:
                head = curs.fetchall()
                if head:
                    data = head[0][0]
                else:
                    data = ''

            if skin_name != '':
                sub_plus = ' (' + skin_name + ')'
            else:
                sub_plus = ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, data = 'main' + title, safe = 1), await wiki_set(), await wiki_custom(conn), wiki_css(['(HTML)' + sub_plus, 0])],
                data = '''
                    <form method="post">
                        ''' + start + '''
                        <textarea class="opennamu_textarea_500" placeholder="''' + get_lang(conn, 'enter_html') + '''" name="content" id="content">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        ''' + (get_lang(conn, 'main_css_warning') + '<hr class="main_hr">' if title == '_head' else '') + '''
                        <button id="opennamu_save_button" type="submit" ''' + form_action + '''>''' + get_lang(conn, 'save') + '''</button>
                        ''' + plus + '''
                    </form>
                ''',
                menu = [['setting', get_lang(conn, 'return')]]
            ))