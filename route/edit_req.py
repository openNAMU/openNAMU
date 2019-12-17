from .tool.func import *

def edit_req_2(conn, name):
    curs = conn.cursor()

    ip = ip_check()
    get_ver = flask.request.args.get('r', None)
    if get_ver:
        section = None
    else:
        section = flask.request.args.get('section', None)

    if acl_check(name) == 1:
        if acl_check(name, 'edit_req') == 1 or re.search('^user:', name) or ban_check() == 1 or get_ver:
            return re_error('/ban')
    else:
        if not get_ver:
            return redirect('/edit/' + url_pas(name))
        else:
            get_ver = int(number_check(get_ver))
        
    if not get_ver:
        curs.execute(db_change("select data from data where title = ?"), [name])
        old = curs.fetchall()
        if not old:
            return redirect('/w/' + url_pas(name))
    else:
        curs.execute(db_change("select data, send, ip, date from history where title = ? and id = ? and type = 'req'"), [name, str(get_ver)])
        old = curs.fetchall()
        if not old:
            return redirect('/w/' + url_pas(name))
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        if slow_edit_check() == 1:
            return re_error('/error/24')

        today = get_time()
        if get_ver:
            content = old[0][0]
        else:
            content = flask.request.form.get('content', '')

            if flask.request.form.get('otent', '') == content:
                return redirect('/w/' + url_pas(name))
            
            if edit_filter_do(content) == 1:
                return re_error('/error/21')

        content = savemark(content)
        
        if old:
            leng = leng_check(len(flask.request.form.get('otent', '')), len(content))
            
            if section:
                content = old[0][0].replace(flask.request.form.get('otent', ''), content)
        else:
            leng = '+' + str(len(content))

        if get_ver:
            if old:
                curs.execute(db_change("update data set data = ? where title = ?"), [content, name])
            else:
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, content])

                curs.execute(db_change('select data from other where name = "count_all_title"'))
                curs.execute(db_change("update other set data = ? where name = 'count_all_title'"), [str(int(curs.fetchall()[0][0]) + 1)])

            curs.execute(db_change("select user from scan where title = ?"), [name])
            for scan_user in curs.fetchall():
                curs.execute(db_change("insert into alarm (name, data, date) values (?, ?, ?)"), [
                    scan_user[0],
                    ip + ' | <a href="/w/' + url_pas(name) + '">' + name + '</a> | Edit', 
                    today
                ])

            curs.execute(db_change("update history set type = '', send = ? where title = ? and id = ? and ip = ? and date = ? and type = 'req'"), [
                old[0][1] + ' (' + ip + ' pass)', 
                name,
                str(get_ver),
                old[0][2],
                old[0][3]
            ])
            
            curs.execute(db_change("delete from back where link = ?"), [name])
            curs.execute(db_change("delete from back where title = ? and type = 'no'"), [name])
            
            render_set(
                title = name,
                data = content,
                num = 1
            )
        else:
            history_plus(
                name,
                content,
                today,
                ip,
                flask.request.form.get('send', ''),
                leng,
                '',
                'req'
            )
        
        conn.commit()
        
        if get_ver:
            return redirect('/w/' + url_pas(name))
        else:
            return redirect('/recent_changes?set=req')
    else:            
        if old:
            data = old[0][0]
        else:
            data = ''
            
        data_old = data
        get_name = ''

        save_button = load_lang('edit_req') if not get_ver else load_lang('edit_req_check') 
        menu_plus = [[]]
        sub = load_lang('edit_req')
        disable = '' if not get_ver else 'disabled'

        curs.execute(db_change('select data from other where name = "edit_bottom_text"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            b_text = '<hr class=\"main_hr\">' + sql_d[0][0]
        else:
            b_text = ''

        curs.execute(db_change('select data from other where name = "edit_help"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            p_text = sql_d[0][0]
        else:
            p_text = load_lang('defalut_edit_help')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + sub + ')', 0])],
            data =  get_name + '''
                <form method="post">
                    <script>do_stop_exit();</script>
                    ''' + edit_button() + '''
                    <textarea rows="25" ''' + disable + ''' id="content" placeholder="''' + p_text + '''" name="content">''' + html.escape(re.sub('\n$', '', data)) + '''</textarea>
                    <textarea id="origin" name="otent">''' + html.escape(re.sub('\n$', '', data_old)) + '''</textarea>
                    <hr class=\"main_hr\">
                    <input ''' + disable + ''' placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    ''' + captcha_get() + ip_warring() + '''
                    <button id="save" type="submit" onclick="go_save_zone = 1;">''' + save_button + '''</button>
                    <button id="preview" type="button" onclick="load_preview(\'''' + url_pas(name) + '\')">' + load_lang('preview') + '''</button>
                </form>
                ''' + b_text + '''
                <hr class=\"main_hr\">
                <div id="see_preview"></div>
            ''',
            menu = [['w/' + url_pas(name), load_lang('return')]] + menu_plus
        ))