from .tool.func import *

def edit_revert_2(name):
    

    num = int(number_check(flask.request.args.get('num', '1')))

    sqlQuery("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
    if sqlQuery("fetchall") and admin_check(6) != 1:
        return re_error('/error/3')

    if acl_check(name) == 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)
    
        sqlQuery("select data from history where title = ? and id = ?", [name, str(num)])
        data = sqlQuery("fetchall")
        if data:
            if edit_filter_do(data[0][0]) == 1:
                return re_error('/error/21')

        sqlQuery("delete from back where link = ?", [name])
        sqlQuery("commit")
        
        if data:                                
            sqlQuery("select data from data where title = ?", [name])
            data_old = sqlQuery("fetchall")
            if data_old:
                leng = leng_check(len(data_old[0][0]), len(data[0][0]))
                sqlQuery("update data set data = ? where title = ?", [data[0][0], name])
            else:
                leng = '+' + str(len(data[0][0]))
                sqlQuery("insert into data (title, data) values (?, ?)", [name, data[0][0]])
                
            history_plus(
                name, 
                data[0][0], 
                get_time(), 
                ip_check(), 
                flask.request.form.get('send', ''), 
                leng,
                'r' + str(num) + ''
            )

            render_set(
                title = name,
                data = data[0][0],
                num = 1
            )
            
            sqlQuery("commit")
            
        return redirect('/w/' + url_pas(name))
    else:
        sqlQuery("select title from history where title = ? and id = ?", [name, str(num)])
        if not sqlQuery("fetchall"):
            return redirect('/w/' + url_pas(name))

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('revert') + ')', 0])],
            data =  '''
                    <form method="post">
                        <span>r''' + flask.request.args.get('num', '0') + '''</span>
                        <hr class=\"main_hr\">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('revert') + '''</button>
                    </form>
                    ''',
            menu = [['history/' + url_pas(name), load_lang('history')], ['recent_changes', load_lang('recent_change')]]
        ))