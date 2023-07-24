from .tool.func import *

def recent_history_add(name = 'Test', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if admin_check() != 1:
            return re_error('/ban')

        if flask.request.method == 'POST' and do_type == '':
            admin_check(None, 'history_add (' + name + ')')

            today = get_time()
            content = flask.request.form.get('content', '')
            leng = '+' + str(len(content))

            history_plus(
                name,
                content,
                today,
                'Add:' + flask.request.form.get('get_ip', ''),
                flask.request.form.get('send', ''),
                leng,
                t_check = 'add',
                mode = 'add'
            )

            conn.commit()

            return redirect('/history/' + url_pas(name))
        else:
            curs.execute(db_change('select data from other where name = "edit_help"'))
            sql_d = curs.fetchall()
            p_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else load_lang('default_edit_help')

            send = ''
            get_ip = ''
            data = ''
            data_preview = ''
            if do_type == 'preview':
                data = flask.request.form.get('content', '')
                data = data.replace('\r', '')

                send = flask.request.form.get('send', '')
                get_ip = flask.request.form.get('get_ip', '')

                data_preview = render_set(
                    doc_name = name, 
                    doc_data = data
                )
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('history_add'), wiki_set(), wiki_custom(), wiki_css(['(' + name + ')', 0])],
                data = '''
                    <form method="post">                        
                        <div>''' + edit_button('opennamu_edit_textarea') + '''</div>
                        
                        <textarea id="opennamu_edit_textarea" class="opennamu_textarea_500" name="content" placeholder="''' + p_text + '''">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        
                        <input placeholder="''' + load_lang('why') + '''" name="send" value="''' + html.escape(send) + '''">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + load_lang('name') + '''" name="get_ip" value="''' + html.escape(get_ip) + '''">
                        <hr class="main_hr">
                        
                        <button id="opennamu_save_button" formaction="/history_add/''' + url_pas(name) + '''" type="submit">''' + load_lang('save') + '''</button>
                        <button id="opennamu_preview_button" formaction="/history_add_preview/''' + url_pas(name) + '''" type="submit">''' + load_lang('preview') + '''</button>
                    </form>
                    
                    <hr class="main_hr">
                    <div id="opennamu_preview_area">''' + data_preview + '''</div>
                ''',
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))