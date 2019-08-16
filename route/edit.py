from .tool.func import *

def edit_2(conn, name):
    curs = conn.cursor()

    ip = ip_check()
    if acl_check(name) == 1:
        return re_error('/ban')
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if flask.request.form.get('otent', '') == flask.request.form.get('content', ''):
            return redirect('/w/' + url_pas(name))
            
        if edit_filter_do(flask.request.form.get('content', '')) == 1:
            return re_error('/error/21')

        today = get_time()
        content = savemark(flask.request.form.get('content', ''))
        
        curs.execute("select data from data where title = ?", [name])
        old = curs.fetchall()
        if old:
            leng = leng_check(len(flask.request.form.get('otent', '')), len(content))
            
            if flask.request.args.get('section', None):
                content = old[0][0].replace(flask.request.form.get('otent', ''), content)
                
            curs.execute("update data set data = ? where title = ?", [content, name])
        else:
            leng = '+' + str(len(content))
            
            curs.execute("insert into data (title, data) values (?, ?)", [name, content])

        curs.execute("select user from scan where title = ?", [name])
        for _ in curs.fetchall():
            curs.execute("insert into alarm (name, data, date) values (?, ?, ?)", [ip, ip + ' - <a href="/w/' + url_pas(name) + '">' + name + '</a> (Edit)', today])

        history_plus(
            name,
            content,
            today,
            ip,
            flask.request.form.get('send', ''),
            leng
        )
        
        curs.execute("delete from back where link = ?", [name])
        curs.execute("delete from back where title = ? and type = 'no'", [name])
        
        render_set(
            title = name,
            data = content,
            num = 1
        )
        
        conn.commit()
        
        return redirect('/w/' + url_pas(name))
    else:            
        curs.execute("select data from data where title = ?", [name])
        new = curs.fetchall()
        if new:
            if flask.request.args.get('section', None):
                data = re.sub('\n(?P<in>={1,6})', '<br>\g<in>', '\n' + re.sub('\r\n', '\n', new[0][0]) + '\n')
                i = 0

                while 1:
                    g_data = re.search('((?:<br>)(?:(?:(?!\n|<br>).)+)(?:\n*(?:(?:(?!<br>).)+\n*)+)?)', data)
                    if g_data:
                        if int(flask.request.args.get('section', '1')) - 1 == i:
                            data = re.sub('<br>(?P<in>={1,6})', '\n\g<in>', g_data.groups()[0])
                            
                            break
                        else:
                            data = re.sub('((?:<br>)(?:(?:(?!\n|<br>).)+)(?:\n*(?:(?:(?!<br>).)+\n*)+)?)', '\n', data, 1)

                        i += 1
                    else:
                        break
            else:
                data = new[0][0]
        else:
            data = ''
            
        data_old = data
        
        if not flask.request.args.get('section', None):
            get_name =  '''
                <a href="/manager/15?plus=''' + url_pas(name) + '">(' + load_lang('load') + ')</a> <a href="/edit_filter">(' + load_lang('edit_filter_rule') + ''')</a>
                <hr class=\"main_hr\">
            '''
        else:
            get_name = ''
            
        if flask.request.args.get('plus', None):
            curs.execute("select data from data where title = ?", [flask.request.args.get('plus', 'test')])
            get_data = curs.fetchall()
            if get_data:
                data = get_data[0][0]
                get_name = ''

        curs.execute('select data from other where name = "edit_bottom_text"')
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            b_text = '<hr class=\"main_hr\">' + sql_d[0][0]
        else:
            b_text = ''

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('edit') + ')', 0])],
            data =  get_name + '''
                <form method="post">
                    ''' + edit_button() + '''
                    <textarea id="content" rows="25" id="content" name="content">''' + html.escape(re.sub('\n$', '', data)) + '''</textarea>
                    <textarea style="display: none;" name="otent">''' + html.escape(re.sub('\n$', '', data_old)) + '''</textarea>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    ''' + captcha_get() + ip_warring() + '''
                    <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    <button id="preview" type="button" onclick="do_preview(\'''' + name + '\')">' + load_lang('preview') + '''</button>
                </form>
                ''' + b_text + '''
                <hr class=\"main_hr\">
                <div id="see_preview"></div>
            ''',
            menu = [['w/' + url_pas(name), load_lang('return')], ['delete/' + url_pas(name), load_lang('delete')], ['move/' + url_pas(name), load_lang('move')]]
        ))