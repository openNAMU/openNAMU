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
            
        if len(flask.request.form.get('send', None)) > 500:
            return re_error('/error/15')

        if flask.request.form.get('otent', None) == flask.request.form.get('content', None):
            return redirect('/w/' + url_pas(name))
            
        if edit_filter_do(flask.request.form.get('content', '')) == 1:
            return re_error('/error/21')

        today = get_time()
        content = savemark(flask.request.form.get('content', None))
        
        curs.execute("select data from data where title = ?", [name])
        old = curs.fetchall()
        if old:
            leng = leng_check(len(flask.request.form.get('otent', None)), len(content))
            
            if flask.request.args.get('section', None):
                content = old[0][0].replace(flask.request.form.get('otent', None), content)
                
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
            flask.request.form.get('send', None),
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
                test_data = '\n' + re.sub('\r\n', '\n', new[0][0]) + '\n'   
                
                section_data = re.findall('((?:={1,6}) ?(?:(?:(?!={1,6}\n).)+) ?={1,6}\n(?:(?:(?!(?:={1,6}) ?(?:(?:(?!={1,6}\n).)+) ?={1,6}\n).)*\n*)*)', test_data)
                data = section_data[int(flask.request.args.get('section', None)) - 1]
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
            action = ''
        else:
            get_name = ''
            action = '?section=' + flask.request.args.get('section', None)
            
        if flask.request.args.get('plus', None):
            curs.execute("select data from data where title = ?", [flask.request.args.get('plus', None)])
            get_data = curs.fetchall()
            if get_data:
                data = get_data[0][0]
                get_name = ''

        js_data = edit_help_button()

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('edit') + ')', 0])],
            data =  get_name + js_data[0] + '''
                    <form method="post" action="/edit/''' + url_pas(name) + action + '''">
                        ''' + js_data[1] + '''
                        <textarea id="content" rows="25" name="content">''' + html.escape(re.sub('\n$', '', data)) + '''</textarea>
                        <textarea style="display: none;" name="otent">''' + html.escape(re.sub('\n$', '', data_old)) + '''</textarea>
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + ip_warring() + '''
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                        <button id="preview" type="submit" formaction="/preview/''' + url_pas(name) + action + '">' + load_lang('preview') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('return')], ['delete/' + url_pas(name), load_lang('delete')], ['move/' + url_pas(name), load_lang('move')]]
        ))