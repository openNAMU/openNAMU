from .tool.func import *

def recent_history_add(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if admin_check() != 1:
            return re_error('/ban')

        if flask.request.method == 'POST':
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
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('history_add'), wiki_set(), wiki_custom(), wiki_css(['(' + name + ')', 0])],
                data = '''
                    <form method="post">
                        <textarea style="display: none;" id="opennamu_js_edit_origin"></textarea>
                        <textarea style="display: none;" id="opennamu_js_edit_textarea" name="content"></textarea>
                        
                        <div>''' + edit_button() + '''</div>
                        
                        <textarea id="opennamu_js_edit_textarea_view" class="content" placeholder="''' + p_text + '''"></textarea>
                        <hr class="main_hr">
                        
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + load_lang('name') + '''" name="get_ip" type="text">
                        <hr class="main_hr">
                        
                        <button id="opennamu_js_save" type="submit">''' + load_lang('save') + '''</button>
                        <button id="opennamu_js_preview" type="button">''' + load_lang('preview') + '''</button>
                    </form>
                    
                    <hr class="main_hr">
                    <div id="opennamu_js_preview_area"></div>
                ''',
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))